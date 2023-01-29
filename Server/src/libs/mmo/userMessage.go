package mmo

import (
	. "Src/libs/command"
	"Src/libs/db"
	"Src/libs/logger"
	"Src/libs/mmo/data"
	"Src/libs/packet"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	LOGINKEYERR   = 1
	LOGINMYSQLERR = 2
	LOGINTABLEERR = 3
	LOGINUIDERR   = 4 //UID不存在 账号或密码错误
)

const (
	MAX_USER_ID  = 99999
	MAX_TABLE_ID = 99999
	MAX_MONEY    = 99999
)

const (
	RESULT_SUCCESS = 0
)

//这一块做用户消息处理
func (u *User) Send(p *packet.ICPacket) bool {

	_, err := u.conn.Write(p.GetData())
	if err != nil {
		logger.Error("send failed err:%s", err)
		return false
	}
	logger.Error("tid:%d,uid:%d,send cmd:0x%0x", u.Tid, u.Id, p.GetCmd())
	logger.Error("send data:[%s]", hex.EncodeToString(p.GetData()))
	return true
}


//消息模块拆分 例如角色消息模块 它的命令字是在0x2001~0x2fff之间  消息命令字做与运算取头部数字0x2001&0xFF000 ==0x2000	分成不同的消息模块取处理
// switch modChoose {
// case 1:
// 	self.HandleBase()
// case 2:
// 	self.HandleBag()
// case 3:
// 	self.HandlePool()
// case 4:
// 	self.HandleMap()
// case 5:
// 	self.HandleRelics()
// case 6:
// 	self.HandleRole()
// case 7:
// 	self.HandleWeapon()
// case 8:
// 	for _, v := range self.modManage {
// 		v.SaveData()
// 	}
// }
func (u *User) ProcessPkg(p packet.IPacket) {
	fmt.Println("Process Client Pkg")

	icp := p.(*packet.ICPacket)
	if icp.Decrypt() == -1 {
		return
	}
	cmd := icp.GetCmd()

	logger.Error("recv cmd[0x%x]", cmd)
	logger.Error("recv data:[%s]", hex.EncodeToString(icp.GetData()))

	switch cmd {
	//登录注册模块
	case CMD_CLIENT_HEARTBEAT_REQ: //心跳
		u.ProcessHeartBeat(icp)
	case CMD_CLIENT_LOGIN_REQ: //登录
		u.ProcessLogin(icp)
	case CMD_CLIENT_REGISTER_REQ: //注册
		u.ProcessRegister(icp)
	case CMD_CLIENT_LOGOUT_REQ: //登出
		u.ProcessLogout(icp)
	case CMD_CLIENT_GAME_ENTER_REQ:	//进入游戏
		u.ProcessGameEnter(icp)
	case CMD_CLIENT_GAME_LEAVE_REQ:	//退出游戏
		u.ProcessGameLeave(icp)



	//房间地图模块
	case CMD_CLIENT_GET_TABLES_INFO_REQ: //获取所有房间信息
		u.ProcessGetTablesInfo(icp)
	case CMD_CLIENT_ENTER_TABLE_REQ: //进入房间
		u.ProcessEnterTable(icp)
	case CMD_CLIENT_LEAVE_TABLE_REQ:
		u.ProcessLeaveTable(icp)


	//角色模块
	case CMD_CLIENT_GET_ROLES_INFO_REQ: //获取角色
		u.ProcessGetRolesInfo(icp)
	case CMD_CLIENT_CREATE_ROLE_REQ: //创建角色
		u.ProcessCreateRole(icp)
	case CMD_CLIENT_DEL_ROLE_REQ: //删除角色
		u.ProcessDelRole(icp)

	//商店模块
	case CMD_CLIENT_GET_STORE_INFO_REQ:
		u.ProcessGetStoreInfo(icp)
	case CMD_CLIENT_STORE_BUY_REQ:
		u.ProcessStoreBuy(icp)
	case CMD_CLIENT_STORE_SELL_REQ:
		u.ProcessStoreSell(icp)

	//战斗模块
	case CMD_CLIENT_MOVE_REQ:
		u.ProcessMove(icp)
	}

	//
	//modChoose:=(cmd&0xF000)>>12 //计算cmd的前缀 获取消息所属的mod
	//switch modChoose {
	//case 1:
	//	u.HandleBase(cmd,icp)
	//case 2:
	//	self.HandleBag()
	//case 3:
	//	self.HandlePool()
	//case 4:
	//	self.HandleMap()
	//case 5:
	//	self.HandleRelics()
	//case 6:
	//	self.HandleRole()
	//case 7:
	//	self.HandleWeapon()
	//case 8:
	//	for _, v := range self.modManage {
	//		v.SaveData()
		//}

}


func (u *User) ProcessHeartBeat(pack *packet.ICPacket) {
	//TODO 直接返回CMD_SERVER_HEARTBEAT_RESP
	logger.Error("MsgPing")
	u.lastTime = time.Now() //记录心跳时间

	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_HEARTBEAT_RESP)
	outICP.End()
	u.Send(outICP)
}
func (u *User) ProcessLogin(pack *packet.ICPacket) {
	//TODO	数据库验证账号密码获取唯一uid	如果有同uid登录用户将其踢出
	//登录 客户端会发送账号和密码 由数据库进行验证 返回唯一的uid
	account := pack.ReadString()
	password := pack.ReadString()

	ret := 0
	if db.GetDB().Mysqldb == nil {
		db.GetDB().InitMysql()
	}
	uid := int32(db.GetDB().GetUserId(account, password))

	if uid == 0 {
		//登录错误 返回错误码
		ret = 1
	}




	//TODO 登录结果处理
	if ret != 0 {
		logger.Error("client login fail account:[%s],password:[%s],ret:[%d]", account, password, ret)
	} else {

		logger.Error("client login succ account:[%s],password:[%s]", account, password)
		//登录成功 发送登录结果 uid 玩家列表 踢出上一个登录的同uid用户
		GetGameServer().KickUser(uid)
		u.Id = uid
		//进入游戏大厅
		u.EnterGameHall()
	}


	//TODO 发送自定义数据包
	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_LOGIN_RESP)
	outICP.WriteByte(byte(ret))
	outICP.WriteInt32(uid)
	outICP.End()
	u.Send(outICP)

}

func (u *User) ProcessRegister(pack *packet.ICPacket) {
	//TODO
	//注册 客户端会发送 邮箱 账号和密码 由数据库进行验证
	account := pack.ReadString()
	password := pack.ReadString()

	ret := 0
	//如果该用户已存在 返回错误 否则注册成功
	uid := int32(db.GetDB().GetUserId(account, password))
	if uid != 0 {
		//注册失败 用户已存在
		ret = 1
	} else {
		if !db.GetDB().InsertUser(account, password) {
			//注册失败 数据库添加失败
			ret = 2
		}
	}


	if ret != 0 {
		logger.Error("client register fail account:[%s],password:[%s],ret:[%d]", account, password, ret)
	}else {
		logger.Error("client register succ account:[%s],password:[%s]", account, password)
	}

	//TODO 发送自定义数据包
	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_REGISTER_RESP)
	outICP.WriteByte(byte(ret))
	outICP.End()
	u.Send(outICP)


}
func (u *User) ProcessLogout(pack *packet.ICPacket) {
	GetGameServer().KickUser(u.Id)
	u.Id = -1
	u.LeaveGameHall()
}

func (u *User) ProcessGameEnter(pack *packet.ICPacket) {
	ret:=0

	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_GAME_ENTER_RESP)

	characterIdx:=pack.ReadInt32()
	logger.Error("character id:%d, Enter Game",characterIdx);
	character,ok:=u.UserInfo.Player.Characters[int(characterIdx)]
	if !ok {
		//角色不存在
		ret = 1
	}else{

	
	outICP.WriteByte(byte(ret))
	
	outICP.WriteInt32(int32(character.Id))
	outICP.WriteString(character.Name)
	outICP.WriteInt32(int32(character.Level))
	outICP.WriteInt64(character.Exp)
	outICP.WriteInt64(character.Gold)

	outICP.WriteInt32(int32(character.MapId))
	outICP.WriteInt32(int32(character.Ride))

	outICP.WriteInt32(int32(character.Class))
	outICP.WriteInt32(int32(character.Type))

	outICP.WriteInt32(int32(character.EntityId))
	outICP.WriteInt32(int32(character.ConfigId))
	}

	

	outICP.End()

	u.Send(outICP)
}

func (u *User) ProcessGameLeave(pack *packet.ICPacket) {

}

func (u *User) ProcessGetRolesInfo(pack *packet.ICPacket) {
	//TODO 返回用户所有角色基本信息

	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_GET_ROLES_INFO_RESP)

	totalNum:=len(u.GetModRole().RoleInfo)
	outICP.WriteInt32(int32(totalNum))

	for _,v:=range u.GetModRole().RoleInfo {
		outICP.WriteInt32(v.RoleId)
		outICP.WriteInt32(v.RoleType)
		outICP.WriteInt32(v.RoleLevel)
		outICP.WriteInt64(v.RoleNum)
		outICP.WriteString(v.RoleName)
	}

	outICP.End()

	u.Send(outICP)

}
func (u *User) ProcessCreateRole(pack *packet.ICPacket) {
	//TODO	创建成功	给用户添加该角色信息
	ret:=0
	//TCharacter character = new TCharacter()
	//{
	//	Name = request.Name,
	//	Class = (int)request.Class,
	//	TID = (int)request.Class,
	//	Level = 1,
	//	MapID = 1,
	//	MapPosX = 1500, //初始出生位置X
	//	MapPosY = 2600, //初始出生位置Y
	//	MapPosZ = 820,
	//	Gold = 100000, //初始10万金币
	//	HP=100,
	//	MP =50,
	//	Equips = new byte[28]
	//};

	character:=new(data.NCharacterInfo)
	character.InitData()

	//character.Id = pack.ReadInt()
	character.Name = pack.ReadString()
	character.Class = pack.ReadInt()
	character.Level =1
	character.MapId =1

	if !db.GetDB().InsertRole(u.Id, int32(character.Class), character.Name) {
		//注册失败 数据库添加失败
		ret = 1
	}
	character.Id = int(db.GetDB().GetRoleId(u.Id, int32(character.Class), character.Name))



	character.Entity.Position.X = 1500
	character.Entity.Position.Y = 2600
	character.Entity.Position.Z = 150

	character.Gold = 10000

	character.AttrDynamic.HP = 100
	character.AttrDynamic.MP = 50

	character.Equips = make([]byte,28)


	bag:=new(data.NBagInfo)
	bag.Items = make([]byte,0)
	bag.Unlocked = 20

	character.Bag =bag

	item:=new(data.NItemInfo)
	item.Id =1
	item.Count =20
	character.Items[item.Id]=item

	u.UserInfo.Player.Characters[character.Id]=character
	logger.Error("create character id:%d",character.Id);

	//创建角色成功	发送角色列表
	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_CREATE_ROLE_RESP)
	outICP.WriteByte(byte(ret))


	if ret==0 {
		outICP.WriteInt32Ex(len(u.UserInfo.Player.Characters))

		for _, v := range u.UserInfo.Player.Characters {
			outICP.WriteInt32Ex(v.Id)
			outICP.WriteString(v.Name)
			outICP.WriteInt32Ex(data.CHARACTER_TYPE_PLAYER)
			outICP.WriteInt32Ex(v.Class)
			outICP.WriteInt32Ex(v.ConfigId)

		}
	}

	outICP.End()
	u.Send(outICP)








	//对用户id
	//if u.Id < 0 || u.Id > MAX_USER_ID {
	//	return
	//}
	//
	//char_class := pack.ReadInt32()	//角色职业
	//char_name := pack.ReadString()	//角色名字
	//
	//ret := 0
	//
	//roleId := int32(db.GetDB().GetRoleId(u.Id, char_class, char_name))
	//
	//if roleId != 0 {
	//	//角色已存在 返回错误码
	//	ret = 1
	//} else {
	//	//角色不存在 创建
	//	if !db.GetDB().InsertRole(u.Id, char_class, char_name) {
	//		//注册失败 数据库添加失败
	//		ret = 2
	//	}
	//}
	//
	//outICP := packet.NewICPacket()
	//outICP.Begin(CMD_SERVER_CREATE_ROLE_RESP)
	//outICP.WriteByte(byte(ret))
	//outICP.End()
	//u.Send(outICP)
	//
	//if ret != 0 {
	//	logger.Error("client create role fail uid:[%d],role_type:[%d],role_name:[%s],ret:[%d]", u.Id, char_class, char_name)
	//	return
	//}
	//
	//logger.Error("client create role succ uid:[%d],role_type:[%d],role_name:[%s]", u.Id, char_class, char_name)
	//u.GetModRole().AddItem(char_class, 1)

}

func (u *User) ProcessDelRole(pack *packet.ICPacket) {
	//TODO	删除成功	给用户删除该角色信息
	roleId:=pack.ReadInt32()
	roleNum:=pack.ReadInt64()

	u.GetModRole().RemoveItem(roleId,roleNum)

	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_DEL_ROLE_RESP)
	outICP.WriteByte(0)	//成功返回0
	outICP.End()

	u.Send(outICP)
}

func (u *User) ProcessGetTablesInfo(pack *packet.ICPacket) {
	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_GET_TABLES_INFO_RESP)
	//发送字段 房间的数量	for tid name
	tableNum := 0
	GetGameServer().Tables.Range(func(k, v interface{}) bool {
		tableNum++
		return true
	})
	outICP.WriteByte(byte(tableNum))

	GetGameServer().Tables.Range(func(k, v interface{}) bool {


		outICP.WriteInt32(v.(*Table).Id)
		outICP.WriteString(v.(*Table).Name)
		outICP.WriteInt32(v.(*Table).Type)
		//以及在线人数
		if v.(*Table).Id==MAP_TYPE_ONLINE {
			outICP.WriteInt32(u.server.GetOnlineMapCount())
		}else {
			outICP.WriteInt32(int32(len(v.(*Table).GetAllUsers())))
		}
		return true
	})

	outICP.End()
	u.Send(outICP)
}
func (u *User) ProcessEnterTable(pack *packet.ICPacket) {
	//TODO
	//字段：uid tid role info
	if u.Id < 0 || u.Id > MAX_USER_ID {
		return
	}
	uid := pack.ReadInt32()
	tid := pack.ReadInt32()
	roleid:=pack.ReadInt32()
	name :=pack.ReadString()

	ret := 0
	//查询是否有该房间 如果有该用户加入房间
	table := GetGameServer().GetTable(tid)
	if table == nil {
		ret = 1 //table不存在
	} else {
		user := table.GetUser(uid)
		if user != nil {
			ret = 2 //用户已经进入table
		} else {
			table.AddUser(u)
			user = table.GetUser(uid)
			if user == nil {
				ret = 3 //用户已经进入table失败
			}
		}

	}



	if ret != 0 {
		//向用户单播

		logger.Error("uid:[%d] enter table:[%d] fail", uid, tid)
		return
	}
	u.Tid = tid
	logger.Error("uid:[%d] enter table:[%d] succ", uid, tid)


	//向房间所有用户广播进入的玩家uid tid role_id class_id 即可 客户端根据
	BroadcastPack := packet.NewICPacket()
	BroadcastPack.Begin(CMD_SERVER_ENTER_TABLE_RESP)
	BroadcastPack.WriteByte(byte(ret))
	BroadcastPack.WriteInt32(uid)
	BroadcastPack.WriteInt32(tid)
	BroadcastPack.WriteInt32(roleid)
	BroadcastPack.WriteString(name)
	BroadcastPack.End()

	u.server.BroadcastTable(u.Tid,BroadcastPack)
}

func (u *User) ProcessLeaveTable(pack *packet.ICPacket) {
	//TODO
	uid := pack.ReadInt32()
	tid := pack.ReadInt32()
	table := GetGameServer().GetTable(tid)
	if table != nil {
		table.DelUser(uid)
	}
	logger.Error("uid:[%d] leave table:[%d]", uid, tid)

}
func (u *User) ProcessMove(pack *packet.ICPacket) {
	uid := pack.ReadInt32()
	tid := pack.ReadInt32()
	x := pack.ReadFloat32()
	y := pack.ReadFloat32()
	z := pack.ReadFloat32()

	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_MOVE_RESP)
	outICP.WriteInt32(uid)
	outICP.WriteInt32(tid)

	outICP.WriteFloat32(x)
	outICP.WriteFloat32(y)
	outICP.WriteFloat32(z)
	outICP.End()

	table := GetGameServer().GetTable(tid)
	if table != nil {
		table.Broadcast(outICP)
	}
}


func (u *User) ProcessGetStoreInfo(pack *packet.ICPacket) {
	//TODO	获取商店信息	返回所有商店列表信息	(根据商店种类	返回不同的商品信息)

	store_type:=pack.ReadByte()


	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_GET_STORE_INFO_RESP)
	switch store_type {

	case STORE_TYPE_NORMAL:
		goodsLen:=len(u.GetModStore().GoodsList)
		outICP.WriteInt32(int32(goodsLen))
		for _,v:=range u.GetModStore().GoodsList{
			outICP.WriteInt32(v.ItemId)
			outICP.WriteString(v.Name)
			outICP.WriteInt64(v.Price)
			outICP.WriteInt64(v.SellPrice)
			outICP.WriteInt64(v.ItemNum)
		}


	}

	outICP.End()

	u.Send(outICP)



}

func (u *User) ProcessStoreBuy(pack *packet.ICPacket) {
	//TODO	根据购买物品的id和数量 返回购买后用户的金钱	以及购买的物品及数量
	itemId:=pack.ReadInt32()
	num:=pack.ReadInt64()

	u.GetModStore().Buy(itemId,num)

	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_STORE_BUY_RESP)
	outICP.WriteInt32(itemId)
	outICP.WriteInt64(num)
	outICP.WriteInt64(u.Money)
}

func (u *User) ProcessStoreSell(pack *packet.ICPacket) {
	//TODO	根据购买物品的id和数量 返回购买后用户的金钱	以及购买的物品及数量
	itemId:=pack.ReadInt32()
	num:=pack.ReadInt64()

	u.GetModStore().Sell(itemId,num)

	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_STORE_SELL_RESP)
	outICP.WriteInt32(itemId)
	outICP.WriteInt64(num)
	outICP.WriteInt64(u.Money)
}
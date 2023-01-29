package GameServer

import (
	"Src/libs/config"
	"Src/libs/db"
	"Src/libs/logger"
	"Src/libs/packet"
	"fmt"
	"net"
	"strconv"
	"sync"
)

//监听类型
const (
	LISTEN_TYPE_CLIENT = 1
	LISTEN_TYPE_SERVER = 2
	LISTEN_TYPE_ADMIN  = 3
)

type GameServer struct {
	DBService *db.BaseDb //数据库连接
	// go中map数据结构不是线程安全的，即多个goroutine同时操作一个map，则会报错 sync.Map适合大量读，少量写
	m_ServerUserList sync.Map //所有用户列表
	m_CloseUserList  sync.Map //不正常断连接用户列表

	//系统数据模块	使用泛型方便做扩展
	modManage map[string]SysModBase
	localPath string //用户数据本地保存路径
}

//单例
var gameserver *GameServer

func GetGameServer() *GameServer {
	if gameserver == nil {
		gameserver = &GameServer{
			modManage: map[string]SysModBase{
				MOD_GUILD:     new(ModGuild),
				SYS_MOD_MATCH: new(ModMatch),
			},
		}

	}
	return gameserver
}

func (s *GameServer) Init() bool {
	//s.DBService = db.GetDB()

	return true
}

func (s *GameServer) Listen(address string, lType int) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("Starting TCP Server failed,err:%s", err.Error())
		return
	}
	defer listener.Close()

	logger.Error("Server start listen addr[%s]", address)

	for {
		conn, err := listener.Accept()

		if nil != err {
			logger.Error("listener.AcceptTCP() failed, err[%s]", err.Error())
			return
		}
		logger.Error("New client connected addr[%s] lType[%d]", conn.RemoteAddr().String(), lType)

		switch lType {

		case LISTEN_TYPE_CLIENT:
			pSocket := NewSocketHandler(conn, s)
			go ReceiveIcPkg(pSocket)

		}

	}
}

func (s *GameServer) Run() {
	//主线程开启监听 可使用goroutine开启协程并行监听其他端口
	s.Listen(config.GlobalConfig.Host, LISTEN_TYPE_CLIENT)
}

func (s *GameServer) ProcessClose(pSocket *SocketHandler) {

	logger.Error("ProcessClose")
	if pSocket.m_pGameUser != nil && pSocket.m_pGameUser.m_pSocket == pSocket {
		pSocket.m_pGameUser.m_pSocket = nil

		s.AddCloseUser(pSocket.m_pGameUser)

		logger.Error("uid:%d ProcessClose ", pSocket.m_pGameUser.m_nUserID)

	}

}

func (s *GameServer) ProcessHeartBeat(pack *packet.ICPacket, pSocket *SocketHandler) {

	logger.Error("ProcessHeartBeat")
	//每次心跳	重置心跳检测定时	超过一分钟定时触发
	if pSocket.m_pGameUser != nil {

		logger.Error("ResetCheckHeartBeatTimer")

		pSocket.m_pGameUser.ResetCheckHeartBeatTimer()

		go ProcessCheckHeartBeatTimeOut(pSocket.m_pGameUser)

	}

	outICP := packet.NewICPacket()
	outICP.Begin(SERVER_COMMAND_HEARTBEAT)
	outICP.End()
	pSocket.Send(outICP)
}

func ProcessCheckHeartBeatTimeOut(pUser *CGameUser) {
	//心跳超时处理 心跳超时加入close用户列表
	for {
		select {
		case <-pUser.m_CheckHeartBeat.C:
			GetGameServer().AddCloseUser(pUser)
			logger.Error("ProcessHeartBeatTimeOut")
		}
	}
}

func (s *GameServer) ProcessUserLogin(pack *packet.ICPacket, pSocket *SocketHandler) {
	uid := pack.ReadInt()

	if s.IsCloseUser(uid, pSocket) {
		return
	}
	pos:=Position{X: 0,Y:0,Angle: 0}
	s.AddNewUser(uid,pos, pSocket)
	logger.Error("User Login Succ")

}

func (s *GameServer) ProcessSearchUser(pack *packet.ICPacket, pSocket *SocketHandler) {
	mode := pack.ReadByte() //搜索模式

	var pUser *CGameUser = nil
	ret := 0
	if mode == 0 { //通过uid搜索	uid唯一
		uid := pack.ReadInt()

		pUser = s.GetUser(uid)
		if pUser != nil {
			ret = 1
		}
	}

	if mode == 1 { //通过name搜索	name唯一
		name := pack.ReadString()
		pUser = s.GetUserByName(name)
		if pUser != nil {
			ret = 1
		}
	}

	//返回用户搜索结果 可能多个
	outICP := packet.NewICPacket()
	outICP.Begin(SERVER_COMMAND_SEARCH_USER_RET)
	outICP.WriteByte(byte(ret)) //0：搜索用户不存在	1：搜索用户存在
	if ret == 1 {
		outICP.WriteInt32Ex(pUser.m_nUserID)
	}

	outICP.End()
	pSocket.Send(outICP)

}

func (s *GameServer) ProcessAddFriend(pack *packet.ICPacket, pSocket *SocketHandler) {
	uid := pack.ReadInt()
	pUser := s.GetUser(uid)
	if pUser == nil || pUser.m_pSocket == nil {
		return
	}

	//
	outICP := packet.NewICPacket()
	outICP.Begin(SERVER_COMMAND_SEND_FRIEND_INVITATION)

	outICP.WriteInt32Ex(pSocket.m_pGameUser.m_nUserID)

	outICP.End()
	pSocket.Send(outICP)

}

func (s *GameServer) ProcessAcceptFriendsInvitation(pack *packet.ICPacket, pSocket *SocketHandler) {
	uid := pack.ReadInt()
	pUser := s.GetUser(uid)
	if pUser == nil || pUser.m_pSocket == nil {
		return
	}
	ret := pSocket.m_pGameUser.GetModFriend().AddFriend(uid)

	//更新好友列表

	outICP := packet.NewICPacket()
	outICP.Begin(SERVER_COMMAND_SEND_FRIEND_INVITATION)
	outICP.WriteByte(byte(ret))
	outICP.WriteInt32Ex(pSocket.m_pGameUser.m_nUserID)
	outICP.End()
	pSocket.Send(outICP)

}

func (s *GameServer) ProcessReturnUserTeamInfo(pack *packet.ICPacket, pSocket *SocketHandler) {
	uid := pack.ReadInt()
	friendUid := pack.ReadInt()
	pUser := s.GetUser(uid)
	if pUser == nil {
		return
	}
	if !s.GetSysModMatch().HasTeam(pUser.m_nTeamId) {
		//没有队伍分配队伍
		s.GetSysModMatch().DelTeam(pUser.m_nTeamId)
		pUser.m_nTeamId = GenerateTeamId()
		s.GetSysModMatch().AddUserToTeam(pUser, pUser.m_nTeamId)
	}

	if friendUid != 0 { //好友邀请加入队伍
		//1.查找好友所在的队伍

	}

}

func (s *GameServer) AddNewUser(uid int,pos Position, pSocket *SocketHandler) {
	//pNewUser:= &CGameUser{
	//	m_nUserID: uid,
	//	m_pSocket: pSocket,
	//	m_CheckHeartBeat: nil,
	//}

	pNewUser := NewUser(uid,pos, pSocket)

	s.AddOnlineUser(pNewUser)
	pSocket.m_pGameUser = pNewUser

}

func (s *GameServer) IsCloseUser(uid int, pSocket *SocketHandler) bool { //是否刷新进入用户
	gameuser := s.GetUser(uid)

	if gameuser != nil {
		logger.Error("IsCloseUser yes!!!")

		if gameuser.m_pSocket != nil {
			gameuser.m_pSocket.m_pGameUser = nil
		}
		s.m_CloseUserList.Delete(uid)
		//gameuser->StopTimer();

		gameuser.m_pSocket = pSocket
		pSocket.m_pGameUser = gameuser

		//SendLoginErrResponse(gameuser.m_pSocket, EXISTERR);
		return true

	}
	return false
}

func (s *GameServer) IsUserExist(uid int, pSocket *SocketHandler) bool { //用户是否存在

	gameuser := s.GetUser(uid)

	if gameuser != nil {
		logger.Error("IsUserExist yes!!!")

		if gameuser.m_pSocket != nil {
			gameuser.m_pSocket.m_pGameUser = nil
		}
		s.m_CloseUserList.Delete(uid)
		//gameuser->StopTimer();

		gameuser.m_pSocket = pSocket
		pSocket.m_pGameUser = gameuser

		//SendLoginErrResponse(gameuser.m_pSocket, EXISTERR);
		return true

	}
	return false
}

// 添加上线用户
func (s *GameServer) AddOnlineUser(pUser *CGameUser) {
	s.m_ServerUserList.Store(pUser.m_nUserID, pUser)
}

func (s *GameServer) DelOnlineUser(uid int) {
	s.m_ServerUserList.Delete(uid)
}

//获取单个用户
func (s *GameServer) GetUser(uid int) *CGameUser {
	v, ok := s.m_ServerUserList.Load(uid)
	if ok {
		return v.(*CGameUser)
	}
	return nil
}

func (s *GameServer) GetUserByName(name string) *CGameUser {
	var pUser *CGameUser = nil
	s.m_ServerUserList.Range(func(k, v interface{}) bool {
		if v.(*CGameUser).Name == name {
			pUser = v.(*CGameUser)
		}
		return true
	})
	return pUser
}

func (s *GameServer) AddCloseUser(pUser *CGameUser) {
	s.m_CloseUserList.Store(pUser.m_nUserID, pUser)
}
func (s *GameServer) GetSysModGuild() *ModGuild {
	return s.modManage[MOD_GUILD].(*ModGuild)
}

func (s *GameServer) GetSysModMatch() *ModMatch {
	return s.modManage[SYS_MOD_MATCH].(*ModMatch)
}

func (s *GameServer) ProcessUserLoginNew(map_ map[string]interface{}, pSocket *SocketHandler) {
	m := make(map[string]interface{})
	slice := make([]map[string]interface{}, 0)

	account, _ := map_["account"].(string)
	password, _ := map_["password"].(string)
	uid := int(db.Db.GetUserId(account, password))

	if uid == 0 {
		return
	}
	pNewUser:=&CGameUser{
		m_nUserID: uid,
		Pos: Position{
			X: 400.00,
			Y: 700.00,
			Angle: 90,
		},
		Action: 0,
		m_pSocket: pSocket,
	}
	m["cmd"] = 0x3001
	m["uid"] = pNewUser.m_nUserID
	m["strkey"] = "test"
	m["x"] = pNewUser.Pos.X
	m["y"] = pNewUser.Pos.Y
	m["angle"] = pNewUser.Pos.Angle
	m["action"] = pNewUser.Action

	//用户登录进来	返回初始坐标	动作以及角度
	
	//发送所有已在线用户的坐标 动作以及角度
	userList :=s.GetAllUsers() 
	m["count"] = len(userList)
	for i := 0; i < len(userList); i++ {
		pUser := userList[i] 
		key:=fmt.Sprintf("uid_%d",i)
		m[key] = pUser.m_nUserID
		m["x"] = pUser.Pos.X
		m["y"] = pUser.Pos.Y
		m["angle"] = pUser.Pos.Angle
		m["action"] = pUser.Action
	}

	

	slice = append(slice, m)
	pSocket.SendJsPkg(slice)

	s.SendFriendInvite(uid,pSocket)

	// if s.IsCloseUser(uid,pSocket) {
	// 	return
	// }
	if s.IsUserExist(uid, pSocket) {
		return
	}
	

	s.m_ServerUserList.Store(pNewUser.m_nUserID,pNewUser)
	logger.Error("User Login Succ")

	//向其他用户广播登录
	s.BroadCastOtherLogin(uid)

}

//发送好友邀请信息列表
func (s *GameServer) SendFriendInvite(uid int, pSocket *SocketHandler) {

	m := make(map[string]interface{})
	slice := make([]map[string]interface{}, 0)
	//离线邀请信息
	inviteUserList:=db.Db.GetFriendInvitationInfo(uid)
	m["cmd"] = 0x3005
	m["uid"] = uid
	m["invite_count"] = len(inviteUserList)
	for i := 0; i < len(inviteUserList); i++ {
		m["invite_userid"] = inviteUserList[i]
	}
	slice = append(slice, m)
	pSocket.SendJsPkg(slice)
}

func (s *GameServer) ProcessUserReigster(map_ map[string]interface{}, pSocket *SocketHandler) {
	account, _ := map_["account"].(string)
	password, _ := map_["password"].(string)
	uid := int(db.Db.GetUserId(account, password))

	m := make(map[string]interface{})
	slice := make([]map[string]interface{}, 0)
	ret := 0

	m["cmd"] = 0x3002
	if uid != 0 {
		ret = -1 //该账号已存在
	} else {
		//写入数据库
		if !db.Db.InsertUser(account, password) {
			ret = -2
		}
	}
	m["ret"] = ret
	slice = append(slice, m)
	pSocket.SendJsPkg(slice)

}

func (s *GameServer) ProcessClientHeartBeat(map_ map[string]interface{}, pSocket *SocketHandler) {
	logger.Error("ProcessClientHeartBeat")
	m := make(map[string]interface{})
	slice := make([]map[string]interface{}, 0)

	m["cmd"] = 0x3003
	
	slice = append(slice, m)
	pSocket.SendJsPkg(slice)
}

func (s *GameServer) ProcessChatMsg(map_ map[string]interface{}, pSocket *SocketHandler) {
	logger.Error("ProcessChatMsg")
	m := make(map[string]interface{})
	slice := make([]map[string]interface{}, 0)
	msg:= map_["msg"].(string)

	m["cmd"] = 0x3004
	m["msg"] = msg
	slice = append(slice, m)

	s.BroadCastWorld(slice)
}

func (s *GameServer)Interface2Int(inter interface{})int{
	str := fmt.Sprintf("%v", inter)
	data, err := strconv.Atoi(str)
	if err != nil {
		logger.Error("Interface2Int err")
		return 0
	}
	return data
}

func (s *GameServer)Interface2Float64(inter interface{})float64{
	str := fmt.Sprintf("%v", inter)
	data, err := strconv.ParseFloat(str, 64)
	if err != nil {
		logger.Error("Interface2Int err")
		return 0.00
	}
	return data
}

func (s *GameServer)GetMapIntData(map_ map[string]interface{}, key string)int{
	value_, _ := map_[key]
	value:=s.Interface2Int(value_)
	return value
}
func (s *GameServer)GetMapFloat64Data(map_ map[string]interface{}, key string)float64{
	value_, _ := map_[key]
	value:=s.Interface2Float64(value_)
	return value
}

func (s *GameServer) ProcessChangeUserAction(map_ map[string]interface{}, pSocket *SocketHandler) {

	

	m := make(map[string]interface{})
	slice := make([]map[string]interface{}, 0)

	uid:=s.GetMapIntData(map_,"uid")
	x:=s.GetMapFloat64Data(map_,"x")
	y:=s.GetMapFloat64Data(map_,"y")
	angle:=s.GetMapIntData(map_,"angle")
	action:=s.GetMapIntData(map_,"action")

	logger.Error("ProcessChangeUserAction uid: [%d] action: [%d] angle: [%d] x: [%d] y: [%d]",uid,action,angle,x,y)
	pUser:=s.GetUser(uid)

	if pUser == nil {
		logger.Error("user is not exist")
		return
	}
	pUser.Pos.X = x
	pUser.Pos.Y = y
	pUser.Pos.Angle = angle
	pUser.Action = action

	m["cmd"] = 0x3010
	m["uid"] = pUser.m_nUserID
	m["x"] = pUser.Pos.X
	m["y"] = pUser.Pos.Y
	m["angle"] = pUser.Pos.Angle
	m["action"] = pUser.Action

	slice = append(slice, m)


	//向其他用户广播
	s.BroadCastWorld(slice)



}

func (s *GameServer) BroadCastOtherLogin(uid int) {
	m := make(map[string]interface{})
	slice := make([]map[string]interface{}, 0)

	m["cmd"] = 0x4001
	m["uid"] = uid
	m["strkey"] = "test"
	pUser:= s.GetUser(uid)
	m["x"] = pUser.Pos.X
	m["y"] = pUser.Pos.Y
	m["angle"] = pUser.Pos.Angle
	m["action"] = pUser.Action
	slice = append(slice, m)
	s.BroadCastWorld(slice)
}

func (s *GameServer) BroadCastWorld(slice []map[string]interface{}) {
	userList := s.GetAllUsers()
	for i := 0; i < len(userList); i++ {
		if userList[i].m_pSocket!=nil {
		userList[i].m_pSocket.SendJsPkg(slice)
		logger.Error("BroadCast time: %d", i+1)
		}
	}
}

// 获取用户列表
func (s *GameServer) GetAllUsers() []*CGameUser {
	sl := make([]*CGameUser, 0)
	s.m_ServerUserList.Range(func(k, v interface{}) bool {
		sl = append(sl, v.(*CGameUser))
		return true
	})
	return sl
}

func (s *GameServer) ProcessAddFriendNew(map_ map[string]interface{}, pSocket *SocketHandler) {
	userid:=s.GetMapIntData(map_,"userid")	//发起邀请的用户
	friendid:=s.GetMapIntData(map_,"friendid")	//被邀请的用户

	if userid == friendid {
		return
	}

	db.Db.InsertFriendInvitationInfo(userid,friendid)
	//向用户通知有好友邀请到来	如果不在线添加邀请到数据库	在它上线时从数据库获取所有邀请	在线则立马发送邀请(也要写入数据库)
	friend:=s.GetUser(friendid)

	if friend != nil {
		s.SendFriendInvitation(friend,userid)
	}
}

//发送单个好友邀请信息
func (s *GameServer) SendFriendInvitation(pUser *CGameUser, userid int) {

	//从数据库取出用户信息
	m := make(map[string]interface{})
	slice := make([]map[string]interface{}, 0)
	m["cmd"] = 0x4022
	m["uid"] = userid
	m["name"] = db.Db.GetUserName(userid)
	slice = append(slice, m)
	pUser.m_pSocket.SendJsPkg(slice)
}

func (s *GameServer) ProcessUserRecvFriendInvitation(map_ map[string]interface{}, pSocket *SocketHandler) {
	//taskId:=s.GetMapIntData(map_,"taskId")
	ret:=s.GetMapIntData(map_,"ret")
	userid:=s.GetMapIntData(map_,"userid")
	friendid:=s.GetMapIntData(map_,"friendid")


	//修改邀请信息状态或者删除邀请信息
	//db.Db.DeleteFriendInvitation(taskId)
	if ret == 1 {	//用户同意加好友
		db.Db.InsertNewFriend(userid,friendid, "")
		db.Db.InsertNewFriend(friendid,userid,"")
	}else{		//用户拒绝加好友

	}
}

func (s *GameServer) ProcessUserSearch(map_ map[string]interface{}, pSocket *SocketHandler) {
	info, _ := map_["info"].(string)	//搜索的用户名称
	userList:=db.Db.GetSearchUserList(info)	//模糊查询用户列表

	m := make(map[string]interface{})
	slice := make([]map[string]interface{}, 0)

	m["cmd"] = 0x4023
	m["size"] = len(userList)
	for i:= 0 ;i <len(userList);i++ {
		uid_key:=fmt.Sprintf("uid_%d",i)
		m[uid_key] = userList[i]
		name:=db.Db.GetUserName(userList[i])
		name_key:=fmt.Sprintf("name_%d",i)
		m[name_key] = name
	}
	slice = append(slice, m)
	pSocket.SendJsPkg(slice)
}

func (s *GameServer) ProcessFriendInvitation(map_ map[string]interface{}, pSocket *SocketHandler) {
}

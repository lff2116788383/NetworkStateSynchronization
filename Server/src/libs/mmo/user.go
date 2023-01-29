package mmo

import (
	"Src/libs/config"
	"Src/libs/mmo/data"
	"Src/libs/packet"

	"fmt"
	"net"
	"os"
	"time"
)



//这一块做用户数据处理
type User struct {
	Id       int32
	Tid      int32
	conn     net.Conn
	server  *GameServer     // 缓存Server的引用
	exitTime int64
	lastTime time.Time

	localPath string //用户数据本地保存路径
	// 用户数据模块	使用泛型方便做扩展
	modManage map[string]ModBase

	//Language string
	Money int64
	// 游戏数据


	UserInfo *data.NUserInfo
	UserName string
	PassWord string



}

func NewUser(conn net.Conn, server *GameServer) *User {
	//***************泛型架构***************************
	//var user *User
	//user = new(User)
	//user.Id = userid
	//user.conn = conn
	//
	//user.modManage = map[string]ModBase{
	//			MOD_ROLE: new(ModRole),
	//		}

	user := &User{
		Id:   -1,
		Tid:  -1,
		conn: conn,
		server: server,
		modManage: map[string]ModBase{
			MOD_ROLE: new(ModRole),
			MOD_BAG: new(ModBag),
			MOD_STORE: new(ModStore),


			MOD_CHARACTER:new(ModCharacter),
		},
		UserInfo:new(data.NUserInfo),
	}
	user.UserInfo.InitData()
	return user
}

func (self *User) InitData() {
	path := config.GlobalConfig.LocalSavePath
	_, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return
		}
	}
	self.localPath = path + fmt.Sprintf("/%d", self.Id)
	_, err = os.Stat(self.localPath)
	if err != nil {
		err = os.Mkdir(self.localPath, os.ModePerm)
		if err != nil {
			return
		}
	}


}

func (u *User) Online() {

	// 用户上线，将用户加入到OnlineMap中，注意加锁操作
	u.server.OnlineMap.Store(u.conn.RemoteAddr().String(), u)
	// 广播当前用户上线消息
	//self.server.BroadcastWorld("上线啦O(∩_∩)O")
	fmt.Println("user Online")
}

func (u *User) Offline() {

	// 用户下线，将用户从OnlineMap中删除，注意加锁
	u.server.OnlineMap.Delete(u.conn.RemoteAddr().String())

	// 广播当前用户下线消息
	//this.server.BroadCast(this, "下线了o(╥﹏╥)o")
	fmt.Println("user Offline")
}

func (u *User) InitMod() {
	for _, v := range u.modManage {
		v.LoadData(u)
	}
}

func (u *User) GetType() int {
	return LISTEN_TYPE_CLIENT
}

func (u *User) GetConn() net.Conn {
	return u.conn
}

func (u *User) GenPacket() packet.IPacket {
	return packet.NewICPacket()
}


func (u *User) EnterGameHall()  {
	hall:=u.server.GetTable(MAP_TYPE_HALL)
	if hall != nil {
		hall.AddUser(u)
		u.Tid = MAP_TYPE_HALL
		//进入大厅开始初始化数据模块
		u.InitData()
		u.InitMod()
	}
}

func (u *User) LeaveGameHall()  {
	hall:=u.server.GetTable(MAP_TYPE_HALL)
	if hall != nil {
		hall.DelUser(u.Id)
		u.Tid = -1
	}
}

//泛型数据模块获取
func (u *User) GetModRole() *ModRole {
	return u.modManage[MOD_ROLE].(*ModRole)
}
func (u *User) GetModBag() *ModBag {
	return u.modManage[MOD_BAG].(*ModBag)
}

func (u *User) GetModStore() *ModStore {
	return u.modManage[MOD_STORE].(*ModStore)
}
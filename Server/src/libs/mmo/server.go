package mmo

import (
	. "Src/libs/command"
	"Src/libs/config"
	"Src/libs/csvs"
	"Src/libs/logger"
	"Src/libs/packet"
	"Src/libs/db"
	"net"
	"sync"
	"time"
)

//世界频道 队伍 工会
type GameServer struct {

	DBService *db.BaseDb


	//有多张桌子
	Tables sync.Map // userId *User

	// 在线用户容器
	OnlineMap sync.Map

	//有多张相同的桌子
}

//单例
var gameserver *GameServer

func GetGameServer() *GameServer {
	if gameserver == nil {
		gameserver = new(GameServer)

	}
	return gameserver
}

func (s *GameServer) LoadTables() {
	//table := Table{Id: 1001, Name: "极乐世界",Type: 1}
	//s.AddTable(&table)

	for _,v:=range csvs.ConfigMapMap{
		table := Table{Id: int32(v.MapId), Name: v.MapName,Type: int32(v.MapType)}
		s.AddTable(&table)
	}
}

func (s *GameServer) Init() bool {
	s.DBService = db.GetDB()



	//创建默认房间 房间号1001 名字：极乐世界
	s.LoadTables()
	if len(s.GetAllTables()) == 0 {
		return false
	}
	return true
}

func (s *GameServer) Run() {
	//检测连接断开
	//go s.salir()
	//主线程开启监听 可使用goroutine开启协程并行监听其他端口
	s.Listen(config.GlobalConfig.Host, LISTEN_TYPE_CLIENT)
}

func (s *GameServer)Listen(address string, lType int) {
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
		//case LISTEN_TYPE_ADMIN:
		//	admin := Admin{conn: conn}
		//	go ReceiveIcPkg(&admin)
		case LISTEN_TYPE_CLIENT:
			//创建一个uid默认为-1的用户 暂时不处理
			//有连接到来	直接进入大厅(代表在线人数)

			//这里只是一个连接	不是大厅	登录成功后才进入游戏大厅
			//user := User{Id: -1,Tid:-1,conn: conn,server: s,lastTime: time.Now()}
			user:=NewUser(conn,s)
			user.Online()
			//+s.GetTable(MAP_TYPE_HALL).AddUser(&user)
			// TODO 启动一个协程去处理
			go ReceiveIcPkg(user)

		}

	}
}

// 检测离开
func (s *GameServer) salir() {
	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_TEST_RESP)
	outICP.End()

	for {
		tables := s.GetAllTables()
		for i := 0; i < len(tables); i++ {
			users := tables[i].GetAllUsers()
			for j := 0; j < len(users); j++ {
				u := users[j]
				ret := u.Send(outICP)
				if !ret {
					logger.Error("user:[%d],tid:[%d] tableName:[%s] disconnect", users[j].Id, users[j].Tid,s.GetTable(users[j].Tid).Name)
					table := s.GetTable(u.Tid)
					if table != nil {
						user := table.GetUser(u.Id)
						if user != nil {
							table.DelUser(u.Id)
						}
					}
				}
				time.Sleep(5 * time.Second)
			}
		}
	}

}

func (s *GameServer) KickUser(userid int32) {

	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_KICK_RESP)
	outICP.End()

	tables := s.GetAllTables()
	for i := 0; i < len(tables); i++ {
		users := tables[i].GetAllUsers()
		for j := 0; j < len(users); j++ {
			u := users[j]
			if userid == u.Id {
				u.Send(outICP)
				u.conn.Close()
				tables[i].DelUser(u.Id)
			}

		}
	}

}

// 添加房间
func (s *GameServer) AddTable(t *Table) {

	s.Tables.Store(t.Id, t)

}


//获取单个房间
func (s *GameServer) GetTable(tid int32) *Table {
	v, ok := s.Tables.Load(tid)
	if ok {
		return v.(*Table)
	}
	return nil
}

//获取所有房间
func (s *GameServer) GetAllTables() []*Table {
	sl := make([]*Table, 0)
	s.Tables.Range(func(k, v interface{}) bool {
		sl = append(sl, v.(*Table))
		return true
	})
	return sl
}

//删除房间
func (s *GameServer) DelTable(tid int32) {

	s.Tables.Delete(tid)
}


//获取在线人数
func (s *GameServer) GetOnlineMapCount() int32 {

	count:=0
	s.OnlineMap.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	return int32(count)
}

//广播桌子
func (s *GameServer)BroadcastTable(tid int32, pack *packet.ICPacket) {

	table := s.GetTable(tid)

	if table != nil {

		table.Users.Range(func(k, v interface{}) bool {
			v.(*User).Send(pack)
			return true
		})
	}
}

//大厅广播
func (s *GameServer)BroadcastHall(pack *packet.ICPacket) {
	s.BroadcastTable(MAP_TYPE_HALL,pack)
}

//全局广播
func (s *GameServer)BroadcastWorld(pack *packet.ICPacket) {

	s.Tables.Range(func(k, v interface{}) bool {
		s.BroadcastTable(v.(*Table).Id, pack)
		return true
	})
}


package server

import (
	"Src/libs/config"
	//"Src/libs/db"
	"Src/libs/logger"
	"Src/libs/packet"

	//"fmt"
	"net"
	//"strconv"
	//"sync"
)

type Server struct {
	ClientSocketMap map[string]*SocketHandler
}

//单例
var server *Server

func GetServer() *Server {
	if server == nil {
		server = &Server{
			ClientSocketMap: make(map[string]*SocketHandler),
		}

	}
	return server
}

//成员函数
func (s *Server) Init() bool {
	//s.DBService = db.GetDB()
	return true
}

func (s *Server) Run() {
	//主线程开启监听 可使用goroutine开启协程并行监听其他端口
	s.Listen(config.GlobalConfig.Host, LISTEN_TYPE_CLIENT)

}

func (s *Server) Listen(address string, lType int) {
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

func (s *Server) ProcessClose(pSocket *SocketHandler) {

	// if pSocket.m_pGameUser != nil && pSocket.m_pGameUser.m_pSocket == pSocket {
	// 	pSocket.m_pGameUser.m_pSocket = nil

	// 	s.AddCloseUser(pSocket.m_pGameUser)

	// 	logger.Error("uid:%d ProcessClose ", pSocket.m_pGameUser.m_nUserID)

	// }
	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_USER_LEAVE)
	outICP.WriteString(pSocket.human.Desc)
	outICP.End()

	s.BroadcastPkg(outICP)

	delete(s.ClientSocketMap, pSocket.human.Desc)
	logger.Error("ProcessClose  delete desc:%s rest client count:%d", pSocket.human.Desc, len(s.ClientSocketMap))
}

func (s *Server) BroadcastUserEnterResponse(pSocket *SocketHandler, pack *packet.ICPacket) {
	desc := pack.ReadString()
	x := pack.ReadFloat64()
	y := pack.ReadFloat64()
	z := pack.ReadFloat64()
	ry := pack.ReadFloat64()
	pSocket.human = &BaseHuman{
		X:    x,
		Y:    y,
		Z:    z,
		Ry:   ry,
		Desc: desc,
		HP:   100,
	}
	s.ClientSocketMap[desc] = pSocket
	logger.Error("Receive Enter|socket desc:%s pkg desc:%s,%g,%g,%g,%g", pSocket.conn.RemoteAddr().String(), desc, x, y, z, ry)

	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_USER_ENTER)
	outICP.WriteString(desc)
	outICP.WriteFloat64(x)
	outICP.WriteFloat64(y)
	outICP.WriteFloat64(z)
	outICP.WriteFloat64(ry)
	outICP.End()

	s.BroadcastPkg(outICP)

}

func (s *Server) SendUserList(pSocket *SocketHandler) {

	count := len(s.ClientSocketMap)
	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_USER_LIST)
	outICP.WriteInt32(int32(count))
	for k := range s.ClientSocketMap {
		h := s.ClientSocketMap[k].human
		outICP.WriteString(h.Desc)
		outICP.WriteFloat64(h.X)
		outICP.WriteFloat64(h.Y)
		outICP.WriteFloat64(h.Z)
		outICP.WriteFloat64(h.Ry)
		outICP.WriteInt32(100)
	}
	outICP.End()

	pSocket.Send(outICP)

}

func (s *Server) ProcessUserMove(pack *packet.ICPacket) {
	desc := pack.ReadString()
	//移动的目标位置
	x := pack.ReadFloat64()
	y := pack.ReadFloat64()
	z := pack.ReadFloat64()

	_, ok := s.ClientSocketMap[desc]
	if !ok {
		return
	}

	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_USER_MOVE)
	outICP.WriteString(desc)
	outICP.WriteFloat64(x)
	outICP.WriteFloat64(y)
	outICP.WriteFloat64(z)
	outICP.End()

	s.BroadcastPkg(outICP)
}

func (s *Server) ProcessUserLeave(pSocket *SocketHandler) {
	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_USER_LEAVE)
	outICP.WriteString(pSocket.human.Desc)
	outICP.End()

	s.BroadcastPkg(outICP)

	delete(s.ClientSocketMap, pSocket.human.Desc)
	logger.Error("ProcessUserLeave  delete desc:%s rest client count:%d", pSocket.human.Desc, len(s.ClientSocketMap))
}

func (s *Server) ProcessUserAttack(pack *packet.ICPacket) {
	desc := pack.ReadString()
	eulY := pack.ReadFloat64()

	outICP := packet.NewICPacket()
	outICP.Begin(CMD_SERVER_USER_ATTACK)
	outICP.WriteString(desc)
	outICP.WriteFloat64(eulY)
	outICP.End()

	s.BroadcastPkg(outICP)
}
func (s *Server) ProcessUserHit(pack *packet.ICPacket) {
	attDesc := pack.ReadString()
	hitDesc := pack.ReadString()
	damage := pack.ReadInt32()

	v, ok := s.ClientSocketMap[hitDesc]
	if !ok {
		return
	}

	v.human.HP -= damage

	if v.human.HP <= 0 {
		outICP1 := packet.NewICPacket()
		outICP1.Begin(CMD_SERVER_USER_DIE)
		outICP1.WriteString(attDesc)
		outICP1.WriteString(hitDesc)
		outICP1.WriteInt32(damage)
		outICP1.End()
		s.BroadcastPkg(outICP1)
		logger.Error("desc:%s died", v.human.Desc)
	}
}

func (s *Server) BroadcastPkg(p *packet.ICPacket) {
	for k := range s.ClientSocketMap {
		v := s.ClientSocketMap[k]
		v.Send(p)
	}
}


func (s *Server) ProcessUserPos(pack *packet.ICPacket) {
	desc := pack.ReadString()
	x := pack.ReadFloat64()
	y := pack.ReadFloat64()
	z := pack.ReadFloat64()
	ry := pack.ReadFloat64()

	v, ok := s.ClientSocketMap[desc]
	if !ok {
		return
	}
	v.human.X = x
	v.human.Y = y
	v.human.Z = z
	v.human.Ry = ry
}
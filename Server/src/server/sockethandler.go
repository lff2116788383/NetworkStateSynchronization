package server

import (
	"Src/libs/logger"
	"Src/libs/packet"

	//"encoding/hex"
	"net"
)

type SocketHandler struct {
	conn   net.Conn
	server *Server

	human *BaseHuman
}

func NewSocketHandler(conn net.Conn, server *Server) *SocketHandler {
	pSocket := &SocketHandler{conn: conn,
		server: server}
	return pSocket
}

func (s *SocketHandler) GetType() int {
	return LISTEN_TYPE_CLIENT
}

func (s *SocketHandler) GetConn() net.Conn {
	return s.conn
}

func (s *SocketHandler) GenPacket() packet.IPacket {
	return packet.NewICPacket()
}

//这一块做用户消息处理
func (s *SocketHandler) Send(p *packet.ICPacket) bool {

	_, err := s.conn.Write(p.GetData())
	if err != nil {
		logger.Error("send failed err:%s", err)
		return false
	}
	logger.Error("send cmd:0x%0x", p.GetCmd())
	//logger.Error("send data:[%s]", hex.EncodeToString(p.GetData()))
	return true
}

func (s *SocketHandler) ProcessPkg(p packet.IPacket) {
	//logger.Error("Process Client Pkg")

	pack := p.(*packet.ICPacket)
	cmd := pack.GetCmd()
	logger.Error("recv cmd[0x%x]", cmd)
	//logger.Error("recv data:[%s]", hex.EncodeToString(icp.GetData()))

	if pack.Decrypt() == -1 {
		logger.Error("packet Decrypt failed")
		return
	}

	switch cmd {

	case CMD_CLIENT_USER_ENTER:
		s.server.BroadcastUserEnterResponse(s, pack) //广播玩家上线进入场景
	case CMD_CLIENT_USER_LIST:
		s.server.SendUserList(s) //向该连接发送所有在线玩家数据
	case CMD_CLIENT_USER_MOVE:
		s.server.ProcessUserMove(pack)
	case CMD_CLIENT_USER_LEAVE:
		s.server.ProcessUserLeave(s)
	case CMD_CLIENT_USER_ATTACK:
		s.server.ProcessUserAttack(pack)
	case CMD_CLIENT_USER_HIT:
		s.server.ProcessUserHit(pack)
	case CMD_CLIENT_USER_POS:
		s.server.ProcessUserPos(pack)


	}

	//用户为空 只能接受心跳和登录命令 可以增加登录前的一些相关命令
	// if s.m_pGameUser == nil {
	// 	logger.Error("Socket User is nil, addr: [%s]", s.conn.RemoteAddr().String())
	// } else {
	// 	logger.Error("Socket User is Online, addr: [%s]", s.conn.RemoteAddr().String())
	// }

	// switch cmd {
	// case CLIENT_COMMAND_HEARTBEAT: //心跳
	// 	s.server.ProcessHeartBeat(icp, s)
	// case CLIENT_COMMAND_LOGIN: //登录
	// 	s.server.ProcessUserLogin(icp, s)

	// case CLIENT_COMMAND_SEARCH_USER:
	// 	s.server.ProcessSearchUser(icp, s)
	// case CLIENT_COMMAND_ADD_FRIEND:
	// 	s.server.ProcessAddFriend(icp, s)
	// case 1:
	// 	s.server.ProcessReturnUserTeamInfo(icp, s)
	// default:
	// 	logger.Error("Unknown Cmd 0x%x", cmd)
	// }
}

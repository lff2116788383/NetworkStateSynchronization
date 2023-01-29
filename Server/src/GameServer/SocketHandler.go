package GameServer

import (
	"Src/libs/logger"
	"Src/libs/packet"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"strconv"

)

type SocketHandler struct {
	conn net.Conn

	m_pGameUser *CGameUser
	server      *GameServer
}

func NewSocketHandler(conn net.Conn, server *GameServer) *SocketHandler {
	pSocket := &SocketHandler{conn: conn,
		m_pGameUser: nil,
		server:      server}
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

func (s *SocketHandler) ProcessPkg(p packet.IPacket) {
	logger.Error("Process Client Pkg")

	icp := p.(*packet.ICPacket)
	cmd := icp.GetCmd()
	logger.Error("recv cmd[0x%x]", cmd)
	logger.Error("recv data:[%s]", hex.EncodeToString(icp.GetData()))

	if icp.Decrypt() == -1 {
		logger.Error("packet Decrypt failed")
		return
	}

	//用户为空 只能接受心跳和登录命令 可以增加登录前的一些相关命令
	if s.m_pGameUser == nil {
		logger.Error("Socket User is nil, addr: [%s]", s.conn.RemoteAddr().String())
	} else {
		logger.Error("Socket User is Online, addr: [%s]", s.conn.RemoteAddr().String())
	}

	switch cmd {
	case CLIENT_COMMAND_HEARTBEAT: //心跳
		s.server.ProcessHeartBeat(icp, s)
	case CLIENT_COMMAND_LOGIN: //登录
		s.server.ProcessUserLogin(icp, s)

	case CLIENT_COMMAND_SEARCH_USER:
		s.server.ProcessSearchUser(icp, s)
	case CLIENT_COMMAND_ADD_FRIEND:
		s.server.ProcessAddFriend(icp, s)
	case 1:
		s.server.ProcessReturnUserTeamInfo(icp, s)
	default:
		logger.Error("Unknown Cmd 0x%x", cmd)
	}

	// switch cmd {
	// //登录注册模块
	// case CMD_CLIENT_HEARTBEAT_REQ: //心跳
	// 	u.ProcessHeartBeat(icp)
	// case CMD_CLIENT_LOGIN_REQ: //登录
	// 	u.ProcessLogin(icp)
	// case CMD_CLIENT_REGISTER_REQ: //注册
	// 	u.ProcessRegister(icp)
	// case CMD_CLIENT_LOGOUT_REQ: //登出
	// 	u.ProcessLogout(icp)
	// case CMD_CLIENT_GAME_ENTER_REQ:	//进入游戏
	// 	u.ProcessGameEnter(icp)
	// case CMD_CLIENT_GAME_LEAVE_REQ:	//退出游戏
	// 	u.ProcessGameLeave(icp)

	// //房间地图模块
	// case CMD_CLIENT_GET_TABLES_INFO_REQ: //获取所有房间信息
	// 	u.ProcessGetTablesInfo(icp)
	// case CMD_CLIENT_ENTER_TABLE_REQ: //进入房间
	// 	u.ProcessEnterTable(icp)
	// case CMD_CLIENT_LEAVE_TABLE_REQ:
	// 	u.ProcessLeaveTable(icp)

	// //角色模块
	// case CMD_CLIENT_GET_ROLES_INFO_REQ: //获取角色
	// 	u.ProcessGetRolesInfo(icp)
	// case CMD_CLIENT_CREATE_ROLE_REQ: //创建角色
	// 	u.ProcessCreateRole(icp)
	// case CMD_CLIENT_DEL_ROLE_REQ: //删除角色
	// 	u.ProcessDelRole(icp)

	// //商店模块
	// case CMD_CLIENT_GET_STORE_INFO_REQ:
	// 	u.ProcessGetStoreInfo(icp)
	// case CMD_CLIENT_STORE_BUY_REQ:
	// 	u.ProcessStoreBuy(icp)
	// case CMD_CLIENT_STORE_SELL_REQ:
	// 	u.ProcessStoreSell(icp)

	// //战斗模块
	// case CMD_CLIENT_MOVE_REQ:
	// 	u.ProcessMove(icp)
	// }

	// //
	// //modChoose:=(cmd&0xF000)>>12 //计算cmd的前缀 获取消息所属的mod
	// //switch modChoose {
	// //case 1:
	// //	u.HandleBase(cmd,icp)
	// //case 2:
	// //	self.HandleBag()
	// //case 3:
	// //	self.HandlePool()
	// //case 4:
	// //	self.HandleMap()
	// //case 5:
	// //	self.HandleRelics()
	// //case 6:
	// //	self.HandleRole()
	// //case 7:
	// //	self.HandleWeapon()
	// //case 8:
	// //	for _, v := range self.modManage {
	// //		v.SaveData()
	// 	//}

}

//这一块做用户消息处理
func (s *SocketHandler) Send(p *packet.ICPacket) bool {

	_, err := s.conn.Write(p.GetData())
	if err != nil {
		logger.Error("send failed err:%s", err)
		return false
	}
	logger.Error("send cmd:0x%0x", p.GetCmd())
	logger.Error("send data:[%s]", hex.EncodeToString(p.GetData()))
	return true
}

func (s *SocketHandler) SendJsPkg(slice []map[string]interface{}) bool {

	data, err := json.Marshal(slice)
	if err != nil {
		fmt.Printf("序列化错误 err=%v\n", err)
	}
	end := "****"
	data = append(data, []byte(end)...)
	_, err1 := s.conn.Write(data)
	if err1 != nil {
		logger.Error("send failed err:%s", err)
		return false
	}
	logger.Error("send data:[%s]", hex.EncodeToString(data))
	return true
}

func (s *SocketHandler) ProcessJsPkg(p map[string]interface{}) {

	key, ok := p["cmd"] // 可以用int64
	if !ok {
		logger.Error("断言失败， key：cmd 不是float64类型")
		return
	} else {
		//logger.Error(" key：cmd 是float64类型, 并且值为：%f", key)
	}

	str := fmt.Sprintf("%v", key)

	cmd, err := strconv.Atoi(str)
	if err != nil {
		logger.Error("cmd err")
		return
	}
	//cmd:=int16(key)
	logger.Error("recv cmd[0x%x]", cmd)

	//var slice []map[string]interface{}
	//slice = append(slice, p)
	//s.SendJsPkg(slice)

	switch cmd {
	case 0x2001:
		s.server.ProcessUserLoginNew(p, s)
	case 0x2002:
		s.server.ProcessUserReigster(p, s)
	case 0x2003:
		s.server.ProcessClientHeartBeat(p, s)
	case 0x2004:
		s.server.ProcessChatMsg(p, s)
	case 0x2010:
		s.server.ProcessChangeUserAction(p, s)
	
	//添加好友相关
	case 0x2020:
		s.server.ProcessAddFriendNew(p, s)
	case 0x2021:
		s.server.ProcessFriendInvitation(p, s)
	case 0x2022:
		s.server.ProcessUserRecvFriendInvitation(p, s)
	case 0x2023:
		s.server.ProcessUserSearch(p, s)

	default:

	}
}

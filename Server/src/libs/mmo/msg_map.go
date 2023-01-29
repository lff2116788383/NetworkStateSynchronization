package mmo
//
//import (
//	"Src/libs/csvs"
//	"Src/libs/logger"
//	"Src/libs/packet"
//	"time"
//)
//
//const (
//	CMD_CLIENT_ENTER_MAP_REQ = 0x4001
//	CMD_CLIENT_LEAVE_MAP_REQ = 0x4002
//	CMD_CLIENT_GET_MAPS_REQ  = 0x4003
//
//)
//
//
////地图消息模块
//func (u *User) HandleMap(cmd uint16,icp *packet.ICPacket) {
//	switch cmd {
//	case CMD_CLIENT_HEARTBEAT_REQ: //心跳
//		u.ProcessHeartBeat(icp)
//	case CMD_CLIENT_LOGIN_REQ: //登录
//		u.ProcessLogin(icp)
//	case CMD_CLIENT_REGISTER_REQ: //注册
//		u.ProcessRegister(icp)
//	case CMD_CLIENT_LOGOUT_REQ: //登出
//		u.ProcessLogout(icp)
//	}
//}
//
//
//
//
////进入地图
//func (u *User) HandleEnterMap(pack *packet.ICPacket) {
//	//TODO
//	uId:=pack.ReadInt32()
//	roleId:=pack.ReadInt32()
//	mapId:=pack.ReadInt32()
//
//	//role_id->GameObject map_id->LoadScene
//	ret:=0
//
//
//	configRole:=csvs.ConfigRoleMap[int(roleId)]
//	if configRole == nil {
//		logger.Error("无法识别的角色")
//		ret=1
//		return
//	}
//	configMap:= csvs.ConfigMapMap[int(mapId)]
//	if configMap == nil {
//		logger.Error("无法识别的地图")
//		ret=2
//		return
//	}
//
//	//用户是否已经在地图上
//
//
//	outICP := packet.NewICPacket()
//	outICP.Begin(CMD_SERVER_HEARTBEAT_RESP)
//	outICP.WriteByte(byte(ret))
//	outICP.WriteInt32(uId)
//	outICP.WriteInt32(roleId)
//	outICP.WriteInt32(mapId)
//	outICP.End()
//	u.Send(outICP)
//
//	table := GetGameServer().GetTable(u.tid)
//		if table != nil {
//			table.Broadcast(outICP)
//		}
//}
//
////离开地图
//func (u *User) HandleLeaveMap(pack *packet.ICPacket) {
//	//TODO
//	logger.Error("MsgPing")
//	u.lastTime = time.Now() //记录心跳时间
//
//	outICP := packet.NewICPacket()
//	outICP.Begin(CMD_SERVER_HEARTBEAT_RESP)
//	outICP.End()
//	u.Send(outICP)
//}
//
////获取所有地图
//func (u *User) HandleGetMaps(pack *packet.ICPacket) {
//	//TODO
//	logger.Error("MsgPing")
//	u.lastTime = time.Now() //记录心跳时间
//
//	outICP := packet.NewICPacket()
//	outICP.Begin(CMD_SERVER_HEARTBEAT_RESP)
//	outICP.End()
//	u.Send(outICP)
//}
package mmo
//
//import (
//	"Src/libs/db"
//	"Src/libs/logger"
//	"Src/libs/packet"
//	"time"
//)
//
////
////import (
////	"Src/libs/db"
////	"Src/libs/logger"
////	"Src/libs/packet"
////	"time"
////)
////
////CLIENT BASE CMD
//const (
//	CMD_CLIENT_HEARTBEAT_REQ = 0x3000 //心跳请求
//	CMD_CLIENT_LOGIN_REQ     = 0x3001 //登录请求
//	CMD_CLIENT_REGISTER_REQ  = 0x3002 //注册请求
//	CMD_CLIENT_LOGOUT_REQ    = 0x3003 //登出请求
//
//	//CMD_CLIENT_CREATE_ROLE_REQ     = 0x3004 //创建角色请求
//	//CMD_CLIENT_GET_TABLES_INFO_REQ = 0x3005 // 获取所有房间信息请求
//	//CMD_CLIENT_ENTER_TABLE_REQ     = 0x3006 // 进入房间请求
//	//CMD_CLIENT_LEAVE_TABLE_REQ     = 0x3007 // 离开房间请求
//)
//
////SERVER CMD
//const (
//	CMD_SERVER_TEST_RESP = 0x3999 //测试响应
//
//	CMD_SERVER_HEARTBEAT_RESP = 0x4000 //心跳响应
//	CMD_SERVER_LOGIN_RESP     = 0x4001 //登录结果响应
//	CMD_SERVER_REGISTER_RESP  = 0x4002 //注册结果响应
//	CMD_SERVER_LOGOUT_RESP    = 0x4003 //登出结果响应
//
//	//CMD_SERVER_CREATE_ROLE_RESP     = 0x4004 //创建角色结果响应
//	//CMD_SERVER_GET_TABLES_INFO_RESP = 0x4005 //获取所有房间信息结果响应
//	//CMD_SERVER_ENTER_TABLE_RESP     = 0x4006 //进入房间结果响应
//	//
//	//CMD_SERVER_MOVE_RESP = 0x4010 //移动结果响应
//
//)
////0x1000 	GM命令
////0x2000	备用命令
////0x3000	基本命令(包括一系列的登录注册心跳等)
//func (u *User) HandleBase(cmd uint16,icp *packet.ICPacket) {
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
//func (u *User) ProcessHeartBeat(pack *packet.ICPacket) {
//	//TODO
//	logger.Error("MsgPing")
//	u.lastTime = time.Now() //记录心跳时间
//
//	outICP := packet.NewICPacket()
//	outICP.Begin(CMD_SERVER_HEARTBEAT_RESP)
//	outICP.End()
//	u.Send(outICP)
//}
//func (u *User) ProcessLogin(pack *packet.ICPacket) {
//	//TODO
//	//登录 客户端会发送账号和密码 由数据库进行验证 返回唯一的uid
//	account := pack.ReadString()
//	password := pack.ReadString()
//
//	ret := 0
//	if db.GetDB().Mysqldb == nil {
//		db.GetDB().InitMysql()
//	}
//	uid := int32(db.GetDB().GetUserId(account, password))
//
//	if uid == 0 {
//		//登录错误 返回错误码
//		ret = 1
//	}
//
//	//登录成功 发送登录结果 uid 玩家列表
//	u.Id = uid
//	outICP := packet.NewICPacket()
//	outICP.Begin(CMD_SERVER_LOGIN_RESP)
//	outICP.WriteByte(byte(ret))
//	outICP.WriteInt32(uid)
//	outICP.End()
//	u.Send(outICP)
//
//	if ret == 0 {
//		logger.Error("client login succ account:[%s],password:[%s]", account, password)
//
//	} else {
//		logger.Error("client login fail account:[%s],password:[%s],ret:[%d]", account, password, ret)
//	}
//
//}
//
//func (u *User) ProcessRegister(pack *packet.ICPacket) {
//	//TODO
//	//注册 客户端会发送 邮箱 账号和密码 由数据库进行验证
//	account := pack.ReadString()
//	password := pack.ReadString()
//
//	ret := 0
//	//如果该用户已存在 返回错误 否则注册成功
//	uid := int32(db.GetDB().GetUserId(account, password))
//	if uid != 0 {
//		//注册失败 用户已存在
//		ret = 1
//	} else {
//		if !db.GetDB().InsertUser(account, password) {
//			//注册失败 数据库添加失败
//			ret = 2
//		}
//	}
//
//	outICP := packet.NewICPacket()
//	outICP.Begin(CMD_SERVER_REGISTER_RESP)
//	outICP.WriteByte(byte(ret))
//	outICP.End()
//	u.Send(outICP)
//
//	if ret == 0 {
//		logger.Error("client register succ account:[%s],password:[%s]", account, password)
//
//	} else {
//		logger.Error("client register fail account:[%s],password:[%s],ret:[%d]", account, password, ret)
//	}
//}
//func (u *User) ProcessLogout(pack *packet.ICPacket) {
//}
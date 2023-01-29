package server

//Client-GameServer 命令
const (
	CMD_CLIENT_USER_ENTER = 0x4001
	CMD_CLIENT_USER_LIST = 0x4002
	CMD_CLIENT_USER_MOVE = 0x4003
	CMD_CLIENT_USER_LEAVE = 0x4004
	CMD_CLIENT_USER_ATTACK = 0x4005
	CMD_CLIENT_USER_HIT = 0x4006
	CMD_CLIENT_USER_DIE= 0x4007
	CMD_CLIENT_USER_POS= 0x4008
)

//GameServer-Client 命令
const (
	CMD_SERVER_USER_ENTER = 0x5001
    CMD_SERVER_USER_LIST = 0x5002
    CMD_SERVER_USER_MOVE = 0x5003
    CMD_SERVER_USER_LEAVE = 0x5004
    CMD_SERVER_USER_ATTACK = 0x5005
	CMD_SERVER_USER_HIT = 0x5006
	CMD_SERVER_USER_DIE= 0x5007
	CMD_SERVER_USER_POS= 0x5008
)

package command

//AC-GameServer 命令字
const (
	CMD_AS_USER_DATA   = 0x1000
	CMD_SA_REGISTER    = 0x1001
	CMD_SA_HEARTBEAT   = 0x1002
	CMD_SA_USER_LOGOUT = 0x1003
	CMD_PS_ADMIN       = 0x1004
)

//AC-Client 命令字
const (
	CMD_CA_AUTH = 0x2001
)

//Client-GameServer 命令字
const (
	CMD_CS_HEARTBEAT = 0x3000 //心跳
	CMD_CS_LOGIN     = 0x3001 //登录
	CMD_CS_REGISTER  = 0x3002 //注册
	//CMD_CS_SPIN      = 0x3002 //摇奖
	//CMD_CS_GET_DATA  = 0x3003 //请求房间数据
	//CMD_CS_JACKPOT   = 0x3004 //请求Jackpot数据
	//CMD_CS_LIST_ADD  = 0x3005 //用户列表上座
	//CMD_CS_LIST_DEL  = 0x3006 //用户列表下座
	//CMD_CS_LIST_WIN  = 0x3007 //用户列表中奖
	//CMD_CS_LOGOUT    = 0x300A //退出
	//CMD_CS_ASSETS  	 = 0x3008 // 资产信息

	//睡美人
	CMD_CS_SB_SMALL_GAME_INIT = 0x3012 //睡美人请求初始化小游戏
	CMD_CS_SB_SMALL_GAME      = 0x3011 //睡美人点击图标

	//财神爷
	CMD_CS_WG_SMALL_GAME_INIT = 0x3013 //初始化聚宝盆
	CMD_CS_WG_SMALL_GAME      = 0x3014 //打开聚宝盆
	CMD_CS_WG_CRAZY_GAME_INIT = 0x3015 //初始化狂热模式
	CMD_CS_WG_CRAZY_GAME      = 0x3016 //狂热模式摇奖

	//凤凰
	CMD_CS_PH_SMALL_GAME_INIT = 0x3017 //凤凰请求初始化小游戏
	CMD_CS_PH_SMALL_GAME      = 0x3018 //凤凰点击图标

	//水果派对
	CMD_CS_FP_SMALL_GAME_INIT = 0x3019 //水果派对请求初始化小游戏
	CMD_CS_FP_SMALL_GAME      = 0x3020 //水果派对点击图标

	//万圣节
	CMD_CS_HL_SMALL_GAME_INIT = 0x3021 //万圣节请求初始化小游戏
	CMD_CS_HL_SMALL_GAME      = 0x3022 //万圣节点击图标

	//亡灵节
	CMD_CS_DD_SMALL_GAME_INIT = 0x3023 //亡灵节请求初始化小游戏
	CMD_CS_DD_SMALL_GAME      = 0x3024 //亡灵节玩小游戏

	//金鸡报喜
	CMD_CS_GC_SMALL_GAME_INIT = 0x3025 //金鸡报喜请求初始化小游戏
	CMD_CS_GC_SMALL_GAME      = 0x3026 //金鸡报喜玩小游戏

	//公路之王
	CMD_CS_RK_ICONS = 0x3027 //请求公路之王上面的图标

	//小游戏广播
	CMD_CS_GC_SMALL_GAME_BROADCAST = 0x3028 //金鸡小游戏奖励广播
)

const (
	CMD_CS_BACCARAT_LOGIN       = 0x3001 //登录
	CMD_CS_BACCARAT_USER_LOGOUT = 0x4002 //退出
	CMD_CS_BACCARAT_CHIPIN      = 0x4003 //下注
	CMD_CS_BACCARAT_TABLE_INFO  = 0x4004 //获取桌子信息
	CMD_CS_BACCARAT_USER_LIST   = 0x4005 //获取前6个用户信息
	CMD_CS_BACCARAT_USER_LIST2  = 0x4006 //获取所有用户信息

	CMD_CS_BACCARAT_BROADCAST_GAME_START  = 0x5001 //开始下注
	CMD_CS_BACCARAT_BROADCAST_GAME_SETTLE = 0x5002 //结算信息
	CMD_CS_BACCARAT_USER_GAME_SETTLE      = 0x5003 //结算信息
)

//GameServer 命令字
const (
	CMD_S_STATUS = 0x4001 //当前状态
)

//系统命令字
const (
	CMD_ADMIN_RESET_CONFIG    = 0x3002 //重置server 重新读取mysql、config.ini
	CMD_ADMIN_RESET_CONFIG_RS = 0x3136 //PHP重置返回状态

	CMD_ADMIN_GET_SERVER_INFO    = 0x3124 //CMS获取Server信息
	CMD_ADMIN_GET_SERVER_INFO_RS = 0x4083 //回复Server信息给PHP

	CMD_ADMIN_STOP_SERVER = 0x3004 //停服
	CMD_ADMIN_KICK_USER   = 0x3005 //踢人
)

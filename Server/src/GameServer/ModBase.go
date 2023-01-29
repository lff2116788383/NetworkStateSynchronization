package GameServer


//每个模块实现增删改查的基本接口
const (
	MOD_USER     = "user"	//用户模块	包含用户重要的唯一信息 如游戏币 生日 注册时间 邮箱等

	MOD_CHARACTER  =  "character"//新版角色模块
	MOD_ROLE       = "role"		//角色模块	创建选择角色进入地图	切换控制角色
	MOD_BAG        = "bag"		//背包模块
	MOD_STORE	   = "store"	//商店模块
	MOD_FRIEND        = "friend"		//好友模块

	MOD_WEAPON     = "weapon "	//武器模块
	MOD_PET		   = "pet"		//宠物模块
	MOD_GUILD 	   = "guild"	//公会模块

	//玩家休闲功能模块	例如钓鱼模块 摆摊


	//server绑定模块	世界地图模块	如刷新某个地图的唯一野怪
	MOD_SYS_WORLD_MAP = "world_map"
	MOD_SYS_WORLD_STORE = "world_store"
	SYS_MOD_MATCH = "sys_match"



)


//用户模块
type ModBase interface {
	LoadData(pUser *CGameUser) //通过用户加载初始化模块
	SaveData()
	InitData()

}


//系统模块
type SysModBase interface {
	LoadData(server *GameServer) //通过server加载初始化模块
	SaveData()
	InitData()

}

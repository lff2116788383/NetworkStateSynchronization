package mmo

//每个模块实现增删改查的基本接口
const (
	MOD_PLAYER     = "player"

	MOD_CHARACTER  =  "character"//新版角色模块


	MOD_ROLE       = "role"		//角色模块	创建选择角色进入地图	切换控制角色
	MOD_MAP        = "map"		//地图模块	地图事件刷新
	MOD_BAG        = "bag"		//背包模块

	MOD_STORE	   = "store"	//商店模块

	MOD_WEAPON     = "weapon "	//武器模块
	MOD_PET		   = "pet"		//宠物模块
	MOD_GUILD 	   = "guild"	//公会模块



)

type ModBase interface {
	LoadData(User *User) //通过用户加载初始化模块
	SaveData()
	InitData()

	//AddItem()
	//RemoveItem()
	//UpdateItem()
	//SelectItem()
}
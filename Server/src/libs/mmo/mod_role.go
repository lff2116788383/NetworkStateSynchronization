package mmo

import (
	"Src/libs/csvs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const(
	None = 0
	Warrior = 1
	Wizard = 2
	Archer = 3
)


//数据池
type DataPool struct {
	Current int64
	Max int64
}

func (d *DataPool ) Limit()  {
	if d.Current > d.Max {
		d.Current = d.Max
	}
}


type RoleInfo struct {
	RoleId    int32			//角色ID
	RoleName  string        //角色名字
	RoleType  int32 		//角色类型
	RoleLevel int32 		//角色等级
	RoleNum   int64			//角色数量

	//GetTimes   int32 		//获取时间
	//RelicsInfo []int32	//圣遗物信息
	//WeaponInfo int32		//武器信息


	//角色数据配置	excel表填充
	Hp DataPool
	Mp DataPool
	Sp DataPool
	Exp DataPool
	Level DataPool
}

//角色模块与用户绑定 角色模块包含多个角色信息
type ModRole struct {
	RoleInfo  map[int32]*RoleInfo
	//HpPool    int32
	//HpCalTime int64
    CtrlRole *RoleInfo
	user *User
	path string
}

//添加角色信息

func (self *ModRole) AddItem(roleId int32, num int64) {
	//判读配置是否有该角色表
	config := csvs.GetRoleConfig(int(roleId))
	if config == nil {
		fmt.Println("配置不存在roleId:", roleId)
		return
	}

	//添加num个RoleInfo
	for i := 0; i < int(num); i++ {
		_, ok := self.RoleInfo[roleId]
		if !ok {//角色不存在
			data := new(RoleInfo)
			data.RoleId = roleId
			data.RoleType =int32(config.Type)
			data.RoleLevel = 1
			data.RoleNum = num
			self.RoleInfo[roleId] = data
			break
		}else {
			self.RoleInfo[roleId].RoleNum+=num
		}
	}
	itemConfig := csvs.GetItemConfig(int(roleId))
	if itemConfig != nil {
		fmt.Println("获得角色", itemConfig.ItemName, "ID:", roleId, "------现有数量:", self.RoleInfo[roleId].RoleNum)
	}

}

//删除角色信息
func (self *ModRole) RemoveItem(roleId int32, num int64) {
	//判读配置是否有该角色表
	config := csvs.GetRoleConfig(int(roleId))
	if config == nil {
		fmt.Println("配置不存在roleId:", roleId)
		return
	}

	//添加num个RoleInfo

		_, ok := self.RoleInfo[roleId]
		if !ok {//不存在
			return
		}else {
			if self.RoleInfo[roleId].RoleNum  < num {
				return
			}
			self.RoleInfo[roleId].RoleNum-=num
		}

	itemConfig := csvs.GetItemConfig(int(roleId))
	if itemConfig != nil {
		fmt.Println("删除角色", itemConfig.ItemName, "ID:", roleId, "------现有数量:", self.RoleInfo[roleId].RoleNum)
	}

}

//获取角色列表
func (self *ModRole) GetRoleList(roleId int32) []*RoleInfo {
	sl := make([]*RoleInfo, 0)

	for _,v:=range self.RoleInfo{
		sl = append(sl, v)
	}
	return sl
}

//获取单个角色
func (self *ModRole) GetRole(roleId int32) *RoleInfo {
	for k,v:=range self.RoleInfo{
		if k == roleId {
			return v
		}
	}
	return nil
}

//是否有该角色信息
func (self *ModRole) IsHasRole(roleId int32) bool {
	//遍历角色模块是否有角色信息
	for _, v := range self.RoleInfo {
		if v.RoleId == roleId {
			return true
		}
	}
	return false
}
//获取角色等级
func (self *ModRole) GetRoleLevel(roleId int32) int32 {
	for _, v := range self.RoleInfo {
		if v.RoleId == roleId {
			return v.RoleLevel
		}
	}
	return 0
}






//--------------------------Mod接口方法
//path路径保存角色数据
func (self *ModRole) SaveData() {
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

//localPath加载角色数据
func (self *ModRole) LoadData(user *User) {

	self.user = user
	self.path = self.user.localPath + "/role.json"

	configFile, err := ioutil.ReadFile(self.path)
	if err != nil {
		self.InitData()
		return
	}
	err = json.Unmarshal(configFile, &self)
	if err != nil {
		self.InitData()
		return
	}

	if self.RoleInfo == nil {
		self.RoleInfo = make(map[int32]*RoleInfo)
	}

	if self.CtrlRole == nil {
		self.CtrlRole = new(RoleInfo)
	}
	return
}

//加载中进行初始化数据
func (self *ModRole) InitData() {
	if self.RoleInfo == nil {
		self.RoleInfo = make(map[int32]*RoleInfo)
	}

	if self.CtrlRole == nil {
		self.CtrlRole = new(RoleInfo)
	}
}

package mmo

import (
	"Src/libs/excels"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ModProp struct {
	ID int32
	Name string

	user *User			//拥有者
	path   string
}

func (self *ModProp) UseItem(propId int32, num int64) {
	//判读配置是否有该角色表
	config := excels.GetPropConfig(int(propId))
	if config == nil {
		fmt.Println("配置不存在propId:", propId)
		return
	}

	//判断背包是否存在足够数量的Prop
	if !self.user.GetModBag().HasEnoughItem(propId,num){
		return
	}
	
	//判断prop类型及使用效果
	switch config.Type {
	case "药品":
		switch	propId {
		case 1:
			self.user.GetModRole().CtrlRole.Hp.Current+=500
		case 2:
			self.user.GetModRole().CtrlRole.Mp.Current+=500
		}

		if self.user.GetModRole().CtrlRole.Hp.Current>self.user.GetModRole().CtrlRole.Hp.Max {
			self.user.GetModRole().CtrlRole.Hp.Current = self.user.GetModRole().CtrlRole.Hp.Max
		}
		if self.user.GetModRole().CtrlRole.Mp.Current>self.user.GetModRole().CtrlRole.Mp.Max {
			self.user.GetModRole().CtrlRole.Mp.Current = self.user.GetModRole().CtrlRole.Mp.Max
		}
		
	}
}









//--------------------------Mod接口方法
//path路径保存角色数据
func (self *ModProp) SaveData() {
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
func (self *ModProp) LoadData(user *User) {

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

	//if self.RoleInfo == nil {
	//	self.RoleInfo = make(map[int32]*RoleInfo)
	//}
	return
}

//加载中进行初始化数据
func (self *ModProp) InitData() {
	//if self.RoleInfo == nil {
	//	self.RoleInfo = make(map[int32]*RoleInfo)
	//}
}
package mmo

import (
	. "Src/libs/mmo/data"
	"encoding/json"
	"io/ioutil"
	"os"
)

//const (
//	Character_Type_Player   = 0
//	Character_Type_Npc  	= 1
//	Character_Type_Monster  = 2
//)




type ModCharacter struct {
	Characters *NCharacterInfo

	user *User			//拥有者
	path   string
}



//--------------------------Mod接口方法
//path路径保存角色数据
func (self *ModCharacter) SaveData() {
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
func (self *ModCharacter) LoadData(user *User) {

	self.user = user
	self.path = self.user.localPath + "/character.json"

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


	return
}

//加载中进行初始化数据
func (self *ModCharacter) InitData() {
	//if self.Characters == nil {
	//	self.Characters = make(map[int]*TCharacter)
	//}
}
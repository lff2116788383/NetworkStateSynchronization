package GameServer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type PetInfo struct {
	Id int
	name string
}
//宠物模块与用户绑定	玩家可以指定宠物的使用角色
type ModPet struct {

	PetInfo  map[int]*PetInfo
	CtrlRole *RoleInfo
	user *CGameUser
	path string
}



//--------------------------Mod接口方法
func (self *ModPet) SaveData() {
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (self *ModPet) LoadData(user *CGameUser) {

	self.user = user
	self.path = self.user.localPath + "/pet.json"

	configFile, err := ioutil.ReadFile(self.path)
	if err != nil {
		fmt.Println("error")
		return
	}
	err = json.Unmarshal(configFile, &self)
	if err != nil {
		self.InitData()
		return
	}

	if self.PetInfo == nil {
		self.PetInfo = make(map[int]*PetInfo)
	}
	return
}

func (self *ModPet) InitData() {

}
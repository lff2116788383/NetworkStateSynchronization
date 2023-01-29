package server

import (
	. "Src/libs/mmo"
	"net"
	"sync"
	. "Src/libs/mmo/data"
)

type Character struct {
	
}
type MapCharacter struct {
	conn	*net.Conn
	character	*Character
}



type Map struct {
	Id int
	Define *MapDefine
	MapCharacters map[int]*MapCharacter
	InstanceId int
}



type ModMap struct {
	Maps sync.Map // userId *User

	user *User
}



func (this *ModMap) AddMap(m *Map) {

	this.Maps.Store(m.Id, m)

}


func (this *ModMap) GetMap(mapId int) *Map {
	v, ok := this.Maps.Load(mapId)
	if ok {
		return v.(*Map)
	}
	return nil
}


func (this *ModMap) GetMaps() []*Map {
	sl := make([]*Map, 0)
	this.Maps.Range(func(k, v interface{}) bool {
		sl = append(sl, v.(*Map))
		return true
	})
	return sl
}


func (this *ModMap) DelMap(mapId int) {
	this.Maps.Delete(mapId)
}




//--------------------------Mod接口方法
//path路径保存角色数据
func (this *ModMap) SaveData() {

}

//localPath加载角色数据
func (this *ModMap) LoadData(user *User) {
	this.user = user

	return
}

//加载中进行初始化数据
func (this *ModMap) InitData() {

}

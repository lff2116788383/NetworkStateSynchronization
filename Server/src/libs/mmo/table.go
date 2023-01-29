package mmo

import (
	"Src/libs/logger"
	"Src/libs/packet"
	"sync"
)
//桌子地图类型
const(
	MAP_TYPE_ONLINE = 0
	MAP_TYPE_HALL = 1

	MAP_TYPE_CLASSIC = 1001  //经典地图	不可删除
	MAP_TYPE_CUSTOM  = 1002  //自定义地图	可删除
)


//修改房间配置进行初始化

type Table struct {
	Id   int32
	Name string
	Type int32
	//done  chan struct{}
	Users sync.Map // userId *User

	Master *User	//房主
	User_Max int32	//最大玩家数
	Camps	sync.Map	//阵营	key

}

//根据配置表new不同的table
func NewTable(tid int32,name string,tabletype int32)* Table  {
	table:=&Table{Id: tid,Name: name,Type: tabletype}
	return table
}


// 添加用户
func (t *Table) AddUser(u *User) {

	t.Users.Store(u.Id, u)
}

//获取单个用户
func (t *Table) GetUser(userId int32) *User {
	v, ok := t.Users.Load(userId)
	if ok {
		return v.(*User)
	}
	return nil
}

//删除用户
func (t *Table) DelUser(userId int32) {

	t.Users.Delete(userId)
}

// 获取用户列表
func (t *Table) GetAllUsers() []*User {
	sl := make([]*User, 0)
	t.Users.Range(func(k, v interface{}) bool {
		sl = append(sl, v.(*User))
		return true
	})
	return sl
}

//func (t *Table) tableInfo() {
//	var userCnt int
//	t.Users.Range(func(k, v interface{}) bool {
//		userCnt++
//		return true
//	})
//	//logger.Info(fmt.Sprintf("table info user count %d", userCnt), t.logField())
//}


//添加阵营
func (t *Table) AddCamp(u *Camp) {

	t.Users.Store(u.Id, u)
}

//获取阵营
func (t *Table) GetCamp(campId int32) *Camp {
	v, ok := t.Users.Load(campId)
	if ok {
		return v.(*Camp)
	}
	return nil
}

//删除阵营
func (t *Table) DelCamp(campId int32) {

	t.Users.Delete(campId)
}

// 获取阵营列表
func (t *Table) GetAllCamps() []*Camp {
	sl := make([]*Camp, 0)
	t.Users.Range(func(k, v interface{}) bool {
		sl = append(sl, v.(*Camp))
		return true
	})
	return sl
}

//广播桌子所有用户
func (t *Table) Broadcast(p *packet.ICPacket) {

	t.Users.Range(func(k, v interface{}) bool {
		logger.Error("Table Broadcast tid:[%d] uid[%d]", t.Id, v.(*User).Id)
		v.(*User).Send(p)
		return true
	})
	//logger.Info(fmt.Sprintf("table info user count %d", userCnt), t.logField())
}

//阵营广播
func (t *Table) BroadcastCamp(campId int32,p *packet.ICPacket) {
	camp:=t.GetCamp(campId)

	if camp != nil {
		camp.Users.Range(func(k, v interface{}) bool {
			logger.Error("Table Broadcast tid:[%d] uid[%d]", t.Id, v.(*User).Id)
			v.(*User).Send(p)
			return true
		})
	}
}
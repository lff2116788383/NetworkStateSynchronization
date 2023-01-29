package mmo

import "sync"

//红蓝阵营 旁观
const(
	CAMP_TYPE_RED  = 1
	CAMP_TYPE_BLUE = 2
	CAMP_TYPE_LOOK = 9
)

//阵营
type Camp struct {
	Id   int32
	Users sync.Map // userId *User
}



// 添加用户
func (t *Camp) AddUser(u *User) {

	t.Users.Store(u.Id, u)
}

//获取单个用户
func (t *Camp) GetUser(userId int32) *User {
	v, ok := t.Users.Load(userId)
	if ok {
		return v.(*User)
	}
	return nil
}

//删除用户
func (t *Camp) DelUser(userId int32) {

	t.Users.Delete(userId)
}

// 获取用户列表
func (t *Camp) GetAllUsers() []*User {
	sl := make([]*User, 0)
	t.Users.Range(func(k, v interface{}) bool {
		sl = append(sl, v.(*User))
		return true
	})
	return sl
}
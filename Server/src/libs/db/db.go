package db

import (
	"Src/libs/config"
	"Src/libs/logger"
	"fmt"
	"strconv"
)

var Db = NewDb()

type BaseDb struct {
	Mysqldb *MysqlClient
}



func GetDB() *BaseDb {
	return Db
}

func NewDb() *BaseDb {
	return &BaseDb{}
}

func (db *BaseDb) InitMysql() bool {

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.GlobalConfig.DBConfig.Mysql.User,
		config.GlobalConfig.DBConfig.Mysql.Password,
		config.GlobalConfig.DBConfig.Mysql.IP,
		config.GlobalConfig.DBConfig.Mysql.Port,
		config.GlobalConfig.DBConfig.Mysql.Databasename)
	logger.Debug("dataSource:%s", dataSource)
	db.Mysqldb = MysqlConnect(dataSource)
	if nil == db.Mysqldb {
		logger.Error("Init Mysql failed. host[%s] port[%d]",
			config.GlobalConfig.DBConfig.Mysql.IP, config.GlobalConfig.DBConfig.Mysql.Port)
		return false
	}
	return true
}

func (db *BaseDb) GetUserInfo(userid int) []map[string]string {
	querySql := fmt.Sprintf("SELECT * FROM flamingo.t_user where f_id =%d", userid)
	return db.Mysqldb.MyQuery(querySql)
}

func (db *BaseDb) GetUserId(name string, password string) int64 {
	querySql := fmt.Sprintf("SELECT uid FROM shadow.User where name ='%s' and password ='%s'", name, password)
	//querySql := "select uid from shadow.User where name ='" + name + "' and password = '" + password + "'"
	fmt.Println("querySql:", querySql)
	//if Db.Mysqldb == nil {
	//	Db.InitMysql()
	//}
	queryRet := db.Mysqldb.MyQuery(querySql)
	if queryRet != nil {
		mapQuery := queryRet[0]

		fmt.Println("mapQuery:", mapQuery)
		value, ok := mapQuery["uid"]
		if ok {
			uid, _ := strconv.ParseInt(value, 10, 64)
			return uid
		}

		fmt.Println("mapQuery key uid doesn't exist.")
		return 0
	}

	fmt.Println("queryRet is nil")
	return 0

}

func (db *BaseDb) InsertUser(name string, password string) bool {
	querySql := fmt.Sprintf("insert into shadow.User(name,password) values('%s','%s')", name, password)
	if !db.Mysqldb.MyExec(querySql) {
		fmt.Println("Update User State fail")
		return false
	}
	return true
}

func (db *BaseDb) GetRoleId(uid int32, role_id int32, name string) int64 {

	querySql := fmt.Sprintf("SELECT id FROM shadow.role where user_id = %d and role_type =%d and role_name ='%s'", uid, role_id, name)
	queryRet := db.Mysqldb.MyQuery(querySql)
	if queryRet != nil {
		mapQuery := queryRet[0]

		fmt.Println("mapQuery:", mapQuery)
		value, ok := mapQuery["id"]
		if ok {
			id, _ := strconv.ParseInt(value, 10, 64)
			return id
		}

		fmt.Println("mapQuery key uid doesn't exist.")
		return 0
	}

	fmt.Println("queryRet is nil")
	return 0
}

func (db *BaseDb) InsertRole(uid int32, role_id int32, name string) bool {
	querySql := fmt.Sprintf("insert into shadow.role(user_id,role_type,role_name) values(%d,%d,'%s')", uid, role_id, name)
	if !db.Mysqldb.MyExec(querySql) {
		fmt.Println("Update User State fail")
		return false
	}
	return true
}

func (db *BaseDb) UpdateUserState(userid int64, state string) bool {

	querySql := fmt.Sprintf("update shadow.User set state='%s' where uid =%d ", state, userid)
	if !db.Mysqldb.MyExec(querySql) {
		fmt.Println("Update User State fail")
	}
	return false
}
func (db *BaseDb) GetSearchUserList(info string) []int {
	querySql := fmt.Sprintf("SELECT * FROM shadow.User where name like '%s' or uid like '%s'", "%"+info+"%", "%"+info+"%")

	fmt.Println("querySql:", querySql)
	slice:= make([]int, 0)
	queryRet := db.Mysqldb.MyQuery(querySql)
	if queryRet != nil {

		for i := 0; i < len(queryRet); i++ {
			mapQuery := queryRet[i]
			fmt.Println("mapQuery:", mapQuery)
			value, ok := mapQuery["uid"]
			if ok {
				uid, _ := strconv.Atoi(value)
				slice = append(slice,uid)
			}else{
				fmt.Println("mapQuery key uid doesn't exist.")
			}

		}
	}
	fmt.Println("queryRet is nil")
	return slice
}


func (db *BaseDb) GetUserName(uid int) string {
	querySql := fmt.Sprintf("SELECT name FROM shadow.User where uid = %d", uid)

	fmt.Println("querySql:", querySql)
	queryRet := db.Mysqldb.MyQuery(querySql)
	if queryRet != nil {
		mapQuery := queryRet[0]

		fmt.Println("mapQuery:", mapQuery)
		value, ok := mapQuery["name"]
		if ok {

			return value
		}

		fmt.Println("mapQuery key uid doesn't exist.")
		return ""
	}

	fmt.Println("queryRet is nil")
	return ""

}


func (db *BaseDb) InsertFriendInvitationInfo(userid int, friendid int)bool{
	querySql := fmt.Sprintf("insert into shadow.Friend_Invitation(userid,friendid,status) values(%d,%d, %d)", userid, friendid,  0)//0代表未激活
	if !db.Mysqldb.MyExec(querySql) {
		fmt.Println("InsertFriendInvitationInfo fail")
		return false
	}
	return true
}

func (db *BaseDb) InsertNewFriend(userid int, friendid int, friendtype string)bool {
	if friendtype == "" {
		friendtype= "我的好友"
	}
	querySql := fmt.Sprintf("insert into shadow.Friend(userid,friendid,status) values(%d,%d,%s)", userid, friendid,friendtype)
	if !db.Mysqldb.MyExec(querySql) {
		fmt.Println("InsertFriendInvitationInfo fail")
		return false
	}
	return true
}

func (db *BaseDb) GetFriendList(userid int)[]int {
	querySql := fmt.Sprintf("SELECT * FROM shadow.Friend where userid = %d", userid)

	fmt.Println("querySql:", querySql)
	slice:= make([]int, 0)
	queryRet := db.Mysqldb.MyQuery(querySql)
	if queryRet != nil {

		for i := 0; i < len(queryRet); i++ {
			mapQuery := queryRet[i]
			fmt.Println("mapQuery:", mapQuery)
			value, ok := mapQuery["friendid"]
			if ok {
				uid, _ := strconv.Atoi(value)
				slice = append(slice,uid)
			}else{
				fmt.Println("mapQuery key uid doesn't exist.")
			}

		}
	}
	fmt.Println("queryRet is nil")
	return slice
}

func (db *BaseDb) GetFriendInvitationInfo(friendid int)[]int {
	querySql := fmt.Sprintf("SELECT * FROM shadow.Friend_Invitation where friendid = %d", friendid)

	fmt.Println("querySql:", querySql)
	slice:= make([]int, 0)
	queryRet := db.Mysqldb.MyQuery(querySql)
	if queryRet != nil {

		for i := 0; i < len(queryRet); i++ {
			mapQuery := queryRet[i]
			fmt.Println("mapQuery:", mapQuery)
			value, ok := mapQuery["userid"]
			if ok {
				uid, _ := strconv.Atoi(value)
				slice = append(slice,uid)
			}else{
				fmt.Println("mapQuery key uid doesn't exist.")
			}

		}
	}
	fmt.Println("queryRet is nil")
	return slice
}
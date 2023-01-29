package db

import (
	"Src/libs/logger"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlClient struct {
	*sql.DB
	dbSourceName string
}

func MysqlConnect(dbSourceName string) *MysqlClient {
	db, err := sql.Open("mysql", dbSourceName)
	if nil != err {
		logger.Error("Open(mysql,%s) failed, error:%s",
			dbSourceName,
			err.Error())
		return nil
	}
	// db对象并不是数据库连接，在需要的时候才会创建，
	// 需要立即验证数据库是否可用
	err = db.Ping()
	if nil != err {
		logger.Error("mysql Ping() failed, error:%s", err.Error())
		return nil
	}
	return &MysqlClient{
		DB:           db,
		dbSourceName: dbSourceName,
	}
}

// 查询某个server桌子的信息，每个桌子的信息以key、value的形式存入map
// 然后map存入list
func (c *MysqlClient) MyQuery(querySql string, args ...interface{}) []map[string]string {
	rows, err := c.DB.Query(querySql, args...)
	if nil != err {
		logger.Error("db.Query(%s) failed, error:%s", querySql, err.Error())
		return nil
	}
	// 获取列名
	columns, err := rows.Columns()
	if nil != err {
		logger.Error("rows.Columns() failed, error:%s", err.Error())
		return nil
	}
	// Make a slice for the value
	values := make([]sql.RawBytes, len(columns))
	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// 格式化的结果集
	var tableList []map[string]string
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			logger.Error("rows.Scan() failed, error:%s", err.Error())
			continue
		}
		// 这个map用来存储一行数据，列名为map的key，map的value为列的值
		rowMap := make(map[string]string)
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col != nil {
				value = string(col)
				rowMap[columns[i]] = value
			}
			//
			if col == nil {
				value = string("NULL")
				rowMap[columns[i]] = value
			}
		}
		tableList = append(tableList, rowMap)
	}
	return tableList

}

func (c *MysqlClient) MyExec(querySql string, args ...interface{}) bool {
	ret, err := c.DB.Exec(querySql, args...)
	if nil != err {
		logger.Error("db.Exec(%s) failed, error:%s", querySql, err.Error())
		return false
	}
	theID, err := ret.LastInsertId()
	// 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return false
	}

	fmt.Printf("Exec success, the id is %d.\n", theID)
	return true
}

package help

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var timeTemplates = []string{
	"2006-01-02 15:04:05", //常规类型
	"2006/01/02 15:04:05",
	"2006-01-02",
	"2006/01/02",
	"15:04:05",
}

/* 时间格式字符串转换 */
func TimeStringToGoTime(tm string) time.Time {
	for i := range timeTemplates {
		t, err := time.ParseInLocation(timeTemplates[i], tm, time.Local)
		if nil == err && !t.IsZero() {
			return t
		}
	}
	return time.Time{}
}

/**
jsonstring --> map
*/
func JsonStr2Map(jsonStr string) (event map[string]interface{}, err error) {
	if err = json.Unmarshal([]byte(jsonStr), &event); err != nil {
		panic(err)
	}
	return
}

/**
jsonstring --> struct
*/
func JsonStr2Struct(jsonStr string, eventStruct interface{}) {
	if err := json.Unmarshal([]byte(jsonStr), &eventStruct); err != nil {
		panic(err)
	}
}

/**
map --> struct
*/
func Map2Struct(mapBean map[string]interface{}, eventStruct interface{}) {
	//将 map 转换为指定的结构体
	str, err := Map2JsonStr(mapBean)
	if err != nil {
		fmt.Println("err = ", err)
	}
	JsonStr2Struct(str, &eventStruct)
}

/**
map --> jsonstring
*/

func Map2JsonStr(mapBean map[string]interface{}) (str string, err error) {
	bytes, err := json.Marshal(mapBean)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	return string(bytes), err
}

/**
struct -> jsonString
*/
func Struct2JsonStr(eventStruct interface{}) (str string, err error) {
	buf, err := json.Marshal(eventStruct) //格式化编码
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	return string(buf), err
}

/**
struct -> map
*/
func Struct2Map(obj interface{}) map[string]interface{} {
	//获取参数o的类型
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	field := t.NumField()

	var data = make(map[string]interface{})
	for i := 0; field > i; i++ {

		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func String2Interface(mapString map[string]string) map[string]interface{} {

	jsonMap, _ := json.Marshal(mapString)
	var mapObject = map[string]interface{}{}
	json.Unmarshal(jsonMap, &mapObject)
	return mapObject
}

// Strval 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

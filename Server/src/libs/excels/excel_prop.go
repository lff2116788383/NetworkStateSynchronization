package excels

import "Src/libs/utils"


const (
	ITEM_TYPE_DRUG    = "药品"
	ITEM_TYPE_FOOD    = "食物"
	ITEM_TYPE_EXP     = "经验"
	ITEM_TYPE_EXCHANGE  = "兑换"
	ITEM_TYPE_MONEY     = "金钱"
	ITEM_TYPE_BOX     = "盒子"
	ITEM_TYPE_SKILL     = "技能"


)

//读取excel转换成结构体
type ConfigProp struct {
	Id int `json:"id"`
	ItemId  int64 `json:"itemid"`
	Level  int `json:"level"`
	Name  string `json:"propName"`
	Price  int `json:"price"`
	Type  string `json:"type"`
	Result  string `json:"result"`
	Desc  string `json:"description"`
	Dwelltime  string `json:"dwelltime"`
	SellPrice  int `json:"sellprice"`
	OverLay  int `json:"overlay"`
	Source  string `json:"source"`


}

var (
	ConfigPropMap         map[int]*ConfigProp

)
//每一个源文件都可以包含一个 init  函数，该函数会在 main 函数执行前，被 Go 运行框架调用，也就是说 init 会在 main 函数前被调用。
func init() {
	ConfigPropMap = make(map[int]*ConfigProp)
	utils.GetExcelUtilMgr().LoadExcel("Prop",&ConfigPropMap)

}

func GetPropConfig(propId int) *ConfigProp {
	return ConfigPropMap[propId]
}
package excels

import "Src/libs/utils"

type ConfigEquipment struct {
	Id int `json:"id"`
	ItemId  int64 `json:"itemid"`
	Level  int `json:"level"`
	Name  string `json:"name"`
	Class  string `json:"class"`
	STR  int `json:"STR"`	//力量
	INT  int `json:"INT"`	//智力
	DEX  int `json:"DEX"`	//敏捷

	HP  int `json:"HP"`
	MP  int `json:"MP"`

	AD int	`json:"AD"`
	Desc  string `json:"Description"`

	Source  string `json:"Source"`

	AP int	`json:"AP"`

	ADEF int	`json:"ADEF"`
	MDEF int	`json:"MDEF"`
	CRIT   float32	`json:"CRIT"`


}


var (
	ConfigEquipmentMap         map[int]*ConfigEquipment

)
//每一个源文件都可以包含一个 init  函数，该函数会在 main 函数执行前，被 Go 运行框架调用，也就是说 init 会在 main 函数前被调用。
func init() {
	ConfigEquipmentMap = make(map[int]*ConfigEquipment)
	utils.GetExcelUtilMgr().LoadExcel("Equipment",&ConfigEquipmentMap)

}

func GetEquipmentConfig(Id int) *ConfigEquipment {
	return ConfigEquipmentMap[Id]
}
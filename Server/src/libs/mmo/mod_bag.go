package mmo

import (
	"Src/libs/csvs"
	"Src/libs/excels"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)
//道具信息
type ItemInfo struct {
	ItemId  int32
	ItemNum int64
	//ItemWeight int32
}

type ModBag struct {
	Id int32
	Name string

	BagInfo map[int32]*ItemInfo//物品信息

	capacity *DataPool 	//容量
	user *User			//拥有者
	path   string
}

func (self *ModBag) AddItem(itemId int32, num int64) {
	
	
	//根据不同范围的itemId	加载不同的excel配置表数据
	if itemId>=42001 &&itemId<= 42010{
		config:=excels.GetPropConfig(int(itemId))
		if config == nil {
			fmt.Println(itemId, "物品不存在")
			return
		}
		//添加道具到背包
		self.AddItemToBag(itemId, num)
	}
	
	//配置表添加ItemWeight属性
	//添加物品超过最大容量	weight*num	>	capacity.Max
	//if (self.capacity.Current + ((itemConfig.ItemWeight)*int32(num)))> self.capacity.Max {
	//	fmt.Println("添加的物品重量超过背包容量	,背包当前容量","添加物品的重量","背包最大容量")
	//}
	
	
}

func (self *ModBag) AddItemToBag(itemId int32, num int64) {
	//TODO

	//如果背包的当前容量加上添加的物品的重量大于背包的最大容量
	////背包当前容量加上要添加的物品的重量

	_, ok := self.BagInfo[itemId]
	if ok {
		//存在直接添加数量
		self.BagInfo[itemId].ItemNum += num
	} else {
		self.BagInfo[itemId] = &ItemInfo{ItemId: itemId, ItemNum: num}
	}
	//config := csvs.GetItemConfig(int(itemId))
	//if config != nil {
	//	fmt.Println("获得物品", config.ItemName, "----数量：", num, "----当前数量：", self.BagInfo[itemId].ItemNum)
	//}

}

func (self *ModBag) RemoveItem(itemId int32, num int64) {
	//itemConfig := csvs.GetItemConfig(int(itemId))
	//if itemConfig == nil {
	//	fmt.Println(itemId, "物品不存在")
	//	return
	//}
	//
	//switch itemConfig.SortType {
	//case csvs.ITEMTYPE_NORMAL:
	//	self.RemoveItemToBagGM(itemId, num)
	//default: //同普通
	//	//self.AddItemToBag(itemId, 1)
	//}
	self.RemoveItemToBagGM(itemId,num)

}

func (self *ModBag) RemoveItemToBagGM(itemId int32, num int64) {
	_, ok := self.BagInfo[itemId]
	if ok {
		self.BagInfo[itemId].ItemNum -= num
	} else {
		self.BagInfo[itemId] = &ItemInfo{ItemId: itemId, ItemNum: 0 - num}
	}
	//config := csvs.GetItemConfig(int(itemId))
	//if config != nil {
	//	fmt.Println("扣除物品", config.ItemName, "----数量：", num, "----当前数量：", self.BagInfo[itemId].ItemNum)
	//}
}

func (self *ModBag) RemoveItemToBag(itemId int32, num int64 ){
	itemConfig := csvs.GetItemConfig(int(itemId))
	switch itemConfig.SortType {
	//case csvs.ITEMTYPE_NORMAL:
	//	self.AddItemToBag(itemId, num)
	case csvs.ITEMTYPE_ROLE:
		fmt.Println("此物品无法扣除")
		return
	case csvs.ITEMTYPE_ICON:
		fmt.Println("此物品无法扣除")
		return
	case csvs.ITEMTYPE_CARD:
		fmt.Println("此物品无法扣除")
		return
	default: //同普通
	}

	if !self.HasEnoughItem(itemId, num) {
		config := csvs.GetItemConfig(int(itemId))
		if config != nil {
			nowNum := int64(0)
			_, ok := self.BagInfo[itemId]
			if ok {
				nowNum = self.BagInfo[itemId].ItemNum
			}
			fmt.Println(config.ItemName, "数量不足", "----当前数量：", nowNum)
		}
		return
	}

	_, ok := self.BagInfo[itemId]
	if ok {
		self.BagInfo[itemId].ItemNum -= num
	} else {
		self.BagInfo[itemId] = &ItemInfo{ItemId: itemId, ItemNum: 0 - num}
	}
	fmt.Println("扣除物品", itemConfig.ItemName, "----数量：", num, "----当前数量：", self.BagInfo[itemId].ItemNum)
}

func (self *ModBag) HasEnoughItem(itemId int32, num int64) bool {
	if itemId == 0 {
		return true
	}
	_, ok := self.BagInfo[itemId]
	if !ok {
		return false
	} else if self.BagInfo[itemId].ItemNum < num {
		return false
	}
	return true
}

func (self *ModBag) UseItem(itemId int32, num int64) {
	itemConfig := excels.GetPropConfig(int(itemId))
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}

	if !self.HasEnoughItem(itemId, num) {
		config := excels.GetPropConfig(int(itemId))
		if config != nil {
			nowNum := int64(0)
			_, ok := self.BagInfo[itemId]
			if ok {
				nowNum = self.BagInfo[itemId].ItemNum
			}
			fmt.Println(config.Name, "数量不足", "----当前数量：", nowNum)
		}
		return
	}

	switch itemConfig.Type {
	case excels.ITEM_TYPE_DRUG:
		switch	itemId {
		case 1:
			self.user.GetModRole().CtrlRole.Hp.Current+=500
		case 2:
			self.user.GetModRole().CtrlRole.Mp.Current+=500
		}


	case excels.ITEM_TYPE_FOOD:
		//给英雄加属性
		switch	itemId {
		case 3:
			self.user.GetModRole().CtrlRole.Hp.Current+=1000
		case 4:
			self.user.GetModRole().CtrlRole.Mp.Current+=1000
		}
	case excels.ITEM_TYPE_EXP:
		switch	itemId {
		case 5:
			self.user.GetModRole().CtrlRole.Exp.Current+=10000
		}
	case excels.ITEM_TYPE_EXCHANGE:
		switch	itemId {
		case 6:
			if self.BagInfo[itemId].ItemNum < 10 {
				fmt.Println(itemId, "此物品数量不足无法兑换")
			}
		}

	case excels.ITEM_TYPE_MONEY:
		switch itemId {
		case 7:
			self.user.Money+= 1000
		case 8:
			self.user.Money+= 10000
		}
	case excels.ITEM_TYPE_BOX:
		switch itemId {
		case 9:
		}
	case excels.ITEM_TYPE_SKILL:
		switch itemId {
		case 10:
		}
	default: //同普通
		fmt.Println(itemId, "此物品无法使用")
		return
	}

	//角色属性值限制
	self.user.GetModRole().CtrlRole.Hp.Limit()
	self.user.GetModRole().CtrlRole.Mp.Limit()
	self.user.GetModRole().CtrlRole.Exp.Limit()
}

func (self *ModBag) UseCookBook(itemId int32, num int64) {
	cookBookConfig := csvs.GetCookBookConfig(int(itemId))
	if cookBookConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}
	self.RemoveItem(itemId, num)
	self.AddItem(int32(cookBookConfig.Reward), num)
}

func (self *ModBag) GetItemNum(itemId int32) int64 {
	itemConfig := csvs.GetItemConfig(int(itemId))
	if itemConfig == nil {
		return 0
	}
	_, ok := self.BagInfo[itemId]
	if !ok {
		return 0
	}
	return self.BagInfo[itemId].ItemNum
}


//--------------------------Mod接口方法
func (self *ModBag) SaveData() {
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (self *ModBag) LoadData(user *User) {

	self.user = user
	self.path = self.user.localPath + "/bag.json"

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

	if self.BagInfo == nil {
		self.BagInfo = make(map[int32]*ItemInfo)
	}
	return
}

func (self *ModBag) InitData() {
	self.capacity=&DataPool{0,100}
}

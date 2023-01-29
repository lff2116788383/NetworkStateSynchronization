package mmo

import (
	"Src/libs/csvs"
	"Src/libs/excels"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//还是挂载到用户上	只关心用户数据
const (
	STORE_TYPE_NORMAL   =1		//普通道具商店
	STORE_TYPE_PET		=2		//宠物商店
	STORE_TYPE_VIP		=9		//VIP商店
)

type Goods struct {
	ItemId  	int32
	Name        string
	Price       int64	//购买价格
	SellPrice	int64   //售卖价格
	ItemNum 	int64
	//ItemWeight  int32
}
//普通商店系统	物品只显示价格	不显示剩余数量
type ModStore struct {
	GoodsList map[int32]*Goods//商品信息
	Type int32

	user *User
	path string
}

//普通购买
func (this*ModStore)Buy(itemId int32,num int64) bool {
	//判断钱是否不够
	goods:=this.GetGoods(itemId)
	if goods == nil {
		return false
	}
	buy_money:=int64(goods.Price *num)
	if this.user.Money < buy_money {
		return false
	}

	//钱足够 减去用户money	 用户添加item
	this.user.Money-=buy_money
	this.user.GetModBag().AddItem(itemId,num)

	return true
}

//Vip购买	打八折
func (this*ModStore)VIPBuy(itemId int32,num int64) bool {
	discount:=0.8

	goods:=this.GetGoods(itemId)
	if goods == nil {
		return false
	}
	buy_money:=int64(float64(goods.Price *num)*discount)
	//判断钱是否不够
	if this.user.Money < buy_money {
		return false
	}

	//钱足够 减去用户money	 用户添加item
	this.user.Money-=buy_money
	this.user.GetModBag().AddItem(itemId,num)

	return true
}


//出售
func (this*ModStore)Sell(itemId int32,num int64) bool {

	goods:=this.GetGoods(itemId)
	if goods == nil {
		return false
	}
	sell_money:=int64(goods.SellPrice *num)

	//增加用户money	 用户删除item
	this.user.Money+=sell_money
	this.user.GetModBag().RemoveItem(itemId,num)

	return true
}

func  (this*ModStore)GetGoods(itemId int32)*Goods  {
	itemConfig := csvs.GetItemConfig(int(itemId))
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return nil
	}
	_,ok:=this.GoodsList[itemId]
	if !ok {
		return nil
	}

	return this.GoodsList[itemId]
}



//--------------------------Mod接口方法
//path路径保存角色数据
func (this *ModStore) SaveData() {
	content, err := json.Marshal(this)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(this.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

//localPath加载角色数据
func (this *ModStore) LoadData(user *User) {

	this.user = user
	this.path = this.user.localPath + "/store.json"

	configFile, err := ioutil.ReadFile(this.path)
	if err != nil {
		this.InitData()
		return
	}
	err = json.Unmarshal(configFile, &this)
	if err != nil {
		this.InitData()
		return
	}

	if this.GoodsList == nil {
		this.GoodsList = make(map[int32]*Goods)
	}
	return
}

//加载中进行初始化数据
func (this *ModStore) InitData() {
	if this.GoodsList == nil {
		this.GoodsList = make(map[int32]*Goods)

		//加载道具商店	武器商店		宠物商店
		for k,v :=range excels.ConfigPropMap{
			this.GoodsList[int32(k)]=&Goods{int32(v.Id),v.Name,int64(v.Price),int64(v.SellPrice),v.ItemId}
		}
	}
}

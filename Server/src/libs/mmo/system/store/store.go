package store

import (
	. "Src/libs/mmo"
)

//还是挂载到用户上	只关心用户数据
const (
	STORE_TYPE_NORMAL   =1		//普通道具商店
	STORE_TYPE_PET		=2		//宠物商店
	STORE_TYPE_VIP		=9		//VIP商店
)

type Goods struct {
	ItemId  	int32
	
	Price       int64	//购买价格
	SellPrice	int64   //售卖价格
	ItemNum 	int64
	ItemWeight  int32
}
//普通商店系统	物品只显示价格	不显示剩余数量
type Store struct {
	GoodsList map[int32]*Goods//商品信息
	Type int
}

//普通购买
func (s*Store)Buy(user *User,itemId int32,num int64) int {
	//判断钱是否不够
	buy_money:=int64(s.GoodsList[itemId].Price *num)
	if user.Money < buy_money {
		return -1
	}

	//钱足够 减去用户money	 用户添加item
	user.Money-=buy_money
	user.GetModBag().AddItem(int(itemId),num)

	return 0
}

//Vip购买	打八折
func (s*Store)VIPBuy(user *User,itemId int32,num int64) int {
	discount:=0.8
	buy_money:=int64(float64(s.GoodsList[itemId].Price *num)*discount)
	//判断钱是否不够
	if user.Money < buy_money {
		return -1
	}

	//钱足够 减去用户money	 用户添加item
	user.Money-=buy_money
	user.GetModBag().AddItem(int(itemId),num)

	return 0
}

//u.server.store.buy(u,itemId)
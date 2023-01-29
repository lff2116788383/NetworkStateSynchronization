package GameServer
import (
"Src/libs/csvs"
"Src/libs/excels"
"encoding/json"
"fmt"
"io/ioutil"
"os"
)

//**************************系统商店模块还是与Server绑定 用户获取Server的商店模块进行买入卖出      --与用户绑定	数据冗余	个人商店可以与用户绑定

const (
	STORE_TYPE_NORMAL   =1		//普通道具商店
	STORE_TYPE_PET		=2		//宠物商店
	STORE_TYPE_VIP		=9		//VIP商店
)

type Item struct {
	ItemId  	int
	Name        string
	Price       int64	//购买价格
	SellPrice	int64   //售卖价格
	ItemNum 	int64
	//ItemWeight  int
}
//普通商店系统	物品只显示价格	不显示剩余数量
type ModStore struct {
	ItemList map[int]*Item//商品信息
	Type int

	server *GameServer
	path string
}

//普通购买
func (this*ModStore)Buy(pUser *CGameUser,itemId int,num int64) bool {
	//判断钱是否不够
	goods:=this.GetItem(itemId)
	if goods == nil {
		return false
	}
	buy_money:=int64(goods.Price *num)
	if pUser.Money < buy_money {
		return false
	}

	//钱足够 减去用户money	 用户添加item 这里需要更新用户信息
	pUser.Money-=buy_money
	pUser.GetModBag().AddItem(itemId,num)

	return true
}

//Vip购买	打八折
func (this*ModStore)VIPBuy(pUser *CGameUser,itemId int,num int64) bool {
	discount:=0.8

	goods:=this.GetItem(itemId)
	if goods == nil {
		return false
	}
	buy_money:=int64(float64(goods.Price *num)*discount)
	//判断钱是否不够
	if pUser.Money < buy_money {
		return false
	}

	//钱足够 减去用户money	 用户添加item
	pUser.Money-=buy_money
	pUser.GetModBag().AddItem(itemId,num)

	return true
}


//出售
func (this*ModStore)Sell(pUser *CGameUser,itemId int,num int64) bool {

	goods:=this.GetItem(itemId)
	if goods == nil {
		return false
	}
	sell_money:=int64(goods.SellPrice *num)

	//增加用户money	 用户删除item
	pUser.Money+=sell_money
	pUser.GetModBag().RemoveItem(itemId,num)

	return true
}

func  (this*ModStore)GetItem(itemId int)*Item  {
	itemConfig := csvs.GetItemConfig(int(itemId))
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return nil
	}
	_,ok:=this.ItemList[itemId]
	if !ok {
		return nil
	}

	return this.ItemList[itemId]
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
func (this *ModStore) LoadData(server *GameServer) {

	this.server = server
	this.path = server.localPath + "/store.json"

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

	if this.ItemList == nil {
		this.ItemList = make(map[int]*Item)
	}
	return
}

//加载中进行初始化数据
func (this *ModStore) InitData() {
	if this.ItemList == nil {
		this.ItemList = make(map[int]*Item)

		//加载道具商店	武器商店		宠物商店
		for k,v :=range excels.ConfigPropMap{
			this.ItemList[int(k)]=&Item{int(v.Id),v.Name,int64(v.Price),int64(v.SellPrice),v.ItemId}
		}
	}
}

package main

import (
	"Src/libs/command"
	"Src/libs/mmo"

	//"fmt"
	//"unsafe"
	//. "Src/libs/utils"
	. "Src/libs/excels"
	"Src/libs/mmo/data"
	"Src/libs/packet"
	"fmt"
)

func main() {
	//var i1 int = 1
	//var i2 int8 = 2
	//var i3 int16 = 3
	//var i4 int32 = 4
	//var i5 int64 = 5
	//fmt.Println(unsafe.Sizeof(i1))
	//fmt.Println(unsafe.Sizeof(i2))
	//fmt.Println(unsafe.Sizeof(i3))
	//fmt.Println(unsafe.Sizeof(i4))
	//fmt.Println(unsafe.Sizeof(i5))

	//var sl  = make(map[int]*ConfigProp)
	//GetExcelUtilMgr().LoadExcel("Prop",&sl)
	//for k,v:=range sl{
	//	fmt.Println("key:",k,"value:",v.PropId, v.Num, v.Level, v.Name, v.Price, v.Type, v.Result, v.Desc, v.Dwelltime, v.SellPrice, v.OverLay, v.GetMethod)
	//}

	v:=GetPropConfig(1)
	fmt.Println("value:",v.Id, v.ItemId, v.Level, v.Name, v.Price, v.Type, v.Result, v.Desc, v.Dwelltime, v.SellPrice, v.OverLay, v.Source)

	v2:=GetEquipmentConfig(1)
	fmt.Println("value:",v2.Id, v2.Name, v2.Desc, v2.Class, v2.Source, v2.CRIT)


	u := &mmo.User{
		Id:   -1,
		Tid:  -1,
		UserInfo:new(data.NUserInfo),
	}
	u.UserInfo.InitData()
	ret:=0

	character:=new(data.NCharacterInfo)
	character.InitData()

	//character.Id = pack.ReadInt()
	character.Name = "eee"
	character.Class = 2
	character.Level =1
	character.MapId =1

	//if !db.GetDB().InsertRole(u.Id, int32(character.Class), character.Name) {
	//	//注册失败 数据库添加失败
	//	ret = 1
	//}
	//character.Id = int(db.GetDB().GetRoleId(u.Id, int32(character.Class), character.Name))



	fmt.Println("debug 1")
	character.Entity.Position.X = 1500
	character.Entity.Position.Y = 2600
	character.Entity.Position.Z = 150

	character.Gold = 10000

	fmt.Println("debug 2")
	character.AttrDynamic.HP = 100
	character.AttrDynamic.MP = 50


	fmt.Println("debug 3")
	character.Equips = make([]byte,28)


	bag:=new(data.NBagInfo)
	bag.Items = make([]byte,0)
	bag.Unlocked = 20

	character.Bag =bag

	fmt.Println("debug 4")

	item:=new(data.NItemInfo)
	item.Id =1
	item.Count =20
	character.Items[item.Id]=item
	fmt.Println("debug 5")

	u.UserInfo.Player.Characters[character.Id]=character
	fmt.Println("debug 6")

	//创建角色成功	发送角色列表
	outICP := packet.NewICPacket()
	outICP.Begin(command.CMD_SERVER_CREATE_ROLE_RESP)
	outICP.WriteByte(byte(ret))


	if ret==0 {
		outICP.WriteInt32Ex(len(u.UserInfo.Player.Characters))

		for _, v := range u.UserInfo.Player.Characters {
			outICP.WriteInt32Ex(v.Id)
			outICP.WriteString(v.Name)
			outICP.WriteInt32Ex(data.CHARACTER_TYPE_PLAYER)
			outICP.WriteInt32Ex(v.Class)
			outICP.WriteInt32Ex(v.ConfigId)

		}
	}

	outICP.End()

	fmt.Println(character.Class)
}
//输出结果：
//8
//1
//2
//4
//8
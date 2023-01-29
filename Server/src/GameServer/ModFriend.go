package GameServer

import (
	"Src/libs/packet"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)


type FriendInfo = CGameUser
//type FriendInfo struct {
//	Id int
//	name string
//	Status int
//}
//好友模块
type ModFriend struct {
	FriendList map[int]*FriendInfo 	//改成CGameUser

	user *CGameUser
	path string
}

func (self* ModFriend) AddFriend(uid int) int {
	//1.玩家是否存在
	pUser:=GetGameServer().GetUser(uid)
	if pUser == nil {
		return -1
	}
	//2.玩家是否已经是好友
	firend:=self.GetFriend(uid)

	if firend != nil {
		return -2
	}

	self.FriendList[pUser.m_nUserID] = pUser
	pUser.GetModFriend().FriendList[self.user.m_nUserID] = self.user

	return 0
}

func (self* ModFriend) DelFriend(uid int) int  {
	pUser:=GetGameServer().GetUser(uid)
	if pUser == nil {
		return -1
	}
	
	if !self.HasFriend(uid) {
		return -2	//好友不存在
	}
	delete(self.FriendList,pUser.m_nUserID)
	delete(pUser.GetModFriend().FriendList,self.user.m_nUserID)
	return 0
}

func (self* ModFriend) GetFriend(uid int)* FriendInfo  {
	friend, ok := self.FriendList[uid] /*如果确定是真实的,则存在,否则不存在 */
	if (ok) {
		return friend
	}
	return nil
}

func (self* ModFriend) HasFriend(uid int)bool  {
	_, ok := self.FriendList[uid] /*如果确定是真实的,则存在,否则不存在 */
	if (ok) {
		return true
	}
	return false
}

func (self* ModFriend) GetFriendListCount() int {
	count:=0
	for _,v:= range self.FriendList {
		if v != nil {
			count++
		}
	}
	return  count
}

func (self* ModFriend) FriendChat(uid int, msg string) int {
	pUser:=GetGameServer().GetUser(uid)
	if pUser == nil {
		return -1
	}

	if !self.HasFriend(uid) {
		return -2	//好友不存在
	}

	if pUser.m_pSocket== nil {
		return -3
	}

	outICP := packet.NewICPacket()
	outICP.Begin(SERVER_COMMAND_FRIEND_CHAT)
	outICP.WriteString(msg)
	outICP.End()
	pUser.m_pSocket.Send(outICP)
	
	return 0
}


//--------------------------Mod接口方法
func (self *ModFriend) SaveData() {
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (self *ModFriend) LoadData(user *CGameUser) {

	self.user = user
	self.path = self.user.localPath + "/pet.json"

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

	if self.FriendList == nil {
		self.FriendList = make(map[int]*FriendInfo)
	}
	return
}

func (self *ModFriend) InitData() {

}
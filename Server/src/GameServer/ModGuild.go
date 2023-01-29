package GameServer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	DEFALUT_DESC = "暂无描述"
)
//公会信息
type GuildInfo struct {
	Id int				//公会序号
	GuildName string	//公会名字
	President *CGameUser	//会长
	GuildPurpose string	//公会宗旨

	MemberLimit int	//成员上限
	OnlineMemberCount int //当前在线公会成员人数
	CreateTime time.Time	//创建时间
	GuildMemberList map[int]*GuildMemberInfo	//公会成员列表
	GuildApplyList map[int]*GuildApplyInfo	//公会申请列表
}
var GuildStartId = 10000

func NewGuild(pUser *CGameUser) *GuildInfo {
	GuildStartId+=1

	guild:= &GuildInfo{
		Id: GuildStartId,
		GuildName:DEFALUT_DESC,
		President: pUser,
		GuildPurpose: DEFALUT_DESC,
		MemberLimit:250,
		OnlineMemberCount:0,
		CreateTime: time.Now(),
		GuildMemberList:make(map[int]*GuildMemberInfo),
		GuildApplyList: make(map[int]*GuildApplyInfo),
	}



	return guild
}


func (self*GuildInfo) InitData()  {
	if self.GuildMemberList == nil {
		self.GuildMemberList = make(map[int]*GuildMemberInfo)
	}

	if self.GuildApplyList == nil {
		self.GuildApplyList = make(map[int]*GuildApplyInfo)
	}
}

func (self*GuildInfo) AddMember(member *GuildMemberInfo)  {

		self.GuildMemberList[member.User.m_nUserID] = member

		member.User.guild = self

}

func (self*GuildInfo) DelMember(uid int)  {

	delete(self.GuildMemberList,uid)


}

func (self*GuildInfo) GetMember(uid int) *GuildMemberInfo {

	for k,v:= range self.GuildMemberList  {
		if k == uid  {
			return v
		}
	}

	return nil

}


//公会申请加入信息
type GuildApplyInfo struct {
	Id int
	GuildId int
	CharacterId int
	Name string
	Class int
	Level int
	Result int
}

type GuildMemberInfo struct {
	User *CGameUser
	GuildName string //用户公会昵称
	JoinTime time.Time	//加入公会时间
	LastTime time.Time	//上次上线时间
	Status int		//状态	0：下线 1：在线 2：隐身
	MemberType int  //公会成员类型 0：普通公会成员	1：公会副会长 2：公会会长
}

func NewGuildMember(pUser *CGameUser, memberType int) *GuildMemberInfo {
	member:=new (GuildMemberInfo)
	member.User = pUser
	member.GuildName = DEFALUT_DESC
	member.Status = 1
	member.MemberType = memberType
	return member
}



//群和公会等多个用户相关的与server绑定
type ModGuild struct {
	GuildList map[int]*GuildInfo	//公会列表


	server * GameServer
	path string
}





func (self* ModGuild) AddGuild(guild * GuildInfo)  {
	self.GuildList[guild.Id] = guild
}

func (self* ModGuild) DelGuild(Id int)  {
	delete(self.GuildList,Id)
}

func (self* ModGuild) GetGuild(Id int)* GuildInfo  {
	guild, ok := self.GuildList[Id] /*如果确定是真实的,则存在,否则不存在 */
	if (ok) {
		return guild
	}
	return nil
}

func (self* ModGuild) GetGuildListCount() int {
	count:=0
	for _,v:= range self.GuildList {
		if v != nil {
			count++
		}
	}
	return  count
}




func (self* ModGuild) AddGuildApply(apply * GuildApplyInfo)  {
	for k,v:=range self.GuildList{
		if k == apply.GuildId {
			v.GuildApplyList[apply.Id] = apply
		}
	}
}




//--------------------------Mod接口方法
func (self *ModGuild) SaveData() {
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (self *ModGuild) LoadData(server *GameServer) {

	self.server = server
	self.path = self.server.localPath + "/guild.json"

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

	if self.GuildList == nil {
		self.GuildList = make(map[int]*GuildInfo)
	}


	return
}

func (self *ModGuild) InitData() {
	for _,v := range self.GuildList{
		v.InitData()
	}
}
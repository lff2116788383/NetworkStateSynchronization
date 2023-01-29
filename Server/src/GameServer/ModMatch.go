package GameServer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type MatchTeamMember struct {
	MemberType int	//成员类型 0：普通成员 1：队长
	member *CGameUser
}

type MatchTeam struct {
	Id int
	Members []*MatchTeamMember
}
//组队模块与server绑定
type ModMatch struct {
	TeamList map[int]*MatchTeam

	server *GameServer
	path string
}


func (self *ModMatch) HasTeam(teamId int) bool {
	_,ok:=self.TeamList[teamId]
	if  ok {
		return true
	}
	return false
}

func (self *ModMatch) DelTeam(teamId int)  {
	_,ok:=self.TeamList[teamId]
	if  ok {
		delete(self.TeamList,teamId)
	}
}

func (self *ModMatch) AddUserToTeam(pUser* CGameUser, teamId int)  {
	t,ok:=self.TeamList[teamId]
	if  ok {	//队伍存在	加入队伍 成为普通队员
		teamMember:=&MatchTeamMember{
			MemberType:0,
			member: pUser,
		}
		t.Members = append(t.Members,teamMember)
		return
	}

	//队伍不存在存在	加入队伍 成为队长
	team:=&MatchTeam{
		Id:GenerateTeamId(),
		Members: make([]*MatchTeamMember,0),
	}

	teamMember:=&MatchTeamMember{
		MemberType:1,
		member: pUser,
	}

	pUser.m_nTeamId = team.Id
	team.Members = append(team.Members,teamMember)

}
var TeamStartID int = 10000
func GenerateTeamId() int  {
	TeamStartID++
	return TeamStartID
}









//--------------------------Mod接口方法
func (self *ModMatch) SaveData() {
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (self *ModMatch) LoadData(server *GameServer) {

	self.server = server
	self.path = self.server.localPath + "/match.json"

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

	if self.TeamList == nil {
		self.TeamList = make(map[int]*MatchTeam)
	}
	return
}

func (self *ModMatch) InitData() {

}
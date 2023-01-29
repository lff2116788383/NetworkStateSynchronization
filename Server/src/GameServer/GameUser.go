package GameServer

import (
	"Src/libs/config"
	"Src/libs/packet"
	"fmt"
	"os"
	"time"
)

const(
	HEARTBEATTIME = 60
)

type Position struct{
	X float64
	Y float64
	Angle int
}


type CGameUser struct {
	m_nUserID int

	m_pSocket *SocketHandler

	Money int64		//这个应该放到ModUser

	//UserInfo *UserInfo	//用户信息相关数据太复杂	使用泛型拆分模块
	// 用户数据模块	使用泛型方便做扩展
	modManage map[string]ModBase
	localPath string //用户数据本地保存路径

	//定时器
	m_CheckHeartBeat *time.Timer	//心跳检测


	Name string
	m_nTeamId int



	//
	guild *GuildInfo

	Pos Position
	Action int
}

func NewUser(uid int,pos Position,  pSocket*SocketHandler) *CGameUser {
	//***************泛型架构***************************

	pNewUser:= &CGameUser{
		m_nUserID: uid,
		m_pSocket: pSocket,
		m_CheckHeartBeat: nil,
		modManage: map[string]ModBase{
			MOD_ROLE: new(ModRole),
			MOD_BAG: new(ModBag),
			MOD_FRIEND: new(ModFriend),
		},
		Pos: pos,
	}
	pNewUser.InitData()
	pNewUser.InitMod()
	return pNewUser
}


func (u *CGameUser) StartCheckHeartBeatTimer(){
	u.m_CheckHeartBeat = time.NewTimer(HEARTBEATTIME *time.Second)
}

func (u *CGameUser) StopCheckHeartBeatTimer(){
	u.m_CheckHeartBeat.Stop()
}

func (u *CGameUser) ResetCheckHeartBeatTimer(){
	if u.m_CheckHeartBeat == nil {
		u.StartCheckHeartBeatTimer()
	}
	u.m_CheckHeartBeat.Reset(HEARTBEATTIME *time.Second)
}



func (u *CGameUser) InitData() {
	path := config.GlobalConfig.LocalSavePath
	_, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return
		}
	}
	u.localPath = path + fmt.Sprintf("/%d", u.m_nUserID)
	_, err = os.Stat(u.localPath)
	if err != nil {
		err = os.Mkdir(u.localPath, os.ModePerm)
		if err != nil {
			return
		}
	}


}
func (u *CGameUser) InitMod() {
	for _, v := range u.modManage {
		v.LoadData(u)
	}
}

func (u *CGameUser) GetMod(modName string) ModBase {
	return u.modManage[modName]
}

func (u *CGameUser) GetModRole() *ModRole {
	return u.modManage[MOD_ROLE].(*ModRole)
}

func (u *CGameUser) GetModBag() *ModBag {
	return u.modManage[MOD_BAG].(*ModBag)
}

func (u *CGameUser) GetModFriend() *ModFriend {
	return u.modManage[MOD_FRIEND].(*ModFriend)
}


//用户创建公会
func (u *CGameUser) CreateGuild(GuildName string)  {
	guild := NewGuild(u)
	member:= NewGuildMember(u,2)
	guild.AddMember(member)
	GetGameServer().GetSysModGuild().AddGuild(guild)
}

//获取用户所在的公会
func (u *CGameUser) GetGuild() *GuildInfo {
	return u.guild
}

func (u *CGameUser) GetGuildMemberType() int {
	if u.guild == nil {
		return -1
	}
	member := u.guild.GetMember(u.m_nUserID)
	if member == nil {
		return -2
	}
	return  member.MemberType
}

func (u *CGameUser) JoinGuild(Id int) {

	if u.guild != nil {
		return
	}
	member:= NewGuildMember(u,0)	//默认进入是普通成员
	member.JoinTime = time.Now()
	guild := GetGameServer().GetSysModGuild().GetGuild(Id)
	guild.AddMember(member)
}
func (u *CGameUser) LeaveGuild() int {

	if u.guild == nil {
		return -1	//用户未加入公会
	}

	if u.GetGuildMemberType() == 2 || u.GetGuildMemberType() == 1 {
		return -2	//会长/副会长不能直接离开
	}

	u.guild.DelMember(u.m_nUserID)
	u.guild = nil
	return 0

}

func (u *CGameUser) ModifyGuildPurpose(purpose string) int  {
	if u.guild == nil {
		return -1	//用户未加入公会
	}

	if u.GetGuildMemberType() != 2 && u.GetGuildMemberType() != 1 {
		return -2	//非会长/副会长无权
	}
	u.guild.GuildPurpose = purpose
	return 0
}

//用户转让公会
func (u *CGameUser) TransferGuild(pUser *CGameUser) int  {
	if u.guild == nil {
		return -1	//用户未加入公会
	}
	if u.GetGuildMemberType() != 2 {
		return -2	//非会长无权转让
	}

	if pUser.GetGuildMemberType() != 1 {
		return -3	//非副会长无权接受转让
	}

	if u.guild != pUser.guild {
		return  -4
	}

	u.guild.President = pUser
	u.guild.GetMember(u.m_nUserID).MemberType = 0
	pUser.guild.GetMember(pUser.m_nUserID).MemberType = 2


	return 0
}

//用户任命公会成员类型
func (u *CGameUser) AppointGuildMember(pUser *CGameUser, memberType int) int  {
	if u.guild == nil {
		return -1	//用户未加入公会
	}
	if u.GetGuildMemberType() != 2 {
		return -2	//非会长无权
	}

	if u.guild != pUser.guild {
		return  -3	//所在公会不同
	}

	//可以做一个数量限制
	pUser.guild.GetMember(pUser.m_nUserID).MemberType = memberType


	return 0
}

//用户逐出公会成员
func (u *CGameUser) OutGuildMember(pUser *CGameUser) int  {
	if u.guild == nil {
		return -1	//用户未加入公会
	}
	if u.GetGuildMemberType() != 2 && u.GetGuildMemberType() != 1 {
		return -2	//非会长/副会长无权
	}

	if u.guild != pUser.guild {
		return  -3	//所在公会不同
	}

	if pUser.GetGuildMemberType() != 0 {
		return  -4	//只能逐出普通成员
	}

	u.guild.DelMember(pUser.m_nUserID)
	pUser.guild = nil

	return 0
}



//用户公会私聊
func (u *CGameUser) GuildChat(pUser *CGameUser, msg string) int  {
	if u.guild == nil {
		return -1	//用户未加入公会
	}

	if u == pUser {
		return -2	//私聊对象不能是自己
	}

	if u.guild != pUser.guild {
		return  -3	//所在公会不同
	}

	if pUser.m_pSocket == nil {
		return -4
	}
	outICP := packet.NewICPacket()
	outICP.Begin(SERVER_COMMAND_HEARTBEAT)
	outICP.WriteString(msg)
	outICP.End()
	pUser.m_pSocket.Send(outICP)
	return 0
}

//获取申请列表
func (u *CGameUser) GetGuildApplyList() int {

	if u.guild == nil {
		return -1	//用户未加入公会
	}

	if u.GetGuildMemberType() != 2 && u.GetGuildMemberType() != 1 {
		return -2	//非会长/副会长无权
	}



	return 0

}
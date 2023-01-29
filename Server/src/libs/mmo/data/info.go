package data



type long = int64

type CharacterClass = int
type CharacterType = int
type QuestStatus = int
type GuildTitle = int
type ApplyResult = int

//const (
//		None = 0
//		Accept = 1
//		Reject = 2
//)

const (
	CHARACTER_TYPE_PLAYER = 0
	CHARACTER_TYPE_NPC = 1
	CHARACTER_TYPE_MONSTER = 2
)

type NUserInfo struct {
	Id int
	Player *NPlayerInfo
}

func (this *NUserInfo) InitData()  {
	if this.Player == nil {
		this.Player = new(NPlayerInfo)
		this.Player.InitData()
	}
}

type NPlayerInfo struct {
	Id int
	Characters	map[int]*NCharacterInfo
}

func (this *NPlayerInfo) InitData()  {
	if this.Characters == nil {
		this.Characters = make(map[int]*NCharacterInfo)
	}
}

type NCharacterInfo struct {
	Id int
	Name string
	Level int
	Exp long
	Gold long
	MapId int
	Ride int

	Class CharacterClass	//角色职业	如弓箭手 战士 法师
	Type  CharacterType		//角色类型	如玩家	怪物		NPC


	EntityId int
	ConfigId int

	Equips []byte

	Bag	*NBagInfo
	Guild *NGuildInfo

	Entity *NEntity
	AttrDynamic *NAttributeDynamic


	Friends map[int]*NFriendInfo
	Quests map[int]*NQuestInfo
	Items map[int]*NItemInfo
	Skills map[int]*NSkillInfo

}

func (this *NCharacterInfo) InitData()  {
	if this.Bag == nil {
		this.Bag = new(NBagInfo)
	}
	if this.Guild == nil {
		this.Guild = new(NGuildInfo)
	}
	if this.Entity == nil {
		this.Entity = new(NEntity)
		this.Entity.InitData()
	}
	if this.AttrDynamic == nil {
		this.AttrDynamic = new(NAttributeDynamic)
	}

	if this.Items == nil {
		this.Items = make(map[int]*NItemInfo)
	}

}

type NGuildInfo  struct {
	Id int
	GuildName string
	LeaderId int
	LeaderName string
	Notice string
	MemberCount int
	Members map[int]*NGuildMemberInfo
	Applies map[int]*NGuildApplyInfo
	CreateTime long
}

type NFriendInfo struct {
	Id int
	friendInfo *NCharacterInfo
	Status int
}

type NQuestInfo struct {
	QuestId int
	QuestGuid int
	Status QuestStatus
	Targets []int
}

type NBagInfo struct {
	Unlocked int
	Items []byte
}

type NItemInfo struct {
	Id int
	Count int
}

type NEntity struct {
	Id int
	Position *NVector3
	Direction *NVector3
	Speed int
}

func (this *NEntity) InitData()  {
	if this.Position == nil {
		this.Position = new(NVector3)
	}
	if this.Direction == nil {
		this.Direction = new(NVector3)
	}
}


type NAttributeDynamic struct {
	HP int
	MP int
}

type NSkillInfo struct {
	Id int
	Level int
}

type NVector3 struct {
	X int
	Y int
	Z int
}

type NGuildMemberInfo struct {
	Id int
	CharacterId int
	Title GuildTitle
	Info *NCharacterInfo
	JoinTime long
	LastTime long
	Status int
}

func (this*NGuildMemberInfo) InitData()  {
	if this.Info == nil {
		this.Info = new(NCharacterInfo)
	}
}


type NGuildApplyInfo struct {
	GuildId int
	CharacterId int
	Name string
	Class int
	Level int
	Result ApplyResult
}
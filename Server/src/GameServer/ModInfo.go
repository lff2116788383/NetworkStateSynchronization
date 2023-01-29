package GameServer

//
//type long = int64
//
//type CharacterClass = int
//type CharacterType = int
//type QuestStatus = int
//type GuildTitle = int
//type ApplyResult = int
//
////const (
////		None = 0
////		Accept = 1
////		Reject = 2
////)
//
//const (
//	CHARACTER_TYPE_PLAYER = 0
//	CHARACTER_TYPE_NPC = 1
//	CHARACTER_TYPE_MONSTER = 2
//)
//
//type UserInfo struct {
//	Id int
//	Player *PlayerInfo
//}
//
//func (this *UserInfo) InitData()  {
//	if this.Player == nil {
//		this.Player = new(PlayerInfo)
//		this.Player.InitData()
//	}
//}
//
//type PlayerInfo struct {
//	Id int
//	Characters	map[int]*CharacterInfo
//}
//
//func (this *PlayerInfo) InitData()  {
//	if this.Characters == nil {
//		this.Characters = make(map[int]*CharacterInfo)
//	}
//}
//
//type CharacterInfo struct {
//	Id int
//	Name string
//	Level int
//	Exp long
//	Gold long
//	MapId int
//	Ride int
//
//	Class CharacterClass	//角色职业	如弓箭手 战士 法师
//	Type  CharacterType		//角色类型	如玩家	怪物		NPC
//
//
//	EntityId int
//	ConfigId int
//
//	Equips []byte
//
//	Bag	*BagInfo
//	Guild *GuildInfo
//
//	Entity *Entity
//	AttrDynamic *AttributeDynamic
//
//
//	Friends map[int]*FriendInfo
//	Quests map[int]*QuestInfo
//	Items map[int]*ItemInfo
//	Skills map[int]*SkillInfo
//
//}
//
//func (this *CharacterInfo) InitData()  {
//	if this.Bag == nil {
//		this.Bag = new(BagInfo)
//	}
//	if this.Guild == nil {
//		this.Guild = new(GuildInfo)
//	}
//	if this.Entity == nil {
//		this.Entity = new(Entity)
//		this.Entity.InitData()
//	}
//	if this.AttrDynamic == nil {
//		this.AttrDynamic = new(AttributeDynamic)
//	}
//
//	if this.Items == nil {
//		this.Items = make(map[int]*ItemInfo)
//	}
//
//}
//
//type GuildInfo  struct {
//	Id int
//	GuildName string
//	LeaderId int
//	LeaderName string
//	Notice string
//	MemberCount int
//	Members map[int]*GuildMemberInfo
//	Applies map[int]*GuildApplyInfo
//	CreateTime long
//}
//
//type FriendInfo struct {
//	Id int
//	friendInfo *CharacterInfo
//	Status int
//}
//
//type QuestInfo struct {
//	QuestId int
//	QuestGuid int
//	Status QuestStatus
//	Targets []int
//}
//
//type BagInfo struct {
//	Unlocked int
//	Items []byte
//}
//
//type ItemInfo struct {
//	Id int
//	Count int
//}
//
//type Entity struct {
//	Id int
//	Position *Vector3
//	Direction *Vector3
//	Speed int
//}
//
//func (this *Entity) InitData()  {
//	if this.Position == nil {
//		this.Position = new(Vector3)
//	}
//	if this.Direction == nil {
//		this.Direction = new(Vector3)
//	}
//}
//
//
//type AttributeDynamic struct {
//	HP int
//	MP int
//}
//
//type SkillInfo struct {
//	Id int
//	Level int
//}
//
//type Vector3 struct {
//	X int
//	Y int
//	Z int
//}
//
//type GuildMemberInfo struct {
//	Id int
//	CharacterId int
//	Title GuildTitle
//	Info *CharacterInfo
//	JoinTime long
//	LastTime long
//	Status int
//}
//
//func (this*GuildMemberInfo) InitData()  {
//	if this.Info == nil {
//		this.Info = new(CharacterInfo)
//	}
//}
//
//
//type GuildApplyInfo struct {
//	GuildId int
//	CharacterId int
//	Name string
//	Class int
//	Level int
//	Result ApplyResult
//}

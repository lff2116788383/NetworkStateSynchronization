package packet

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
)

const (
	//主版本号
	SERVER_PACEKTVER_NORMAL = 2 //常规版本
	SERVER_PACEKTVER_MERGE  = 3 //国内外命令字合并使用版本号。

	//子版本号
	SERVER_SUBPACKETVER_FIRST  = 1
	SERVER_SUBPACKETVER_NORMAL = 2
	SERVER_SUBPACKETVER_ZIP    = 3 // 压缩版协议ZIP
	SERVER_SUBPACKETVER_FZIP   = 4 // 完全压缩协议 收发都经过zip压缩
)

var (
	FlagsType_None       byte = 0
	FlagsType_HeadPacket byte = 1
	FlagsType_HeartBeat  byte = 2
	FlagsType_BPTPacket  byte = 3
	FlagsType_ICPacket   byte = 4
	FlagsType_Max        byte = 5
)

/**********************数据包interface**********************/
type IPacket interface {
	GetTotalLen() int32
	GetBodyLen() int32
	GetHeadLen() int32

	GetData() []byte
	GetBody() []byte
	Copy(data []byte)
	Refer(data []byte)

	ReadByte() byte
	ReadAll() []byte
	ReadString() string
	ReadInt16() int16
	ReadInt32() int32
	ReadInt64() int64
	ReadInt16B() int16
	ReadInt32B() int32
	ReadInt64B() int64

	WriteByte(value byte)
	WriteBytes(value []byte)
	WriteString(value string)
	WriteInt16(value int16)
	WriteInt32(value int32)
	WriteInt64(value int64)
	WriteInt16B(value int16)
	WriteInt32B(value int32)
	WriteInt64B(value int64)
}

/**********************基本数据包的属性和方法实现**********************/
type BasePacket struct {
	data         []byte //报数据缓冲区
	index        int32  //写数据索引游标
	headLen      int32  //包头长度
	bodyLenIndex int32  //包体长度字段起始位置，默认2字节
}

func (p *BasePacket) GetBodylLenIndex() int32 {
	return p.bodyLenIndex
}
func (p *BasePacket) GetTotalLen() int32 {
	return p.headLen + p.GetBodyLen()
}

func (p *BasePacket) GetData() []byte {
	return p.data
}

func (p *BasePacket) GetBody() []byte {
	return p.data[p.headLen:]
}

func (p *BasePacket) GetBodyLen() int32 {
	fmt.Println("data:", hex.EncodeToString(p.GetData()))
	fmt.Printf("GetBodyLen two byte:%#02x %#02x\n", p.data[p.bodyLenIndex], p.data[p.bodyLenIndex+1])
	fmt.Println("GetBodyLen:", int32(binary.LittleEndian.Uint16(p.data[p.bodyLenIndex:])))

	return int32(binary.LittleEndian.Uint16(p.data[p.bodyLenIndex:]))
}

func (p *BasePacket) GetHeadLen() int32 {
	return p.headLen
}

func (p *BasePacket) Copy(data []byte) {
	p.data = append(p.data[0:], data...)
	p.index += p.GetHeadLen()
}

func (p *BasePacket) Refer(data []byte) {
	p.data = data
	p.index = p.GetHeadLen()
}

/**********************基本数据读取**********************/
func (p *BasePacket) ReadByte() byte {
	var value byte = byte(0)
	if p.index < int32(len(p.data)) {
		value = p.data[p.index]
		p.index++
	}
	return value
}

//读取所有(除了包头)不改变index
func (p *BasePacket) ReadAll() []byte {
	var value []byte
	if p.index < int32(len(p.data)) {
		value = p.data[p.index:]
	}
	return value
}

func (p *BasePacket) ReadString() string {
	value := string("")
	if p.index+4 <= int32(len(p.data)) {
		strLen := int32(binary.LittleEndian.Uint32(p.data[p.index:]))
		p.index += 4
		if p.index+strLen <= int32(len(p.data)) {
			//包中的字符串是'\0'结尾，在C/C++中不会有任何问题，但是golang的string内部
			//都是字节序，'\0'并无特殊对待，作为正常的字符处理，这里需要剔除末尾的'\0'
			if 0 == p.data[p.index+strLen-1] {
				value = string(p.data[p.index : p.index+strLen-1])
			} else {
				value = string(p.data[p.index : p.index+strLen])
			}
			p.index += strLen
		}
	}
	return value
}

//小端读取
func (p *BasePacket) ReadInt16() int16 {
	var value int16 = -1
	if p.index+2 <= int32(len(p.data)) {
		value = int16(binary.LittleEndian.Uint16(p.data[p.index:]))
		p.index += 2
	}
	return value
}

func (p *BasePacket) ReadInt() int {
	return int(p.ReadInt32())
}
func (p *BasePacket) ReadInt32() int32 {
	var value int32 = -1
	if p.index+4 <= int32(len(p.data)) {
		value = int32(binary.LittleEndian.Uint32(p.data[p.index:]))
		p.index += 4
	}
	return value
}

func (p *BasePacket) ReadInt64() int64 {
	var value int64 = -1
	if p.index+8 <= int32(len(p.data)) {
		value = int64(binary.LittleEndian.Uint64(p.data[p.index:]))
		p.index += 8
	}
	return value
}
func (p *BasePacket) ReadFloat32() float32 {
	value, _ := strconv.ParseFloat(p.ReadString(), 32)
	return float32(value)
}

func (p *BasePacket) ReadFloat64() float64 {
	value, _ := strconv.ParseFloat(p.ReadString(), 64)
	return float64(value)
}

//大端读取
func (p *BasePacket) ReadInt16B() int16 {
	var value int16 = -1
	if p.index+2 <= int32(len(p.data)) {
		value = int16(binary.BigEndian.Uint16(p.data[p.index:]))
		p.index += 2
	}
	return value
}

func (p *BasePacket) ReadInt32B() int32 {
	var value int32 = -1
	if p.index+4 <= int32(len(p.data)) {
		value = int32(binary.BigEndian.Uint32(p.data[p.index:]))
		p.index += 4
	}
	return value
}

func (p *BasePacket) ReadInt64B() int64 {
	var value int64 = -1
	if p.index+8 <= int32(len(p.data)) {
		value = int64(binary.BigEndian.Uint64(p.data[p.index:]))
		p.index += 8
	}
	return value
}

/**********************基本数据写入**********************/
func (p *BasePacket) WriteByte(value byte) {
	p.data = append(p.data, value)
}

func (p *BasePacket) WriteBytes(value []byte) {
	p.data = append(p.data, value...)
}

func (p *BasePacket) WriteString(value string) {
	s := []byte(value)
	s = append(s, 0)
	p.WriteInt32(int32(len(s)))
	p.WriteBytes(s)
}

//小端写入
func (p *BasePacket) WriteInt16(value int16) {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf[0:], uint16(value))
	p.data = append(p.data, buf...)
}

func (p *BasePacket) WriteInt32Ex(value int) {
	p.WriteInt32(int32(value))
}

func (p *BasePacket) WriteInt32(value int32) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf[0:], uint32(value))
	p.data = append(p.data, buf...)
}

func (p *BasePacket) WriteInt64(value int64) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf[0:], uint64(value))
	p.data = append(p.data, buf...)
}

//大端写入
func (p *BasePacket) WriteInt16B(value int16) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf[0:], uint16(value))
	p.data = append(p.data, buf...)
}

func (p *BasePacket) WriteInt32B(value int32) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf[0:], uint32(value))
	p.data = append(p.data, buf...)
}

func (p *BasePacket) WriteInt64B(value int64) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf[0:], uint64(value))
	p.data = append(p.data, buf...)
}
func (p *BasePacket) WriteFloat32(value float32) {
	p.WriteString(strconv.FormatFloat((float64)(value), 'E', -1, 32))
}
func (p *BasePacket) WriteFloat64(value float64) {
	p.WriteString(strconv.FormatFloat((float64)(value), 'E', -1, 64))
}

/**********************具体业务包**********************/

const (
	IC_PACKET_HEADER_LENGTH           = 13
	MONEY_PACKET_HEADER_LENGTH        = 9
	CMS_PACKET_HEADER_LENGTH          = 10
	BPT_MAX_PACKET_LENGTH             = 2 * 1024
	BPT_EXTERNAL_PACKET_HEADER_LENGTH = 12
	BPT_INTERNAL_PACKET_HEADER_LENGTH = 22
	LOG_PACKET_HEADER_LENGTH          = 6
	LOG_MAX_PACKET_LENGTH             = 10 * 1024
	SLT_PACKET_HEADER_LENGTH          = 16
	MAX_USER_PACKET_LEN               = 20 * 1024 // 业务层包体最大20KB
)

/**********************IC包相关实现**********************/
//struct ICExternalHeader {
//	char magic[2];
//	unsigned short cmd;
//	unsigned char ver;
//	unsigned char subver;
//	unsigned short bodylen;
//	char checksum;
//	unsigned int sequence;
//	char data[0];
//};

type ICPacket struct {
	BasePacket
	isEncrypted bool
}

func NewICPacket() *ICPacket {
	p := new(ICPacket)

	p.headLen = IC_PACKET_HEADER_LENGTH
	p.bodyLenIndex = 6
	p.data = make([]byte, 0, BPT_MAX_PACKET_LENGTH)

	return p
}

/**********************IC包加解密**********************/
func (p *ICPacket) Encrypt() int {
	return encryptPacket(p)
}

func (p *ICPacket) Decrypt() int {
	return decryptPacket(p)
}

func (p *ICPacket) SetBegin(cmd int16, version byte, subVersion byte) {
	p.WriteBytes([]byte("IC"))
	p.WriteInt16(cmd)
	p.WriteByte(version)
	p.WriteByte(subVersion)
	p.WriteInt16(0)
	p.WriteByte(0)
	p.WriteInt32(0)
}

// 默认写版本号
func (p *ICPacket) Begin(cmd int16) {
	p.WriteBytes([]byte("IC"))
	p.WriteInt16(cmd)
	p.WriteByte(SERVER_PACEKTVER_NORMAL)
	p.WriteByte(SERVER_SUBPACKETVER_FIRST)
	p.WriteInt16(0)
	p.WriteByte(0)
	p.WriteInt32(0)
}

func (p *ICPacket) SetEnd() {
	bodyLen := int16(len(p.data) - int(p.GetHeadLen()))
	binary.LittleEndian.PutUint16(p.data[6:], uint16(bodyLen))
}

//默认加密
func (p *ICPacket) End() {
	fmt.Println("data test1:", hex.EncodeToString(p.data))

	bodyLen := int16(len(p.data) - int(p.GetHeadLen()))
	binary.LittleEndian.PutUint16(p.data[6:], uint16(bodyLen))

	fmt.Println("data test2:", hex.EncodeToString(p.data))

	//加密处理
	p.Encrypt()
}

func (p *ICPacket) SetCheckCode(code byte) {
	p.data[8] = code
	p.isEncrypted = true
}

func (p *ICPacket) GetCheckCode() byte {
	return p.data[8]
}

func (p *ICPacket) IsEncrypted() bool {
	return p.isEncrypted
}

func (p *ICPacket) GetCmd() uint16 {
	return binary.LittleEndian.Uint16(p.data[2:])
}

func (p *ICPacket) GetVersion() byte {
	return p.data[4]
}

func (p *ICPacket) GetSubVersion() byte {
	return p.data[5]
}

func (p *ICPacket) GetSequence() uint32 {
	return binary.LittleEndian.Uint32(p.data[9:])
}

//SLT包定义
/*
|  3  |  1   |    2     |  2  |  4  |    4     |   ...   |
| SLT | type | body len | cmd | uid | reserved | content |
*/
type SltPacket struct {
	BasePacket
}

func NewSltPacket() *SltPacket {
	p := new(SltPacket)
	p.headLen = SLT_PACKET_HEADER_LENGTH
	p.bodyLenIndex = 4
	p.data = make([]byte, 0, BPT_MAX_PACKET_LENGTH)
	return p
}
func (p *SltPacket) Begin(linkId uint64, cmd int16) {
	p.WriteBytes([]byte("SLT"))
	p.WriteByte(0)
	p.WriteInt16(0)
	p.WriteInt16(cmd)
	p.WriteInt64(int64(linkId))
}
func (p *SltPacket) End() {
	bodyLen := int16(len(p.data) - int(p.headLen))
	binary.LittleEndian.PutUint16(p.data[p.bodyLenIndex:], uint16(bodyLen))
}
func (p *SltPacket) GetCmdType() byte {
	return p.data[3]
}
func (p *SltPacket) GetCmd() uint16 {
	return binary.LittleEndian.Uint16(p.data[6:])
}
func (p *SltPacket) GetLinkId() uint64 {
	return binary.LittleEndian.Uint64(p.data[8:])
}

//金币server的通信包
type MoneyPacket struct {
	ICPacket
}

func NewMoneyPacket() *MoneyPacket {
	p := new(MoneyPacket)
	p.headLen = MONEY_PACKET_HEADER_LENGTH
	p.bodyLenIndex = 6
	p.data = make([]byte, 0, BPT_MAX_PACKET_LENGTH)
	return p
}

//overwrite
func (p *MoneyPacket) SetBegin(cmd int16, version byte, subVersion byte) {
	p.WriteBytes([]byte("IC"))
	p.WriteInt16(cmd)
	p.WriteByte(version)
	p.WriteByte(subVersion)
	p.WriteInt16(0)
	p.WriteByte(0)
}

/**********************BPT内部包相关实现**********************/
/**
 * BPT包结构:
 * | 2 byte | 4 byte | 4 byte | 8 byte |  2 byte  |
 * |  cmd   |   ver  | subver |  svid  |    0     |
 */
//struct BPTInternalHeader {
//	unsigned short cmdtype;
//	unsigned int maincmd;
//	unsigned int subcmd;
//	unsigned long long linkid;
//	unsigned short varhrlen;
//	unsigned short datalen;
//	char data[0];
//};

type BPTInternalPacket struct {
	BasePacket
}

func NewBPTInternalPacket() *BPTInternalPacket {
	p := new(BPTInternalPacket)

	p.headLen = BPT_INTERNAL_PACKET_HEADER_LENGTH
	p.bodyLenIndex = 20
	p.data = make([]byte, 0, BPT_MAX_PACKET_LENGTH)

	return p
}

func (p *BPTInternalPacket) Encode(
	cmdtype uint16,
	maincmd uint32,
	subcmd uint32,
	linkid uint64,
	varhrlen uint16,
	data []byte) bool {

	p.WriteInt16(int16(cmdtype))
	p.WriteInt32(int32(maincmd))
	p.WriteInt32(int32(subcmd))
	p.WriteInt64(int64(linkid))
	p.WriteInt16(int16(varhrlen))
	p.WriteInt16(int16(len(data)))

	p.WriteBytes(data)
	return p.GetHeadLen()+int32(len(data)) <= p.GetTotalLen()
}

func (p *BPTInternalPacket) GetCmdType() uint16 {
	return binary.LittleEndian.Uint16(p.data[0:])
}

func (p *BPTInternalPacket) GetMainCmd() uint16 {
	return uint16(binary.LittleEndian.Uint16(p.data[2:]))
}

func (p *BPTInternalPacket) GetSubCmd() uint16 {
	return uint16(binary.LittleEndian.Uint16(p.data[6:]))
}

func (p *BPTInternalPacket) GetLinkid() uint64 {
	return uint64(binary.LittleEndian.Uint64(p.data[10:]))
}

/**********************BPT外部包相关实现**********************/
//struct BPTExternalHeader {
//	char magic[3];
//	unsigned char headlen;
//	unsigned short bodylen;
//	unsigned short maincmd;
//	unsigned short subcmd;
//	char cipherverion;
//	char version;
//	char data[0];
//};

type BPTExternalPacket struct {
	BasePacket
}

func NewBPTExternalPacket() *BPTExternalPacket {
	p := new(BPTExternalPacket)

	p.headLen = BPT_EXTERNAL_PACKET_HEADER_LENGTH
	p.bodyLenIndex = 4
	p.data = make([]byte, 0, BPT_MAX_PACKET_LENGTH)

	return p
}

//构造BPTExternalPacket
func (p *BPTExternalPacket) Encode(maincmd uint16, subcmd uint16, data []byte) bool {
	p.WriteBytes([]byte("BPT"))
	p.WriteByte(BPT_EXTERNAL_PACKET_HEADER_LENGTH)
	p.WriteInt16B(int16(len(data)))
	p.WriteInt16B(int16(maincmd))
	p.WriteInt16B(int16(subcmd))
	p.WriteByte(1)
	p.WriteByte(1)

	p.WriteBytes(data)

	return p.GetHeadLen()+int32(len(data)) <= p.GetTotalLen()
}

func (p *BPTExternalPacket) GetMainCmd() uint16 {
	return uint16(binary.BigEndian.Uint16(p.data[6:]))
}

func (p *BPTExternalPacket) GetSubCmd() uint16 {
	return uint16(binary.BigEndian.Uint16(p.data[8:]))
}

func (p *BPTExternalPacket) GetBodyLen() int32 {
	return int32(binary.BigEndian.Uint16(p.data[p.bodyLenIndex:]))
}

/**********************CMS数据包相关实现**********************/
//struct CMSICInterHeader {
//		uint64 context;
//		uint16 dlen;
//		char data[0];
//	};
type CMSPacket struct {
	BasePacket
}

func NewCMSPacket() *CMSPacket {
	p := new(CMSPacket)

	p.headLen = CMS_PACKET_HEADER_LENGTH
	p.bodyLenIndex = 8
	p.data = make([]byte, 0, BPT_MAX_PACKET_LENGTH)

	return p
}

func (p *CMSPacket) Encode(content uint64, data []byte) {
	p.WriteInt64(int64(content))
	p.WriteInt16(int16(len(data)))
	p.WriteBytes(data)
}

func (p *CMSPacket) GetContext() uint64 {
	return uint64(binary.LittleEndian.Uint64(p.data[0:]))
}

/**********************Log数据包相关实现**********************/
/**
 * Log包结构:
 * | 4 bytes |  1 byte   |     1 byte      | ServerName |  data  |
 * |  包长度   'C' 或 'R'   ServerName长度      "Name"      数据  |
 */
type LogPacket struct {
	BasePacket
	ServerName string
}

func NewLogPacket() *LogPacket {
	p := new(LogPacket)
	p.headLen = LOG_PACKET_HEADER_LENGTH
	return p
}

func (p *LogPacket) GetTotalLen() int32 {
	return int32(binary.LittleEndian.Uint32(p.data[0:]))
}

func (p *LogPacket) GetBodyLen() int32 {
	return p.GetTotalLen() - p.headLen
}

func (p *LogPacket) Begin(serverName string) {
	p.WriteInt32(0)
	p.WriteByte(byte('C'))
	p.WriteByte(byte(len(serverName)))
	p.WriteBytes([]byte(serverName))
	p.ServerName = serverName
}

func (p *LogPacket) End() {
	packageLen := int32(len(p.data))
	binary.LittleEndian.PutUint32(p.data[0:], uint32(packageLen))
}

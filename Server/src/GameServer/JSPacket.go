package GameServer

import (
	"Src/libs/logger"
	"encoding/json"
)
const(
	JS_MAX_PACKET_LENGTH             = 2 * 1024
)
type JSPacket struct {
	JSPMap map[string]interface{}
	data []byte
}

func (j* JSPacket) Copy(buf []byte)  {
	j.data = buf
	j.Decrypt()
	err := json.Unmarshal(j.data,&j.JSPMap)
	if nil != err {
		logger.Error("json.Unmarshal err[%s]", err)
		return
	}
}

func (j* JSPacket) ReadInt()  {

}

func (j* JSPacket) Encrypt()  {
	var i byte =0
	for	;i< (byte)(len(j.data));i++ {
		j.data[i] = (byte)(j.data[i] + i +5)
	}
}
func (j* JSPacket) Decrypt()  {
	var i byte =0
	for	;i< (byte)(len(j.data));i++ {
		j.data[i] = (byte)(j.data[i] - i -5)
	}
}

func NewJSPacket() *JSPacket  {
	jsp:=&JSPacket{
		data: make([]byte,0,JS_MAX_PACKET_LENGTH),
		JSPMap: make(map[string]interface{}),
	}
	return jsp
}


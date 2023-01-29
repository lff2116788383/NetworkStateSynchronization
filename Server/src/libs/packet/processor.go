package packet

import (
	"Src/libs/logger"
	"sync"
	"time"
)

var (
	p    *MsgProcessor // 单例
	once sync.Once
)

const (
	MAX_SEND_TIMEOUT = 1
)

type Processor interface {
	Process(Message)
	Epoll() Message
}

type MsgProcessor struct {
	messageChannel chan Message
}

func GetMsgProcessor() *MsgProcessor {
	once.Do(func() {
		p = NewMsgProcessor()
	})
	return p
}

func NewMsgProcessor() *MsgProcessor {
	return &MsgProcessor{make(chan Message, 10000)}
}

type Message interface {
	Data() interface{} //消息的数据
	Args() interface{} //消息的参数
}

var LenMessageChannel = 0

func (p *MsgProcessor) Process(message Message) {

	select {
	case p.messageChannel <- message:
		LenMessageChannel = len(p.messageChannel)
	case <-time.After(time.Second * MAX_SEND_TIMEOUT):
		//超时，导致数据丢弃
		logger.Error("send to packet channel timeout!!! p.messageChannel.len[%d]", len(p.messageChannel))
	}
}

func (p *MsgProcessor) Epoll() Message {
	return <-p.messageChannel
}

package mmo

import (
	"Src/libs/logger"
	"Src/libs/packet"
	"io"
	"net"


)

//泛型
const (
	LISTEN_TYPE_CLIENT = 1
	LISTEN_TYPE_SERVER = 2
	LISTEN_TYPE_ADMIN  = 3
)

//消息处理使用泛型 方便针对不同的连接对象做处理(如GM账号和普通账号)
type Processer interface {
	GetType() int
	GetConn() net.Conn
	GenPacket() packet.IPacket
	ProcessPkg(packet.IPacket)
}

func ReceiveIcPkg(prs Processer) {
	conn := prs.GetConn()
	defer func() {
		conn.Close()
		if prs.GetType() == LISTEN_TYPE_CLIENT {

			//连接关闭 先从table删除 再从大厅删除  最后从在线连接map删除
			table:=GetGameServer().GetTable(prs.(*User).Tid)
			if table != nil {
				table.DelUser(prs.(*User).Id)
			}

			GetGameServer().GetTable(MAP_TYPE_HALL).DelUser(prs.(*User).Id)

			prs.(*User).Offline()

		}
	}()

	for {
		if conn == nil {
			return
		}
		p := prs.GenPacket()

		head := make([]byte, p.GetHeadLen())
		_, err := io.ReadFull(conn, head)
		if nil != err {
			logger.Error("br.Read 1 err[%s]", err)
			return
		}
		//logger.Error("head len:[%d]", len(head))

		p.Copy(head)
		//logger.Error("head Data:[%s]", hex.EncodeToString(head))

		bodyLen := p.GetBodyLen()
		if bodyLen >= packet.MAX_USER_PACKET_LEN {
			logger.Error("length of uesr packet more than MAX_USER_PACKET_LEN, bodyLen[%d]", bodyLen)
			return
		}
		//logger.Error("bodyData len:[%d]", bodyLen)

		bodyData := make([]byte, bodyLen)
		_, err = io.ReadFull(conn, bodyData)
		if nil != err {
			logger.Error("io.ReadFull(%d) failed, error[%s]", bodyLen, err)
			return
		}

		//logger.Error("bodyData:[%s]", hex.EncodeToString(bodyData))

		p.WriteBytes(bodyData)

		prs.ProcessPkg(p)
	}
}

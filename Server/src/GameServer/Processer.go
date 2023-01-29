package GameServer

import (
	"Src/libs/logger"
	"Src/libs/packet"
	"encoding/json"
	"net"
	"fmt"
	"os"
)
//消息处理使用泛型 方便针对不同的连接对象做处理(如GM账号和普通账号)
type Processer interface {
	GetType() int
	GetConn() net.Conn
	GenPacket() packet.IPacket
	ProcessPkg(packet.IPacket)

	ProcessJsPkg(map[string]interface{})
}

func ReceiveIcPkg(prs Processer) {
	conn := prs.GetConn()
	defer func() {
		conn.Close()
		if prs.GetType() == LISTEN_TYPE_CLIENT { //处理客户端连接断开
			GetGameServer().ProcessClose(prs.(*SocketHandler))//类型断言	接口转换为具体类型
		}
	}()

	//for {
	//	if conn == nil {
	//		return
	//	}
	//	p := prs.GenPacket()
	//
	//	head := make([]byte, p.GetHeadLen())
	//	_, err := io.ReadFull(conn, head)
	//	if nil != err {
	//		logger.Error("br.Read 1 err[%s]", err)
	//		return
	//	}
	//	//logger.Error("head len:[%d]", len(head))
	//
	//	p.Copy(head)
	//	//logger.Error("head Data:[%s]", hex.EncodeToString(head))
	//
	//	bodyLen := p.GetBodyLen()
	//	if bodyLen >= packet.MAX_USER_PACKET_LEN {
	//		logger.Error("length of uesr packet more than MAX_USER_PACKET_LEN, bodyLen[%d]", bodyLen)
	//		return
	//	}
	//	//logger.Error("bodyData len:[%d]", bodyLen)
	//
	//	bodyData := make([]byte, bodyLen)
	//	_, err = io.ReadFull(conn, bodyData)
	//	if nil != err {
	//		logger.Error("io.ReadFull(%d) failed, error[%s]", bodyLen, err)
	//		return
	//	}
	//
	//	//logger.Error("bodyData:[%s]", hex.EncodeToString(bodyData))
	//
	//	p.WriteBytes(bodyData)
	//
	//	prs.ProcessPkg(p)
	//}

	// for{
	// 	if conn == nil {
	// 		return
	// 	}
	// 	buf := make([]byte,1024)
	// 	//_, err := io.ReadFull(conn, buf)
	// 	readLen, err := conn.Read(buf)
	// 	if nil != err {
	// 		logger.Error("br.Read 1 err[%s]", err)
	// 		return
	// 	}
	// 	//logger.Error("recv json data: %s", string(buf))
	// 	logger.Error("Unmarshal json data: %s", string(buf[:readLen]))
		

	// 	words := make(map[string]interface{})
	// 	err = json.Unmarshal(buf[:readLen],&words)
	// 	if nil != err {
	// 		logger.Error("json.Unmarshal err[%s]", err)
	// 		return
	// 	}
	// 	// for k, v := range words {
	// 	// 	logger.Error("key: [%s], value: [%s]", k,v)
	// 	// }

	// 	prs.ProcessJsPkg(words)

		
	// }


	//connRemoteAddr := conn.RemoteAddr().String()
	//// 循环读取客户端的消息
	//for {
	//	// 创建一个1024长度的 byte 切片用于存储客户端的消息
	//	bytes := make([]byte, 1024)
	//	// 等待客户端通过 conn 发送消息
	//	// 如果客户端一直没有发送，那么此协程就阻塞在这里
	//	readLen, err := conn.Read(bytes)
	//	if err != nil {
	//		errStr := err.Error()
	//		contains := strings.Contains(errStr, "An existing connection was forcibly closed by the remote host")
	//		if contains || errStr == "EOF" {
	//			fmt.Printf("客户端[%v]: %v 断开链接\n", nowTime(), connRemoteAddr)
	//			return
	//		} else {
	//			fmt.Printf("服务器读取客户端[%v]: %v 消息出现异常, err = %v \n", nowTime(), connRemoteAddr, err)
	//			return
	//		}
	//	}
	//	// 服务器端显示客户端的消息
	//	// bytes[:readLen] ==> 只打印有效的数据长度
	//	fmt.Printf("server[%v]-客户端[%v]说: %v\n", nowTime(), connRemoteAddr, string(bytes[:readLen]))
	//}



	//声明一个临时缓冲区，用来存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)
	go reader(readerChannel)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		tmpBuffer = Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
		for i := 0; i < len(tmpBuffer); i++ {
			
			if tmpBuffer[i] == '{' {
				for j := i; j < len(tmpBuffer); j++{
					if tmpBuffer[j] == '}' {
						logger.Error("Unmarshal json data: %s", string(tmpBuffer[i:j+1]))
						words := make(map[string]interface{})
						err = json.Unmarshal(tmpBuffer[i:j+1],&words)
						if nil != err {
							logger.Error("json.Unmarshal err[%s]", err)
							return
						}
						tmpBuffer= tmpBuffer[j+1:]
						prs.ProcessJsPkg(words)
						break;
					}
				}
				
			}
		}
		

	}
}


func reader(readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			Log(string(data))
		}
	}
}

func Log(v ...interface{}) {
	fmt.Println(v...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

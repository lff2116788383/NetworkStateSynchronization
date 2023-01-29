package GameServer

import (
	"bytes"
	"encoding/binary"
)

var DATAPACK_HEAD int32 = 0xffff

type DataPack struct {
	//数据缓冲区起始地址
	m_starptr []byte
	//数据缓冲区结束地址
	m_endptr []byte
	//数据缓冲区中数据的结束地址
	m_dataptr []byte
	//读取数据的偏移指针
	m_offset []byte
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

func NewDataPack(buffer []byte, length int )  {
	pack:=new(DataPack)
	if buffer == nil {
		pack.m_starptr = make([]byte,0,128)
		pack.m_endptr = append(pack.m_starptr, 128)
		pack.m_dataptr = pack.m_starptr
		pack.m_offset = pack.m_starptr
		pack.WriteInt32(DATAPACK_HEAD)
		pack.WriteInt32(4)
		pack.m_offset = append(pack.m_starptr, 6)
	}else {
		if length == -1 {
			length = len(buffer)
		}
		pack.m_starptr = make([]byte,0,length)
		copy(pack.m_starptr, buffer[:length])
		pack.m_endptr = append(pack.m_starptr, (byte)(length))
		pack.m_dataptr = append(pack.m_starptr, (byte)(length))
		pack.m_offset = append(pack.m_starptr, 6)
	}

}
func (d *DataPack) SetSize(size int)  {
	//保存现在的各个地址位置偏移大小
	//oldlen := len(d.m_endptr) - len(d.m_starptr)
	datalen := len(d.m_dataptr) - len(d.m_starptr)
	offset := len(d.m_offset) - len(d.m_starptr)

	//动态扩充数据缓冲区内存并还原各个地址相对m_starptr偏移位置
	d.m_starptr = make([]byte,0, len(d.m_starptr)+size)
	d.m_endptr = append(d.m_starptr, byte(size))
	d.m_dataptr = append(d.m_starptr, byte(datalen))
	d.m_offset = append(d.m_starptr, byte(offset))
}

func (d *DataPack) GetBuffer()[]byte  {
	return d.m_starptr
}
func (d *DataPack) GetBufferLength() int {
	return len(d.m_dataptr)- len(d.m_starptr)
}
func (d *DataPack) SetOffset(pos int)  {
	d.m_offset = append(d.m_starptr, byte(pos))
}

func (d *DataPack) WriteBuffer(buff []byte,size int)  {

}

func (d *DataPack) WriteString(buff []byte,length int)  {

}

func (d *DataPack) WritePackLen(length int)  {

}
func (d *DataPack) CheckCapacity(size int)  {

}
func (d *DataPack) WriteAtomic()  {

}
func (d *DataPack) WriteInt(value int)  {

}

func (d* DataPack) WriteInt32(value int32)  {

}
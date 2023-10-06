package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

//封包，拆包的具体模块
type DataPack struct {
}

//拆包封包实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	//DataLen uint32(4字节) + ID uint32(4字节)
	return 8
}

//封包方法
//datalen|msgID|data
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一bytes字节缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//将msg的各消息写入缓冲
	err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		return nil, err
	}

	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}

	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetDta())
	if err != nil {
		return nil, err
	}

	//dataBuff.Write(msg.GetDta())

	return dataBuff.Bytes(), nil
}

//拆包方法
//将包的Head消息读出，之后再根据head消息里的data的长度，读取data
func (d *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//解压head消息，得到dataLen和MsgId
	msg := &Message{}

	err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}

	err = binary.Read(dataBuff, binary.LittleEndian, &msg.Id)
	if err != nil {
		return nil, err
	}

	//判断dataLen是否超出程序允许的最大包长度
	outRange := utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize
	if outRange {
		return nil, errors.New("out pack Range")
	}

	return msg, nil
}

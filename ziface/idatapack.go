package ziface

/*
 封包、拆包 模块
 直接面向Tcp链接中的数据流，用于处理TCP粘包的问题
*/

type IDataPack interface {
	//获取包的头的长度
	GetHeadLen() uint32
	//封包
	Pack(msg IMessage) ([]byte, error)
	//拆包
	Unpack([]byte) (IMessage, error)
}

package ziface

/*
 IRequest接口
 将客户端请求链接信息与请求的数据包装到一个request中
*/

type IRequest interface {
	//得到当前链接
	GetConnection() IConnection

	//得到请求的消息数据
	GetData() []byte

	//得到请求的消息ID
	GetMsgID() uint32
}

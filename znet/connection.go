package znet

import (
	"errors"
	"io"
	"log"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

/*
 链接模块
*/

type Connection struct {
	//当前Conn隶属于哪个Server
	TcpServer ziface.IServer

	//当前链接的socket TCP套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前的链接状态
	isClose bool

	////当前链接所绑定的处理业务方法api
	//handleAPI ziface.HandleFunc

	//链接监测告知通道 channel
	ExitChan chan bool

	//无缓冲管道，用于读写协程间的消息通信
	msgChan chan []byte

	////链接处理的方法Router
	//Router ziface.IRouter

	//消息的管理MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
}

//链接模块初始化方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandle,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}

	//将conn加入ConnManager中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

//链接的读业务方法
func (c *Connection) StartReader() {
	log.Println(" [Reader Goroutine is running...]")
	defer log.Println("[Reader is exit!] connID= ", c.ConnID, ", remote add is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		////读取客户端的数据到buf中，最大字节
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	log.Println("recv buf err: ", err)
		//	return
		//}

		//创建一个拆包解包对象
		dp := NewDataPack()
		//读取客户端的Msg Head 二级制流 8个字节，
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			log.Println("read msg head error:", err)
			break
		}

		//拆包，得到msgID 和 msgDatalen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			log.Println("unpack error:", err)
			break
		}

		//根据datalen，再次读取Data,放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				log.Println("read msg data error:", err)
				break
			}
		}
		msg.SetData(data)

		////调用当前链接所绑定的业务HandleAPI
		//if err = c.handleAPI(c.Conn, buf, cnt); err != nil {
		//	log.Println("ConnID: ", c.ConnID, " handle is error: ", err)
		//	return
		//}

		//得到当前conn的reque请求数据
		req := Request{
			c,
			msg,
		}

		//执行注册的路由方法
		//go func(request ziface.IRequest) {
		//	c.Router.PreHandle(request)
		//	c.Router.Handle(request)
		//	c.Router.PostHandle(request)
		//}(&req)

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已开启工作池机制，将消息发送worker工作池处理
			go c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

//写消息协程，发送消息客户端模块
func (c *Connection) StartWriter() {
	log.Println("[Writer Goroutine is running...]")
	defer log.Printf("%s [conn Writer exit!]", c.RemoteAddr().String())

	//不断的阻塞等待channel的消息
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				log.Println("Send data error:", err)
				return
			}
		case <-c.ExitChan:
			//代表reader已经退出
			return
		}
	}
}

func (c *Connection) Start() {
	log.Println("Conn Start() ... ConnID= ", c.ConnID)

	//启动当前链接的读写数据的业务
	go c.StartReader()
	go c.StartWriter()

	//调用开发者设置的钩子函数
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	log.Println("Conn Stop()... ConnID= ", c.ConnID)

	//如果当前链接已关闭
	if c.isClose {
		return
	}
	c.isClose = true

	//调用开发者设置的钩子函数
	c.TcpServer.CallOnConnStop(c)

	//关闭socket链接
	c.Conn.Close()

	//告知writer关闭
	c.ExitChan <- true

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)

	//将当前链接从ConnMgr中删除
	c.TcpServer.GetConnMgr().Remove(c)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//提供一个SendMsg方法，将发送给客户端的数据进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClose {
		return errors.New("Connection closed when send msg")
	}

	//将data进行封包 MsgDataLen|MsgID|Data
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		log.Println("Pack error msg id=", msgId)
		return errors.New("Pack msg error")
	}

	//将数据发送给客户端
	//if _, err = c.Conn.Write(binaryMsg); err != nil {
	//	log.Printf("Write msg id:%d err:%v", msgId, err)
	//	return errors.New("conn Write error")
	//}

	//将发送数据发往管道
	c.msgChan <- binaryMsg

	return nil
}

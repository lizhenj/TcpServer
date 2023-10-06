package znet

import (
	"log"
	"net"
	"testing"
)

//只是负责测试datapack拆包 封包的单元测试
func TestDataPack(t *testing.T) {
	/*
	 模拟的服务器
	*/
	//1 创建socketTCP
	listenner, err := net.Listen("tcp", "127.0.0.1:8989")
	if err != nil {
		panic(err)
	}

	//起一协程 负责从客户端处理业务
	go func() {
		//2 从客户端读取数据，拆包处理
		for {
			conn, err := listenner.Accept()
			if err != nil {
				log.Println("server accept error: ", err)
				continue
			}

			go func(conn net.Conn) {
				//处理客户端请求
				//--->拆包过程《<----
				//定义一拆包对象dp
				dp := NewDataPack()
				for {
					//1第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := conn.Read(headData)
					if err != nil {
						log.Println("read head error: ", err)
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						log.Println("server unpacke error: ", err)
						break
					}
					if msgHead.GetMsgLen() > 0 {
						//2第二次从conn读，根据head中的datalen，再读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						_, err := conn.Read(msg.Data)
						if err != nil {
							log.Println("server unpacke data error: ", err)
							break
						}

					}

					//完整一消息已读取完毕
					log.Printf("--->Recv MsgID=%d, datalen=%d, data=%s", msgHead.GetMsgId(),
						msgHead.GetMsgLen(), string(msgHead.GetDta()))
				}
			}(conn)
		}
	}()

	/*
	 模拟客户端
	*/
	conn, err := net.Dial("tcp4", "127.0.0.1:8989")
	if err != nil {
		log.Println("client dial err: ", err)
		return
	}

	//创建一个封包对象 dp
	dp := NewDataPack()

	//模拟粘包过程，封装两个msg一同发送
	msg1 := &Message{
		1,
		5,
		[]byte{'H', 'E', 'L', 'L', 'O'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		log.Println("client pack msg1 error", err)
		return
	}

	msg2 := &Message{
		2,
		4,
		[]byte{'n', 'i', 'h', 'o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		log.Println("client pack msg2 error", err)
		return
	}

	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	//客户端阻塞
	select {}
}

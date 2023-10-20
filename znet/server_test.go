package znet

import (
	"fmt"
	"time"
)

/*
  模拟客户端
*/
func ClientTest() {
	fmt.Println("Client Test ... start")
	//3s后再发起测试请求，等待服务器启动
	time.Sleep(3 * time.Second)

}

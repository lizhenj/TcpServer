package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

/*
 存储Zinx框架的全局参数，供其他模块使用
 一些参数是可以通过zinx.json由用户进行配置
*/

type GlobalObj struct {
	/*
	 Server
	*/
	TcpServer ziface.IServer //当前服务器全局的server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器的名称

	/*
	 Zinx
	*/
	Version          string //当前版本号
	MaxConn          int    //当前服务器主机允许的最大链接数
	MaxPackageSize   uint32 //当前框架数据包的最大值
	WorkerPoolSize   uint32 //当前业务工作池的Goroutine数量
	MaxWorkerTaskLen uint32 //框架允许的开辟worker最大数量
}

/*
 定义一个全局的对外GlobalObj
*/
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("D:/Study/go-work/src/myDemo/ZinxV0.1/conf/zinx.json")
	if err != nil {
		panic(err)
	}

	//将json文件数据解析到struct
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
 init方法，初始化当前的GlobalObject对象
*/
func init() {
	//配置默认值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.5",
		TcpPort:          8989,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	//从用户自定义参数文件，加载配置
	GlobalObject.Reload()
}

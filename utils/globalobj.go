package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"zinx/ziface"
	"zinx/zlog"
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

	/*
	  config file path
	*/
	ConfFilePath string

	/*
	  logger
	*/
	LogDir        string //日志所在文件夹 默认“./log”
	LogFile       string //日志文件名称 默认“” --如果没有设置日志文件，打印信息将打印之stderr
	LogDebugClose bool   //是否关闭Debug日志级别调试信息 默认false --默认打开debug信息
}

/*
 定义一个全局的对外GlobalObj
*/
var GlobalObject *GlobalObj

//判断一个文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//读取用户的配置文件
func (g *GlobalObj) Reload() {

	if confFileExists, _ := PathExists(g.ConfFilePath); confFileExists != true {
		return
	}

	data, err := ioutil.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}

	//将json文件数据解析到struct
	err = json.Unmarshal(data, g)
	if err != nil {
		panic(err)
	}

	//Logger 设置
	if g.LogFile != "" {
		zlog.SetLogFile(g.LogDir, g.LogFile)
	}
	if g.LogDebugClose == true {
		zlog.CloseDebug()
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
		ConfFilePath:     "D:/Study/go-work/src/myDemo/ZinxV0.1/conf/zinx.json",
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		LogDir:           "./log",
		LogFile:          "",
	}

	//从用户自定义参数文件，加载配置
	GlobalObject.Reload()
}

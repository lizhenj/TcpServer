package zlog

import (
	"fmt"
	"testing"
)

func TestStdZLog(t *testing.T) {
	//测试 默认debug输出
	Debug("zinx debug content1")
	Debug("zinx debug content2")

	Debugf("zinx debug a = %d\n", 10)

	//设置log标记位，加上文件名称和微秒标记
	ResetFlags(BitDate | BitLongFile | BitLevel)

	//设置日志前缀，主要标记当前日志模块
	SetPrefix("MODULE")
	Error("zinx error content")

	//添加标记位
	AddFlag(BitShortFile | BitTime)
	Stack("Zinx Stack! ")

	//设置日志写入文件
	SetLogFile("./log", "testfile.log")
	Debug("===> zinx debug content ~~666")
	Debug("===> zinx debug content ~~888")
	Error("===> zinx Error!!!! ~~~555~~~")

	//关闭debug调试
	CloseDebug()
	Debug("====>不会出现~！")
	Debug("====>不会出现~！")
	Error("---> zinx error after debug close ！！！")
}

//------------

//位图
type BitMap struct {
	bits []byte
	max  int
}

//初始化一个BitMap
//一个byte有8位，可代表8个数字，取余后加1为存放最大数所需的容量
func NewBitMap(max int) *BitMap {
	bits := make([]byte, (max>>3 + 1))
	return &BitMap{bits: bits, max: max}
}

//添加一个数字到位图
//计算添加数字在数组中的索引index,一个索引可以存放8个数字
//计算存放到索引下的第几个位置，一共0-7个位置
//原索引下的内容与1左移到指定位置后做或运算
func (b *BitMap) Add(num uint) {
	index := num >> 3
	pos := num % 8
	b.bits[index] |= 1 << pos
}

//判断一个数字是否在位图
//找到数字所在的位置，然后做运算
func (b *BitMap) IsExist(num uint) bool {
	index := num >> 3
	pos := num % 8
	return b.bits[index]&(1<<pos) != 0
}

//删除一个在位图的数字
//找到数字所在的位置取反，然后与索引下的数字做与运算
func (b *BitMap) Remove(num uint) {
	index := num >> 3
	pos := num % 8 //获取在8个比特位中的位置
	b.bits[index] = b.bits[index] & ^(1 << pos)
}

//位图的最大数字
func (b *BitMap) Max() int {
	return b.max
}

func (b *BitMap) String() string {
	return fmt.Sprint(b.bits)
}

func TestBitMap(t *testing.T) {
	b := NewBitMap(100000000)
	b.Add(1000)
	fmt.Println("加入测试：", b.IsExist(1000))
	b.Remove(1000)
	fmt.Println("删去测试：", b.IsExist(1000))
	fmt.Println("最大值：", b.Max())
	fmt.Println("数据：", b.String())
}

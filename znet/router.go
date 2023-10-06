package znet

import "zinx/ziface"

//实现router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类的方法进行重写就好了
type BaseRouter struct{}

//BaseRouter的方法都为空
//router可继承BaseRouter进行需求，个性化重写业务函数

func (b *BaseRouter) PreHandle(request ziface.IRequest) {
}

func (b *BaseRouter) Handle(request ziface.IRequest) {
}

func (b *BaseRouter) PostHandle(request ziface.IRequest) {
}

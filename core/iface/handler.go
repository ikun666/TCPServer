package iface

type IHandler interface {
	PreHandle(IRequest)  //处理请求前hook函数
	Handle(IRequest)     //处理请求
	PostHandle(IRequest) //处理请求后hook函数
}

package iface

type IRouter interface {
	AddRoute(uint32, IHandler) error //添加路由-handle
	DoRoute(IRequest) error          //处理请求
}

package impl

import (
	"fmt"

	"github.com/ikun666/tcp_server/core/iface"
)

type router struct {
	routerMap map[uint32]iface.IHandler //对应路由处理handle
}

func newRouter() iface.IRouter {
	return &router{
		routerMap: make(map[uint32]iface.IHandler),
	}
}

// 添加路由
func (r *router) AddRoute(id uint32, handle iface.IHandler) error {
	if _, ok := r.routerMap[id]; ok {
		fmt.Println("handle has existed")
		return fmt.Errorf("handle has existed")
	}
	r.routerMap[id] = handle
	return nil
}

// 处理路由 从请求查找对应handle
func (r *router) DoRoute(req iface.IRequest) error {
	handle, ok := r.routerMap[req.GetMessage().GetMsgTag()]
	if !ok {
		fmt.Println("handle no exist")
		return fmt.Errorf("handle no exist")
	}
	handle.PreHandle(req)
	handle.Handle(req)
	handle.PostHandle(req)
	return nil
}

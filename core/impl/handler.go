package impl

import "github.com/ikun666/tcp_server/core/iface"

//BaseHandler实现IHandler接口，后续Handle内嵌BaseHandler
type BaseHandler struct {
}

//处理请求前hook函数
func (b *BaseHandler) PreHandle(req iface.IRequest) {

}

//处理请求
func (b *BaseHandler) Handle(req iface.IRequest) {

}

//处理请求后hook函数
func (b *BaseHandler) PostHandle(req iface.IRequest) {

}

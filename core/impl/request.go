package impl

import (
	common_iface "github.com/ikun666/tcp_server/common/iface"
	"github.com/ikun666/tcp_server/core/iface"
)

type request struct {
	conn iface.IConnection     //tcp连接
	msg  common_iface.IMessage //请求消息
}

// 获取连接
func (r *request) GetConnetion() iface.IConnection {
	return r.conn
}

// 请求消息
func (r *request) GetMessage() common_iface.IMessage {
	return r.msg
}
func newRequest(conn iface.IConnection, msg common_iface.IMessage) iface.IRequest {
	return &request{
		conn: conn,
		msg:  msg,
	}
}

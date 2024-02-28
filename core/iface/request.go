package iface

import common_iface "github.com/ikun666/tcp_server/common/iface"

type IRequest interface {
	GetConnetion() IConnection         //获取连接
	GetMessage() common_iface.IMessage //请求消息
}

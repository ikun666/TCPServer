package iface

import (
	"net"

	common_iface "github.com/ikun666/tcp_server/common/iface"
)

type IConnection interface {
	Start()                               //连接开始
	Stop()                                //连接结束
	GetConn() *net.TCPConn                //得到连接
	GetID() uint32                        //得到连接id
	Read() (common_iface.IMessage, error) //从连接读取数据
	Write(uint32, []byte) error           //写入数据到连接
	SetProperty(string, any)              //设置k-v
	GetProperty(string) (any, error)      //获取k-v
	RemoveProperty(string)                //删除k-v
	SetOnConnCreate(func(IConnection))    //设置连接创建hook
	SetOnConnDestroy(func(IConnection))   //设置连接销毁hook
}

package iface

type IConnManager interface {
	Add(IConnection)           //添加连接
	Remove(IConnection)        //删除连接
	Get(id uint32) IConnection //获取连接
	Len() int32                //得到连接数
	ClearConn()                //清空连接
}

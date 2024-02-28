package iface

type IServer interface {
	Start()                             //服务开始
	Stop()                              //服务结束
	Run()                               //服务运行
	AddRoute(uint32, IHandler)          //添加路由 tag对应handle
	GetConnManager() IConnManager       //得到连接管理器
	SetOnConnCreate(func(IConnection))  //设置连接创建hook
	SetOnConnDestroy(func(IConnection)) //设置连接销毁hook
}

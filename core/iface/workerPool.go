package iface

type IWorkerPool interface {
	StartWorkerPool(IRouter) //开启工作池
	Add2WorkerPool(IRequest) //添加工作到工作池
}

package impl

import (
	"fmt"
	"math/rand"

	"github.com/ikun666/tcp_server/core/iface"
)

type wokerPool struct {
	workerPoolData []chan iface.IRequest //工作池
	workerPoolSize uint32                //工作池大小 cpu线程数
	workerChanSize uint32                //每个线程绑定chan大小
}

func newWorkerPool(workerPoolSize, workerChanSize uint32) iface.IWorkerPool {
	return &wokerPool{
		workerPoolData: make([]chan iface.IRequest, workerChanSize),
		workerPoolSize: workerPoolSize,
		workerChanSize: workerChanSize,
	}
}

// 开启工作池
func (w *wokerPool) StartWorkerPool(router iface.IRouter) {
	fmt.Println("start worker pool")
	for i := 0; i < int(w.workerPoolSize); i++ {
		//初始化每一个chan
		w.workerPoolData[i] = make(chan iface.IRequest, w.workerChanSize)
		//开启一个go程
		go func(i int) {
			fmt.Println("worker ", i, " start")
			for {
				//只要有工作就处理
				req := <-w.workerPoolData[i]
				router.DoRoute(req)
			}
		}(i)
	}
}

// 添加工作到工作池
func (w *wokerPool) Add2WorkerPool(req iface.IRequest) {
	id := rand.Intn(int(w.workerPoolSize))
	// fmt.Println("add req ", id)
	w.workerPoolData[id] <- req
}

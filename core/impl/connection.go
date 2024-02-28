package impl

import (
	"fmt"
	"log/slog"
	"net"
	"sync"

	common_iface "github.com/ikun666/tcp_server/common/iface"
	common_impl "github.com/ikun666/tcp_server/common/impl"
	"github.com/ikun666/tcp_server/conf"
	"github.com/ikun666/tcp_server/core/iface"
	"github.com/ikun666/tcp_server/utils"
)

type connection struct {
	conn              *net.TCPConn       //tcp连接
	id                uint32             //连接id
	isClosed          bool               //连接已经关闭
	msgChan           chan []byte        //读写分离 reader读取后经过路由handle处理发送write到这里 writer从这里读取
	exitChan          chan struct{}      //通知writer退出
	workerPool        iface.IWorkerPool  //工作池
	connManager       iface.IConnManager //连接管理
	property          map[string]any     //存储k-v
	propertyLock      sync.RWMutex
	onConnCreateFunc  func(iface.IConnection) //连接创建hook
	onConnDestroyFunc func(iface.IConnection) //连接销毁hook
}

// 从连接读取数据
func (c *connection) Read() (common_iface.IMessage, error) {
	//读取消息
	//解包数据到消息
	msg, err := utils.UnPack(c.conn)
	return msg, err
}
func (c *connection) Reader() {
	for {
		//读取消息
		msg, err := c.Read()
		if err != nil {
			slog.Error("read", "err", err)
			return
		}
		// slog.Info("read msg", "msg", string(msg.GetMsgData()))
		//创建请求
		req := newRequest(c, msg)
		//添加到工作池
		c.workerPool.Add2WorkerPool(req)
	}
}

// 写入数据到msgChan
func (c *connection) Write(tag uint32, data []byte) error {
	if c.isClosed {
		return fmt.Errorf("conn close")
	}
	//封装消息
	msg := common_impl.NewMessage(tag, len(data), data)
	//打包消息
	sendMsg, err := utils.Pack(msg)
	if err != nil {
		return fmt.Errorf("pack msg err:%v", err)
	}
	//发送消息到msgChan
	c.msgChan <- sendMsg
	return nil
}
func (c *connection) Writer() {
	for {
		select {
		//如果有要发送的消息
		case msg := <-c.msgChan:
			_, err := c.conn.Write(msg)
			if err != nil {
				slog.Error("write msg", "err", err)
				return
			}
		//收到退出信号
		case <-c.exitChan:
			return
		}
	}
}

// 连接开始
func (c *connection) Start() {
	defer c.Stop()
	//添加到连接管理器
	c.connManager.Add(c)
	//执行hook
	c.OnConnCreate(c)
	go c.Writer()
	c.Reader()
}

// 连接结束
func (c *connection) Stop() {
	//关闭了不再重复关闭
	if c.isClosed {
		return
	}
	c.isClosed = true
	slog.Info("stop", "id", c.id)
	c.OnConnDestroy(c)

	//通知写也退出
	c.exitChan <- struct{}{}
	//释放资源
	c.connManager.Remove(c)

	c.conn.Close()
}

// 得到连接
func (c *connection) GetConn() *net.TCPConn {
	return c.conn
}

// 得到连接id
func (c *connection) GetID() uint32 {
	return c.id
}
func (c *connection) SetProperty(key string, value any) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}
func (c *connection) GetProperty(key string) (any, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if v, ok := c.property[key]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("no this %s property", key)
}
func (c *connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
func (c *connection) SetOnConnCreate(hook func(iface.IConnection)) {
	c.onConnCreateFunc = hook
}
func (c *connection) SetOnConnDestroy(hook func(iface.IConnection)) {
	c.onConnDestroyFunc = hook
}
func (c *connection) OnConnCreate(conn iface.IConnection) {
	if c.onConnCreateFunc != nil {
		c.onConnCreateFunc(conn)
	}
}
func (c *connection) OnConnDestroy(conn iface.IConnection) {
	if c.onConnDestroyFunc != nil {
		c.onConnDestroyFunc(conn)
	}
}
func newConnection(conn *net.TCPConn, id uint32, workerPool iface.IWorkerPool, connManager iface.IConnManager) iface.IConnection {
	return &connection{
		conn:        conn,
		id:          id,
		isClosed:    false,
		msgChan:     make(chan []byte, conf.GConfig.MaxMsgChanSize),
		exitChan:    make(chan struct{}, 1),
		workerPool:  workerPool,
		connManager: connManager,
		property:    make(map[string]any),
	}
}

package impl

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	common_impl "github.com/ikun666/tcp_server/common/impl"
	"github.com/ikun666/tcp_server/conf"
	"github.com/ikun666/tcp_server/core/iface"
	"github.com/ikun666/tcp_server/utils"
)

type Server struct {
	TCPVersion        string                  //tcp版本
	Name              string                  //服务器名称
	IP                string                  //IP地址
	Port              uint16                  //端口
	Router            iface.IRouter           //路由
	WokerPool         iface.IWorkerPool       //工作池
	ConnManager       iface.IConnManager      //连接管理器
	ExitChan          chan struct{}           //退出chan
	OnConnCreateFunc  func(iface.IConnection) //连接创建hook
	OnConnDestroyFunc func(iface.IConnection) //连接销毁hook
}

func (s *Server) Start() {
	//0.开启工作池
	s.WokerPool.StartWorkerPool(s.Router)
	//1.解析tcp ip port
	addr, err := net.ResolveTCPAddr(s.TCPVersion, fmt.Sprintf("%v:%v", s.IP, s.Port))
	if err != nil {
		slog.Error("net.ResolveTCPAddr err", "err", err)
	}
	//2.监听tcp
	listener, err := net.ListenTCP(s.TCPVersion, addr)
	if err != nil {
		slog.Error("net.ListenTCP err", "err", err)
	}
	go func() {
		//连接id
		var id uint32 = 0
		//3 循环处理
		for {
			conn, err := listener.AcceptTCP()
			//如果服务器已经关闭
			if errors.Is(err, net.ErrClosed) {
				slog.Error("Listener closed", "err", err)
				return
			}
			//连接有问题跳过本连接
			if err != nil {
				slog.Error("listener.AcceptTCP err", "err", err)
				continue
			}
			// fmt.Println("conn user:", s.ConnManager.Len())
			//超过最大连接拒绝处理
			if s.ConnManager.Len()+1 > conf.GConfig.MaxConn {
				go deferCloseConn(conn)
				continue
			}
			dealconn := newConnection(conn, id, s.WokerPool, s.ConnManager)
			dealconn.SetOnConnCreate(s.OnConnCreateFunc)
			dealconn.SetOnConnDestroy(s.OnConnDestroyFunc)
			id++
			go dealconn.Start()
		}
	}()

	<-s.ExitChan
	if err := listener.Close(); err != nil {
		slog.Error("Listener close", "err", err)
		return
	}
}
func (s *Server) Stop() {
	s.ConnManager.ClearConn()
}
func (s *Server) Run() {
	defer s.Stop()
	go s.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	s.ExitChan <- struct{}{}
	time.Sleep(5 * time.Second)
	fmt.Println("Server stop")
}
func (s *Server) AddRoute(tag uint32, handle iface.IHandler) {
	s.Router.AddRoute(tag, handle)
}
func NewServer() iface.IServer {
	err := conf.Init("conf/conf.json")
	if err != nil {
		slog.Error("conf.Init", "err", err)
		return &Server{
			TCPVersion:  "tcp",
			Name:        "ikun666",
			IP:          "127.0.0.1",
			Port:        8080,
			Router:      newRouter(),
			WokerPool:   newWorkerPool(8, 1024),
			ConnManager: newConnManager(),
			ExitChan:    make(chan struct{}),
		}
	} else {
		return &Server{
			TCPVersion:  conf.GConfig.TCPVersion,
			Name:        conf.GConfig.ServerName,
			IP:          conf.GConfig.IP,
			Port:        conf.GConfig.Port,
			Router:      newRouter(),
			WokerPool:   newWorkerPool(conf.GConfig.WorkerPoolSize, conf.GConfig.WorkerChanSize),
			ConnManager: newConnManager(),
			ExitChan:    make(chan struct{}),
		}
	}
}
func (s *Server) SetOnConnCreate(hook func(iface.IConnection)) {
	s.OnConnCreateFunc = hook
}
func (s *Server) SetOnConnDestroy(hook func(iface.IConnection)) {
	s.OnConnDestroyFunc = hook
}

// 得到连接管理器
func (s *Server) GetConnManager() iface.IConnManager {
	return s.ConnManager
}
func deferCloseConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("=======users too much,try to connect later======")
	data := []byte("users too much,try to connect later")
	msg := common_impl.NewMessage(401, len(data), data)
	sendMsg, err := utils.Pack(msg)
	if err != nil {
		return
	}
	conn.Write(sendMsg)
	time.Sleep(5 * time.Second)
}

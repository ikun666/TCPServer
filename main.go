package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ikun666/tcp_server/core/iface"
	"github.com/ikun666/tcp_server/core/impl"
)

type echo struct {
	impl.BaseHandler
}

func (e *echo) PreHandle(req iface.IRequest) {
	fmt.Println("PreHandle")
}
func (e *echo) Handle(req iface.IRequest) {
	fmt.Println(string(req.GetMessage().GetMsgData()))
	err := req.GetConnetion().Write(1, req.GetMessage().GetMsgData())
	if err != nil {
		slog.Error("echo Handle", "err", err)
		return
	}
}

type hello struct {
	impl.BaseHandler
}

func (h *hello) Handle(req iface.IRequest) {
	fmt.Println(string(req.GetMessage().GetMsgData()))
	err := req.GetConnetion().Write(1, []byte(fmt.Sprintf("hello: %v", string(req.GetMessage().GetMsgData()))))
	if err != nil {
		slog.Error("echo Handle", "err", err)
		return
	}
}

type stop struct {
	impl.BaseHandler
}

func (s *stop) Handle(req iface.IRequest) {
	fmt.Println(string(req.GetMessage().GetMsgData()))
	fmt.Println("stop")
	req.GetConnetion().Stop()
}
func OnConnCreate(conn iface.IConnection) {
	fmt.Println("----------OnConnCreate----------")
	conn.SetProperty("GitHub", "github.com/ikun666")
	conn.Write(123, []byte("----------OnConnCreate----------"))
}
func OnConnDestroy(conn iface.IConnection) {
	fmt.Println("----------OnConnDestroy----------")
	if v, err := conn.GetProperty("GitHub"); err == nil {
		fmt.Println(v)
	}
	conn.Write(456, []byte("----------OnConnDestroy----------"))
	time.Sleep(2 * time.Second)
}
func main() {
	// utils.InitLogger()

	server := impl.NewServer()
	server.SetOnConnCreate(OnConnCreate)
	server.SetOnConnDestroy(OnConnDestroy)
	server.AddRoute(1, &echo{})
	server.AddRoute(2, &hello{})
	server.AddRoute(3, &stop{})
	server.Run()

}

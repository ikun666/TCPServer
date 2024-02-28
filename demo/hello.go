package main

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/ikun666/tcp_server/common/impl"
	"github.com/ikun666/tcp_server/utils"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		slog.Error("Dial", "err", err)
		return
	}
	defer conn.Close()
	go func() {

		for {
			//读取消息
			//解包数据到消息
			msg, err := utils.UnPack(conn)
			if err != nil {
				slog.Error("Read", "err", err)
				return
			}
			fmt.Printf("recieve msg:%v\n", string(msg.GetMsgData()))
		}
	}()
	for i := 0; i < 5; i++ {
		//封装消息
		data := "ikun"
		msg := impl.NewMessage(2, len(data), []byte(data))
		//打包消息
		sendMsg, err := utils.Pack(msg)
		if err != nil {
			slog.Error("utils.Pack", "err", err)
			return
		}
		//发送消息
		_, err = conn.Write(sendMsg)
		if err != nil {
			slog.Error("Write", "err", err)
			return
		}
		time.Sleep(2 * time.Second)
	}
	//封装消息
	data := "ikun"
	msg := impl.NewMessage(3, len(data), []byte(data))
	//打包消息
	sendMsg, err := utils.Pack(msg)
	if err != nil {
		slog.Error("utils.Pack", "err", err)
		return
	}
	//发送消息
	_, err = conn.Write(sendMsg)
	if err != nil {
		slog.Error("Write", "err", err)
		return
	}
	time.Sleep(2 * time.Second)
}

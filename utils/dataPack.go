package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ikun666/tcp_server/common/iface"
	"github.com/ikun666/tcp_server/common/impl"
)

func Pack(message iface.IMessage) ([]byte, error) {
	var buf bytes.Buffer
	//TLV 格式 依次写入 类型 长度 内容
	err := binary.Write(&buf, binary.LittleEndian, message.GetMsgTag())
	if err != nil {
		// slog.Error("write tag", "err", err)
		return nil, err
	}
	err = binary.Write(&buf, binary.LittleEndian, message.GetMsgLen())
	if err != nil {
		// slog.Error("write len", "err", err)
		return nil, err
	}
	err = binary.Write(&buf, binary.LittleEndian, message.GetMsgData())
	if err != nil {
		// slog.Error("write data", "err", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func UnPack(conn io.Reader) (iface.IMessage, error) {
	//先读消息头 TL 共计8字节
	//head[0] tag head[1] len
	head := make([]uint32, 2)
	err := binary.Read(conn, binary.LittleEndian, head)
	if err != nil {
		// slog.Error("read head", "err", err)
		return nil, err
	}
	//过大丢弃
	if head[1] > 2048 {
		// slog.Info("read len over MaxPackageSize", "len", head[1])
		return nil, fmt.Errorf("read len over MaxPackageSize")
	}
	//根据长度读消息体
	body := make([]byte, head[1])
	err = binary.Read(conn, binary.LittleEndian, body)
	if err != nil {
		// slog.Error("read body", "err", err)
		return nil, err
	}
	msg := impl.NewMessage(head[0], int(head[1]), body)
	return msg, nil
}

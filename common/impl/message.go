package impl

import "github.com/ikun666/tcp_server/common/iface"

type Message struct {
	tag  uint32 //消息类型
	len  uint32 //消息长度
	data []byte //消息内容
}

func NewMessage(tag uint32, len int, data []byte) iface.IMessage {
	return &Message{
		tag:  tag,
		len:  uint32(len),
		data: data,
	}
}
func (m *Message) GetMsgLen() uint32 {
	return m.len
}
func (m *Message) GetMsgTag() uint32 {
	return m.tag
}
func (m *Message) GetMsgData() []byte {
	return m.data
}

func (m *Message) SetMsgLen(len uint32) {
	m.len = len
}
func (m *Message) SetMsgTag(tag uint32) {
	m.tag = tag
}
func (m *Message) SetMsgData(data []byte) {
	m.data = data
}

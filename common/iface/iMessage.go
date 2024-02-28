package iface

type IMessage interface {
	GetMsgLen() uint32  //消息长度
	GetMsgTag() uint32  //消息类型
	GetMsgData() []byte //消息内容

	SetMsgLen(uint32)  //设置消息长度
	SetMsgTag(uint32)  //设置消息类型
	SetMsgData([]byte) //设置消息内容
}

package conf

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

// 全局配置
type GlobalConfig struct {
	IP         string `json:"ip,omitempty"`
	Port       uint16 `json:"port,omitempty"`
	ServerName string `json:"server_name,omitempty"`

	TCPVersion     string `json:"tcp_version,omitempty"`
	MaxConn        int32  `json:"max_conn,omitempty"`
	MaxPackageSize uint32 `json:"max_package_size,omitempty"`
	MaxMsgChanSize uint32 `json:"max_msg_chan_size,omitempty"`
	WorkerPoolSize uint32 `json:"worker_pool_size,omitempty"`
	WorkerChanSize uint32 `json:"worker_chan_size,omitempty"`
}

// 全局配置变量
var GConfig *GlobalConfig

func Init(path string) error {
	GConfig = &GlobalConfig{}
	data, err := os.ReadFile(path)
	if err != nil {
		slog.Error("ReadFile", "err", err)
		return fmt.Errorf("read conf err:%v", err)
	}
	err = json.Unmarshal(data, GConfig)
	if err != nil {
		slog.Error("Unmarshal", "err", err)
		return fmt.Errorf("unmarshal conf err:%v", err)
	}
	return nil
}

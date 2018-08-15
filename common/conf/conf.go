package conf

import (
	"digicon/common/xtime"
)

type CommConf struct {
	Ver     string
	LogPath string
}

// =================================== HTTP ==================================
// HTTPServer http server settings.
type HTTPServer struct {
	Addrs        []string
	MaxListen    int32
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
}

// HTTPClient http client settings.
type HTTPClient struct {
	Dial      xtime.Duration
	Timeout   xtime.Duration
	KeepAlive xtime.Duration
	Timer     int
}

// MultiHttp outer/inner/local http server settings.
type MultiHTTP struct {
	Outer *HTTPServer
	Inner *HTTPServer
	Local *HTTPServer
}

type Server struct {
	Proto string
	Addr  string
}

type RPCServer struct {
	Proto string
	Addr  string
}

type ConfDiscovery struct {
	Role     string
	Interval xtime.Duration
}

type ServiceDiscoveryServer struct {
	ServiceName string `json:"service_name"`
	RPCAddr     string `json:"rpc_addr"`
	ConsulAddr    string `json:"consul_addr"`
	Interval    xtime.Duration `json:"interval"`
	TTL         xtime.Duration `json:"ttl"`
 }

type ServiceDiscoveryClient struct {
	ServiceName string
	EtcdAddr    string
	Balancer    string
}

type Etcd struct {
	Name    string
	Root    string
	Addrs   []string
	Timeout xtime.Duration
}

type Zookeeper struct {
	Root    string
	Addrs   []string
	Timeout xtime.Duration
}

// Redis client settings.
type Redis struct {
	Name         string  `json:"name"`
	Proto        string `json:"proto"`
	Addr         string `json:"addr"`
	Active       int  `json:"active"`
	Idle         int  `json:"idle"`
	DialTimeout  xtime.Duration `json:"dial_timeout"`
	ReadTimeout  xtime.Duration `json:"read_timeout"`
	WriteTimeout xtime.Duration `json:"write_timeout"`
	IdleTimeout  xtime.Duration `json:"idle_timeout"`
}

// KafkaProducer kafka producer settings.
type KafkaProducer struct {
	Zookeeper *Zookeeper
	Brokers   []string
	Sync      bool // true: sync, false: async
}

// KafkaConsumer kafka client settings.
type KafkaConsumer struct {
	Group     string
	Topics    []string
	Offset    bool // true: new, false: old
	Zookeeper *Zookeeper
}

type MySQL struct {
	Name   string `json:"name"` // for trace
	DSN    string `json:"dsn"` // data source name
	Active int    `json:"active"` // pool
	Idle   int    `json:"idle"` // pool
}

type MongoDB struct {
	Addrs       string
	DB          string
	DialTimeout xtime.Duration
}

type ES struct {
	Addrs string
}

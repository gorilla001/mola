package main

import (
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
)

const (
	keepaliveTime    = 30 * time.Second
	keepaliveTimeout = 10 * time.Second
	dialTimeout      = 20 * time.Second
)

func newETCD3Client(c etcdConfig) (*clientv3.Client, error) {
	tlsInfo := transport.TLSInfo{
		CertFile: c.CertFile,
		KeyFile:  c.KeyFile,
		CAFile:   c.CAFile,
	}
	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		return nil, err
	}
	// NOTE: Client relies on nil tlsConfig
	// for non-secure connections, update the implicit variable
	if len(c.CertFile) == 0 && len(c.KeyFile) == 0 && len(c.CAFile) == 0 {
		tlsConfig = nil
	}

	cfg := clientv3.Config{
		DialTimeout:          dialTimeout,
		DialKeepAliveTime:    keepaliveTime,
		DialKeepAliveTimeout: keepaliveTimeout,
		Endpoints:            c.ServerList,
		TLS:                  tlsConfig,
	}

	client, err := clientv3.New(cfg)

	return client, err
}

package zookeeper

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-zookeeper/zk"
)

var (
	zooClient     *ZookeeperClient
	zooClientOnce sync.Once
)

type ZookeeperClient struct {
	conn *zk.Conn
	mu   sync.Mutex
}

func Connect(hosts []string) {
	zooClientOnce.Do(func() {
		conn, _, err := zk.Connect(hosts, time.Second)
		if err != nil {
			panic(fmt.Sprintf("Error connecting to zookeeper: %s", err.Error()))
		}

		zooClient = &ZookeeperClient{
			conn: conn,
		}
	})
}

func Register(path string, data []byte) error {
	_, err := zooClient.conn.Create(path, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}

	return nil
}

func Discover(basePath string) ([]string, error) {
	children, _, err := zooClient.conn.Children(basePath)
	if err != nil {
		return nil, err
	}

	return children, nil
}

func Close() {
	if zooClient != nil {
		zooClient.mu.Lock()
		defer zooClient.mu.Unlock()
		if zooClient.conn != nil {
			zooClient.conn.Close()
			zooClient.conn = nil // Prevent double close
		}
	}
}

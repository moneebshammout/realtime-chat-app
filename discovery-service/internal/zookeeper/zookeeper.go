package zookeeper

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"discovery-service/pkg/utils"

	"github.com/go-zookeeper/zk"
)

var (
	zooClient     *ZookeeperClient
	zooClientOnce sync.Once
	logger        = utils.GetLogger()
)

type ZookeeperClient struct {
	conn *zk.Conn
	mu   sync.Mutex
}

func Connect(hosts []string) {
	zooClientOnce.Do(func() {
		conn, _, err := zk.Connect(hosts, time.Second)
		if err != nil {
			logger.Panicf("ZookeeperClient: Error connecting to zookeeper: %s", err.Error())
		}

		zooClient = &ZookeeperClient{
			conn: conn,
		}
	})
}

func Register(path string, node string) error {
	// if path exist get urls and append to it
	exists, _, err := zooClient.conn.Exists(path)
	if err != nil {
		logger.Errorf("ZookeeperClient: Error checking if path %s exists: %s", path, err.Error())
		return err
	}

	if exists {
		nodes, err := Discover(path)
		if err != nil {
			logger.Errorf("ZookeeperClient: Error discovering path %s: %s", path, err.Error())
			return err
		}

		nodes = append(nodes, node)
		node = strings.Join(nodes, ",")
		newNodes, err := json.Marshal(node)
		if err != nil {
			logger.Errorf("ZookeeperClient: Error marshalling nodes: %s", err.Error())
			return err
		}

		_, err = zooClient.conn.Set(path, newNodes, -1)
		if err != nil {
			logger.Errorf("ZookeeperClient: Error updating path %s: %s", path, err.Error())
			return err
		}

	} else {
		data, err := json.Marshal(node)
		if err != nil {
			logger.Errorf("ZookeeperClient: Error marshalling node: %s", err.Error())
			return err
		}

		_, err = zooClient.conn.Create(path, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
		if err != nil {
			logger.Errorf("ZookeeperClient: Error creating path %s: %s", path, err.Error())
			return err
		}
	}

	return nil
}

func Discover(basePath string) ([]string, error) {
	children, _, err := zooClient.conn.Get(basePath)
	if err != nil {
		logger.Errorf("ZookeeperClient: Error discovering path %s: %s", basePath, err.Error())
		return nil, err
	}

	var data string
	err = json.Unmarshal(children, &data)
	if err != nil {
		logger.Errorf("ZookeeperClient: Error unmarshalling data: %s", err.Error())
		return nil, err
	}

	return strings.Split(data, ","), nil
}

func Close() {
	if zooClient != nil {
		zooClient.mu.Lock()
		defer zooClient.mu.Unlock()
		if zooClient.conn != nil {
			logger.Info("Closing zookeeper connection")
			zooClient.conn.Close()
			zooClient.conn = nil // Prevent double close
		}
	}
}

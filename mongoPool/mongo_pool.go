package mongoPool

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jericho-yu/aid/dict"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	MongoClientPool struct {
		clients *dict.AnyDict[string, *MongoClient]
	}

	MongoClient struct {
		url    string
		client *mongo.Client
	}
)

var (
	mongoClientPool *MongoClientPool
	mongoPoolOnce   sync.Once
)

// NewMongoClient 实例化：mongo客户端
func NewMongoClient(url string) (*MongoClient, error) {
	var err error
	mc := &MongoClient{url: url}
	// 设置客户端选项
	clientOptions := options.Client().ApplyURI(mc.url)

	// 连接到 MongoDB
	if mc.client, err = mongo.Connect(context.TODO(), clientOptions); err != nil {
		return nil, err
	}

	// 检查连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = mc.client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return mc, nil
}

// Close 关闭客户端
func (my *MongoClient) Close() error {
	return my.client.Disconnect(context.Background())
}

// GetClient 获取客户端链接
func (my *MongoClient) GetClient() *mongo.Client {
	return my.client
}

// OnceMongoPool 单例化：mongodb连接池
func OnceMongoPool() *MongoClientPool {
	mongoPoolOnce.Do(func() {
		mongoClientPool = &MongoClientPool{clients: dict.MakeAnyDict[string, *MongoClient]()}
	})
	return mongoClientPool
}

// Append 增加客户端
func (*MongoClientPool) Append(key string, mongoClient *MongoClient) (*MongoClientPool, error) {
	if mongoClientPool.clients.Has(key) {
		return mongoClientPool, errors.New("客户端已存在")
	}

	mongoClientPool.clients.Set(key, mongoClient)

	return mongoClientPool, nil
}

// 清除客户端
func (*MongoClientPool) Remove(key string) (*MongoClientPool, error) {
	if !mongoClientPool.clients.Has(key) {
		return mongoClientPool, errors.New("客户端不存在")
	}

	if mongoClient, exist := mongoClientPool.clients.Get(key); exist {
		return mongoClientPool, errors.New("客户端不存在")
	} else {
		if err := mongoClient.Close(); err != nil {
			return mongoClientPool, err
		}

		mongoClientPool.clients.RemoveByKey(key)
	}

	return mongoClientPool, nil
}

// Clean 清理客户端
func (*MongoClientPool) Clean() *MongoClientPool {
	for _, key := range mongoClientPool.clients.Keys() {
		mongoClientPool.Remove(key)
	}
	return mongoClientPool
}

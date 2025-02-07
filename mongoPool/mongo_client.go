package mongoPool

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	MongoClient struct {
		url               string
		client            *mongo.Client
		currentDatabase   *mongo.Database
		currentCollection *mongo.Collection
		condition         Map
	}

	Data  = []any
	Datum = primitive.D
	KV    = primitive.E
	Map   = primitive.M
)

// NewMongoClient 实例化：mongo客户端
func NewMongoClient(url string) (*MongoClient, error) {
	var err error
	mc := &MongoClient{url: url, condition: Map{}}
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

// Ping 测试链接
func (my *MongoClient) Ping() error {
	var err error
	// 检查连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = my.client.Ping(ctx, nil); err != nil {
		return err
	}
	return nil
}

// SetDatabase 设置数据库
func (my *MongoClient) SetDatabase(database string, opts ...*options.DatabaseOptions) *MongoClient {
	my.currentDatabase = my.client.Database(database, opts...)
	return my
}

// SetCollection 设置文档
func (my *MongoClient) SetCollection(collection string, opts ...*options.CollectionOptions) *MongoClient {
	my.currentCollection = my.currentDatabase.Collection(collection, opts...)
	return my
}

// InsertOne 插入一条数据
func (my *MongoClient) InsertOne(data Datum) (*mongo.InsertOneResult, error) {
	return my.currentCollection.InsertOne(context.TODO(), data)
}

// InsertMany 插入多条数据
func (my *MongoClient) InsertMany(data []any) (*mongo.InsertManyResult, error) {
	return my.currentCollection.InsertMany(context.TODO(), data)
}

// Where 设置查询条件
func (my *MongoClient) Where(condition Map) *MongoClient {
	my.condition = condition
	return my
}

// FindOne 查询一条数据
func (my *MongoClient) FindOne(result *Map, findOneOptionFn func(opt *options.FindOneOptions) *options.FindOneOptions) error {
	defer func() { my.condition = Map{} }()
	var findOneOption *options.FindOneOptions
	if findOneOptionFn != nil {
		findOneOption = findOneOptionFn(options.FindOne())
	}
	return my.currentCollection.FindOne(context.TODO(), my.condition, findOneOption).Decode(&result)
}

// FindMany 查询多条数据
func (my *MongoClient) FindMany(results *[]Map, findOptionFn func(opt *options.FindOptions) *options.FindOptions) error {
	var (
		err        error
		findOption *options.FindOptions
		cursor     *mongo.Cursor
	)
	defer func() { my.condition = Map{} }()
	if findOptionFn != nil {
		findOption = findOptionFn(options.Find())
	}
	if cursor, err = my.currentCollection.Find(context.TODO(), my.condition, findOption); err != nil {
		return err
	} else {
		return cursor.All(context.TODO(), results)
	}
}

// DeleteOne 删除单条数据
func (my *MongoClient) DeleteOne() (*mongo.DeleteResult, error) {
	defer func() { my.condition = Map{} }()
	return my.currentCollection.DeleteOne(context.TODO(), my.condition)
}

// DeleteMany 删除多条数据
func (my *MongoClient) DeleteMany() (*mongo.DeleteResult, error) {
	defer func() { my.condition = Map{} }()
	return my.currentCollection.DeleteMany(context.TODO(), my.condition)
}

// NewMap 新建Map数据
func NewMap(Key string, Value any) Map {
	return Map{Key: Value}
}

// NewKV 新建KeyValue字段数据
func NewKV(Key string, Value any) KV {
	return KV{Key: Key, Value: Value}
}

// NewData 新建Datum单条数据
func NewDatum(kv ...KV) Datum {
	var d = make(Datum, len(kv))
	copy(d, kv)
	return d
}

// NewData 新建Data多条数据集
func NewData(datum ...Datum) Data {
	return Data{datum}
}

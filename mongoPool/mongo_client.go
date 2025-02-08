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
		Err               error
	}

	Data   = primitive.D
	Entity = primitive.E
	Map    = primitive.M
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
func (my *MongoClient) Ping() {
	// 检查连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	my.Err = my.client.Ping(ctx, nil)
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
func (my *MongoClient) InsertOne(data Data) *mongo.InsertOneResult {
	var insertOneRes *mongo.InsertOneResult
	insertOneRes, my.Err = my.currentCollection.InsertOne(context.TODO(), data)
	return insertOneRes
}

// InsertMany 插入多条数据
func (my *MongoClient) InsertMany(data []any) *mongo.InsertManyResult {
	var res *mongo.InsertManyResult
	res, my.Err = my.currentCollection.InsertMany(context.TODO(), data)
	return res
}

// Where 设置查询条件
func (my *MongoClient) Where(condition Map) *MongoClient {
	my.condition = condition
	return my
}

// CleanCondition 清理查询条件
func (my *MongoClient) CleanCondition() {
	my.condition = Map{}
}

// FindOne 查询一条数据
func (my *MongoClient) FindOne(result *Map, findOneOptionFn func(opt *options.FindOneOptions) *options.FindOneOptions) *MongoClient {
	defer my.CleanCondition()
	var findOneOption *options.FindOneOptions
	if findOneOptionFn != nil {
		findOneOption = findOneOptionFn(options.FindOne())
	}
	my.Err = my.currentCollection.FindOne(context.TODO(), my.condition, findOneOption).Decode(&result)
	return my
}

// FindMany 查询多条数据
func (my *MongoClient) FindMany(results *[]Map, findOptionFn func(opt *options.FindOptions) *options.FindOptions) *MongoClient {
	var (
		err        error
		findOption *options.FindOptions
		cursor     *mongo.Cursor
	)
	defer my.CleanCondition()
	if findOptionFn != nil {
		findOption = findOptionFn(options.Find())
	}
	if cursor, my.Err = my.currentCollection.Find(context.TODO(), my.condition, findOption); err != nil {
		return my
	} else {
		my.Err = cursor.All(context.TODO(), results)
		return my
	}
}

// DeleteOne 删除单条数据
func (my *MongoClient) DeleteOne() *mongo.DeleteResult {
	var res *mongo.DeleteResult
	defer my.CleanCondition()
	if res, my.Err = my.currentCollection.DeleteOne(context.TODO(), my.condition); my.Err != nil {
		return nil
	}
	return res
}

// DeleteMany 删除多条数据
func (my *MongoClient) DeleteMany() *mongo.DeleteResult {
	var res *mongo.DeleteResult
	defer my.CleanCondition()
	if res, my.Err = my.currentCollection.DeleteMany(context.TODO(), my.condition); my.Err != nil {
		return nil
	}
	return res
}

// NewMap 新建Map数据
func NewMap(Key string, Value any) Map {
	return Map{Key: Value}
}

// NewEntity 新建实体数据
func NewEntity(Key string, Value any) Entity {
	return Entity{Key: Key, Value: Value}
}

// NewData 新建单条数据
func NewData(kv ...Entity) Data {
	var d = make(Data, len(kv))
	copy(d, kv)
	return d
}

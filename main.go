package main

import (
	"log"

	"github.com/jericho-yu/aid/mongoPool"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	var (
		err           error
		insertOneRes  *mongo.InsertOneResult
		insertManyRes *mongo.InsertManyResult
		findOneRes    mongoPool.Map
		findManyRes   []mongoPool.Map
		deleteOneRes  *mongo.DeleteResult
		deleteManyRes *mongo.DeleteResult
	)
	mp := mongoPool.OnceMongoPool()
	mc, err := mongoPool.NewMongoClient("mongodb://admin:admin@localhost:27017")
	if err != nil {
		log.Fatalf("创建mongo客户端失败：%v", err)
	}
	mp.Append("default", mc)
	mc = mp.GetClient("default").SetDatabase("test_db").SetCollection("test_collection")

	// 清空数据
	if _, err = mc.DeleteMany(); err != nil {
		log.Fatalf("清空数据失败：%v", err)
	}

	// 插入单条数据
	insertOneRes, err = mc.InsertOne(mongoPool.NewDatum(mongoPool.NewKV("name", "张三"), mongoPool.NewKV("age", 18)))
	if err != nil {
		log.Fatalf("插入单条数据失败：%v", err)
	}
	log.Printf("插入单条数据成功：%v\n", insertOneRes.InsertedID)

	// 插入多条数据
	insertManyRes, err = mc.InsertMany(mongoPool.NewData(
		mongoPool.NewDatum(mongoPool.NewKV("name", "李四"), mongoPool.NewKV("age", 19)),
		mongoPool.NewDatum(mongoPool.NewKV("name", "王五"), mongoPool.NewKV("age", 20)),
		mongoPool.NewDatum(mongoPool.NewKV("name", "赵六"), mongoPool.NewKV("age", 21)),
	))
	if err != nil {
		log.Fatalf("插入多条数据失败：%v", err)
	}
	log.Printf("插入多条数据成功：%v\n", insertManyRes.InsertedIDs)

	// 查询单条数据
	if err = mc.Where(mongoPool.NewMap("_id", insertOneRes.InsertedID)).FindOne(&findOneRes, nil); err != nil {
		log.Fatalf("查询单条数据失败：%v", err)
	}
	log.Printf("查询单条数据成功：%v", findOneRes)

	// 查询多条数据
	if err = mc.Where(mongoPool.NewMap("_id", mongoPool.NewMap("$in", insertManyRes.InsertedIDs[1:]))).FindMany(&findManyRes, nil); err != nil {
		log.Fatalf("查询多条数据失败：%v", err)
	}
	log.Printf("查询多条数据成功：%v\n", findManyRes)

	// 删除单条数据
	if deleteOneRes, err = mc.Where(mongoPool.NewMap("_id", insertOneRes.InsertedID)).DeleteOne(); err != nil {
		log.Fatalf("删除单条数据失败：%v", err)
	}
	log.Printf("删除单条数据成功：%d", deleteOneRes.DeletedCount)

	// 删除多条数据
	if deleteManyRes, err = mc.Where(mongoPool.NewMap("_id", mongoPool.NewMap("$in", insertManyRes.InsertedIDs[1:]))).DeleteMany(); err != nil {
		log.Fatalf("删除多条数据失败：%v", err)
	}
	log.Printf("删除多条数据成功：%d", deleteManyRes.DeletedCount)

	// 查询剩余数据
	if err = mc.FindMany(&findManyRes, nil); err != nil {
		log.Fatalf("查询剩余数据失败：%v", err)
	}
	log.Printf("查询剩余数据：%v\n", findManyRes)
}

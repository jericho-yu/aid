package mongoPool

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

func Demo() {
	var (
		err           error
		insertOneRes  *mongo.InsertOneResult
		insertManyRes *mongo.InsertManyResult
		findOneRes    Map
		findManyRes   []Map
		deleteOneRes  *mongo.DeleteResult
		deleteManyRes *mongo.DeleteResult
	)
	mp := OnceMongoPool()
	mc, err := NewMongoClient("mongodb://admin:admin@localhost:27017")
	if err != nil {
		log.Fatalf("创建mongo客户端失败：%v", err)
	}
	mp.AppendClient("default", mc)
	mc = mp.GetClient("default").SetDatabase("test_db").SetCollection("test_collection")

	// 清空数据
	_ = mc.DeleteMany()
	if mc.Err != nil {
		log.Fatalf("清空数据失败：%v", err)
	}

	// 插入单条数据
	insertOneRes = mc.InsertOne(NewData(NewEntity("name", "张三"), NewEntity("age", 18)))
	if mc.Err != nil {
		log.Fatalf("插入单条数据失败：%v", err)
	}
	log.Printf("插入单条数据成功：%v\n", insertOneRes.InsertedID)

	// 插入多条数据
	insertManyRes = mc.InsertMany([]any{
		NewData(NewEntity("name", "李四"), NewEntity("age", 19)),
		NewData(NewEntity("name", "王五"), NewEntity("age", 20)),
		NewData(NewEntity("name", "赵六"), NewEntity("age", 21)),
	})
	if mc.Err != nil {
		log.Fatalf("插入多条数据失败：%v", err)
	}
	log.Printf("插入多条数据成功：%v\n", insertManyRes.InsertedIDs)

	// 查询单条数据
	if mc.Where(NewMap("_id", insertOneRes.InsertedID)).FindOne(&findOneRes, nil).Err != nil {
		log.Fatalf("查询单条数据失败：%v", mc.Err)
	}
	log.Printf("查询单条数据成功：%v", findOneRes)

	// 查询多条数据
	if mc.Where(NewMap("_id", NewMap("$in", insertManyRes.InsertedIDs[1:]))).FindMany(&findManyRes, nil).Err != nil {
		log.Fatalf("查询多条数据失败：%v", mc.Err)
	}
	log.Printf("查询多条数据成功：%v\n", findManyRes)

	// 删除单条数据
	deleteOneRes = mc.Where(NewMap("_id", insertOneRes.InsertedID)).DeleteOne()
	if mc.Err != nil {
		log.Fatalf("删除单条数据失败：%v", mc.Err)
	}
	log.Printf("删除单条数据成功：%d", deleteOneRes.DeletedCount)

	// 删除多条数据
	deleteManyRes = mc.Where(NewMap("_id", NewMap("$in", insertManyRes.InsertedIDs[1:]))).DeleteMany()
	if mc.Err != nil {
		log.Fatalf("删除多条数据失败：%v", mc.Err)
	}
	log.Printf("删除多条数据成功：%d", deleteManyRes.DeletedCount)

	// 查询剩余数据
	if mc.FindMany(&findManyRes, nil).Err != nil {
		log.Fatalf("查询剩余数据失败：%v", err)
	}
	log.Printf("查询剩余数据：%v\n", findManyRes)
}

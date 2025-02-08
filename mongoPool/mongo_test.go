package mongoPool

import (
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func getDB(t *testing.T) (*MongoClientPool, *MongoClient) {
	var err error
	mp := OnceMongoPool()
	mc, err := NewMongoClient("mongodb://admin:admin@localhost:27017")
	if err != nil {
		t.Fatalf("创建mongo客户端失败：%v", err)
	}
	if _, err = mp.Append("default", mc); err != nil {
		t.Fatalf("添加mongo客户端失败：%v", err)
	}
	mc = mp.GetClient("default").SetDatabase("test_db").SetCollection("test_collection")

	// 清空数据
	_ = mc.DeleteMany()
	if mc.Err != nil {
		t.Fatalf("清空数据失败：%v", err)
	}

	return mp, mc
}

func Test1One(t *testing.T) {
	t.Run("操作单条数据", func(t *testing.T) {
		var (
			err          error
			insertOneRes *mongo.InsertOneResult
			findOneRes   Map
			deleteOneRes *mongo.DeleteResult
			mp, mc       = getDB(t)
		)
		// 插入单条数据
		insertOneRes = mc.InsertOne(NewData(NewEntity("name", "张三"), NewEntity("age", 18)))
		if mc.Err != nil {
			log.Fatalf("插入单条数据失败：%v", err)
		}
		t.Logf("插入单条数据成功：%v\n", insertOneRes.InsertedID)

		// 查询单条数据
		if mc.Where(NewMap("_id", insertOneRes.InsertedID)).FindOne(&findOneRes, nil).Err != nil {
			t.Fatalf("查询单条数据失败：%v", mc.Err)
		}
		t.Logf("查询单条数据成功：%v\n", findOneRes)

		// 删除单条数据
		deleteOneRes = mc.Where(NewMap("_id", insertOneRes.InsertedID)).DeleteOne()
		if mc.Err != nil {
			t.Fatalf("删除单条数据失败：%v", mc.Err)
		}
		t.Logf("成功删除数据：%d\n", deleteOneRes.DeletedCount)

		mp.Clean()
	})
}

func Test2Many(t *testing.T) {
	t.Run("操作多条数据", func(t *testing.T) {
		var (
			err           error
			insertManyRes *mongo.InsertManyResult
			findManyRes   []Map
			deleteManyRes *mongo.DeleteResult
			mp, mc        = getDB(t)
		)
		// 插入多条数据
		insertManyRes = mc.InsertMany([]any{
			NewData(NewEntity("name", "李四"), NewEntity("age", 19)),
			NewData(NewEntity("name", "王五"), NewEntity("age", 20)),
			NewData(NewEntity("name", "赵六"), NewEntity("age", 21)),
		})
		if mc.Err != nil {
			t.Fatalf("插入多条数据失败：%v", err)
		}
		t.Logf("插入多条数据成功：%v\n", insertManyRes.InsertedIDs)

		// 查询多条数据
		if mc.Where(NewMap("_id", NewMap("$in", insertManyRes.InsertedIDs))).FindMany(&findManyRes, nil).Err != nil {
			t.Fatalf("查询多条数据失败：%v", mc.Err)
		}
		t.Logf("查询多条数据成功：%v\n", findManyRes)

		// 删除多条数据
		deleteManyRes = mc.Where(NewMap("_id", NewMap("$in", insertManyRes.InsertedIDs))).DeleteMany()
		if mc.Err != nil {
			t.Fatalf("删除多条数据失败：%v", mc.Err)
		}
		t.Logf("删除数据成功：%d\n", deleteManyRes.DeletedCount)

		mp.Clean()
	})
}

package mongoDriver

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
	if _, err = mp.AppendClient("default", mc); err != nil {
		t.Fatalf("添加mongo客户端失败：%v", err)
	}
	mc = mp.GetClient("default").SetDatabase("test_db").SetCollection("test_collection")

	return mp, mc
}

func Test1InsertOne(t *testing.T) {
	t.Run("操作单条数据", func(t *testing.T) {
		var (
			err          error
			insertOneRes *mongo.InsertOneResult
			mp, mc       = getDB(t)
		)
		// 清空数据
		_ = mc.DeleteMany(nil)
		if mc.Err != nil {
			t.Fatalf("清空数据失败：%v", err)
		}

		// 插入单条数据
		if mc.InsertOne(Map{"name": "张三", "age": 18}, &insertOneRes).Err != nil {
			log.Fatalf("插入单条数据失败：%v", err)
		}
		t.Logf("插入单条数据成功：%s\n", insertOneRes.InsertedID.(OID).String())

		mp.Clean()
	})
}

func Test2InsertMany(t *testing.T) {
	t.Run("操作多条数据", func(t *testing.T) {
		var (
			err           error
			insertManyRes *mongo.InsertManyResult
			mp, mc        = getDB(t)
		)
		// 插入多条数据
		if mc.InsertMany([]any{
			Map{"name": "李四", "age": 19},
			Map{"name": "王五", "age": 20},
			Map{"name": "赵六", "age": 21},
		}, &insertManyRes).Err != nil {
			t.Fatalf("插入多条数据失败：%v", err)
		}
		t.Logf("插入多条数据成功：%v\n", insertManyRes.InsertedIDs)

		mp.Clean()
	})
}

func Test3UpdateOne(t *testing.T) {
	var (
		updateOneRes *mongo.UpdateResult
		mp, mc       = getDB(t)
	)

	if mc.Where(Map{"name": "张三"}).UpdateOne(Map{"age": 1}, &updateOneRes).Err != nil {
		t.Fatalf("更新单条数据失败：%v", mc.Err)
	}
	t.Logf("更新成功：%d\n", updateOneRes.ModifiedCount)

	mp.Clean()
}

func Test4UpdateMany(t *testing.T) {
	var (
		updateManyRes *mongo.UpdateResult
		mp, mc        = getDB(t)
	)

	if mc.Where(Map{"name": Map{"$ne": "张三"}}).UpdateMany(Map{"age": 0}, &updateManyRes).Err != nil {
		t.Fatalf("更新单条数据失败：%v", mc.Err)
	}
	t.Logf("更新成功：%d\n", updateManyRes.ModifiedCount)

	mp.Clean()
}

func Test5DeleteOne(t *testing.T) {
	var (
		deleteOneRes *mongo.DeleteResult
		mp, mc       = getDB(t)
	)
	t.Run("删除单条数据", func(t *testing.T) {
		// 删除单条数据
		if mc.Where(Map{"name": "张三"}).DeleteOne(&deleteOneRes).Err != nil {
			t.Fatalf("删除单条数据失败：%v", mc.Err)
		}
		t.Logf("成功删除数据：%d\n", deleteOneRes.DeletedCount)

		mp.Clean()
	})
}

func Test6DeleteMany(t *testing.T) {
	var (
		deleteManyRes *mongo.DeleteResult
		mp, mc        = getDB(t)
	)

	t.Run("删除多条数据", func(t *testing.T) {
		// 删除多条数据
		if mc.Where(Map{"name": Map{"$ne": "张三"}}).DeleteMany(&deleteManyRes).Err != nil {
			t.Fatalf("删除多条数据失败：%v", mc.Err)
		}
		t.Logf("删除数据成功：%d\n", deleteManyRes.DeletedCount)

	})

	mp.Clean()
}

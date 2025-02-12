package mongoDriver

import (
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Student struct {
		Name    string `bson:"name"`
		Age     uint64 `bson:"age"`
		ClassId any    `bson:"class_id"`
		Class   Class  `bson:"class"`
	}

	Class struct {
		Name     string `bson:"name"`
		Students []Student
	}
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
			user         = Student{Name: "张三", Age: 18}
		)
		// 清空数据
		_ = mc.DeleteMany(nil)
		if mc.Err != nil {
			t.Fatalf("清空数据失败：%v", err)
		}

		// 插入单条数据
		if mc.InsertOne(user, &insertOneRes).Err != nil {
			log.Fatalf("插入单条数据失败：%v", err)
		}
		t.Logf("插入单条数据成功：%s\n", insertOneRes.InsertedID.(OID).String())

		mp.Clean()
	})
}

func Test2InsertMany(t *testing.T) {
	t.Run("操作多条数据", func(t *testing.T) {
		var (
			insertOneRes  *mongo.InsertOneResult
			insertManyRes *mongo.InsertManyResult
			mp, mc        = getDB(t)
		)
		// 插入多条数据
		if mc.SetCollection("classes").InsertOne(Class{Name: "一班"}, &insertOneRes).Err != nil {
			t.Fatalf("插入班级失败：%v", mc.Err)
		}
		t.Logf("插入班级成功：%s\n", insertOneRes.InsertedID.(OID).String())

		if mc.SetCollection("students").InsertMany([]any{
			Student{Name: "张三", Age: 18, ClassId: insertOneRes.InsertedID},
			Student{Name: "李四", Age: 19, ClassId: insertOneRes.InsertedID},
		}, &insertManyRes).Err != nil {
			t.Fatalf("插入多条数据失败：%v", mc.Err)
		}

		if mc.SetCollection("classes").InsertOne(Map{"name": "二班"}, &insertOneRes).Err != nil {
			t.Fatalf("插入班级失败：%v", mc.Err)
		}

		if mc.SetCollection("students").InsertMany([]any{
			Student{Name: "王五", Age: 20, ClassId: insertOneRes.InsertedID.(OID)},
			Student{Name: "赵六", Age: 21, ClassId: insertOneRes.InsertedID.(OID)},
		}, &insertManyRes).Err != nil {
			t.Fatalf("插入学生失败：%v", mc.Err)
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

	if mc.Where(Map{"name": "张三"}).UpdateOne(Student{Name: "张三", Age: 1}, &updateOneRes).Err != nil {
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

	if mc.SetCollection("students").Where(Map{"name": Map{"$ne": "张三"}}).UpdateMany(Map{"age": 0}, &updateManyRes).Err != nil {
		t.Fatalf("更新单条数据失败：%v", mc.Err)
	}
	t.Logf("更新成功：%d\n", updateManyRes.ModifiedCount)

	mp.Clean()
}

func Test5FindOne(t *testing.T) {
	var (
		student *Student
		mp, mc  = getDB(t)
	)

	t.Run("查询单条数据", func(t *testing.T) {
		if mc.SetCollection("students").Where(Map{"name": "张三"}).FindOne(&student, nil).Err != nil {
			t.Fatalf("查询单条数据失败：%v", mc.Err)
		}
		t.Errorf("查询单条数据成功：%#v\n", student)

		mp.Clean()
	})
}

func Test6FindMany(t *testing.T) {
	var (
		// students []Student
		// classA Class
		classes []Map
		mp, mc  = getDB(t)
	)

	t.Run("查询多条数据", func(t *testing.T) {
		if mc.SetCollection("classes").Where(Map{
			"$lookup": Map{
				"from":         "students",
				"localField":   "_id",
				"foreignField": "class_id",
				"as":           "students",
			},
		}, Map{
			"$match": Map{"name": "一班"},
		}).Aggregate(&classes).Err != nil {
			t.Fatalf("查询多条数据失败：%v", mc.Err)
		}
		t.Logf("查询多条数据成功：%#v\n", classes)
	})

	mp.Clean()
}

func Test7DeleteOne(t *testing.T) {
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

func Test8DeleteMany(t *testing.T) {
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

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 设置客户端选项
	clientOptions := options.Client().ApplyURI("mongodb://admin:admin@localhost:27017")

	// 连接到 MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("成功连接到 MongoDB!")

	// 使用特定的数据库和集合
	collection := client.Database("test_db").Collection("test_collection")

	// 插入新数据
	newData := bson.D{
		{Key: "name", Value: "Alice"},
		{Key: "age", Value: 25},
		{Key: "city", Value: "New York"},
	}

	insertResult, err := collection.InsertOne(context.TODO(), newData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("插入数据的ID:", insertResult.InsertedID)

	// 根据ID搜索数据
	var result bson.M
	id := insertResult.InsertedID.(primitive.ObjectID)
	if err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("OK: %v\n", result)

	// 搜索多条信息
	filter := bson.D{{Key: "age", Value: bson.D{{Key: "$gt", Value: 20}}}} // 搜索年龄大于20的记录
	findOptions := options.Find().SetLimit(5)                              // 设置返回记录的最大数量
	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	var results []bson.M
	cursor.All(context.TODO(), &results)
	// for cursor.Next(context.TODO()) {
	// 	var elem bson.M
	// 	err := cursor.Decode(&elem)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	results = append(results, elem)
	// }

	// if err := cursor.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("搜索到的记录:")
	for _, result := range results {
		fmt.Println(result)
	}

	// // 删除单条记录
	// deleteResult, err := collection.DeleteOne(context.TODO(), bson.M{"name": "Alice"})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("删除了 %v 条记录\n", deleteResult.DeletedCount)

	// // 删除多条记录
	// deleteManyResult, err := collection.DeleteMany(context.TODO(), bson.M{"age": bson.M{"$gt": 20}})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("删除了 %v 条记录\n", deleteManyResult.DeletedCount)

	// 断开与 MongoDB 的连接
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("已断开与 MongoDB 的连接。")
}

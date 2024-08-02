package mongodb

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestMogoDB(t *testing.T) {
	// 控制初始化超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, e *event.CommandStartedEvent) {
			fmt.Println(e.Command)
		},
		Succeeded: func(ctx context.Context, e *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, e *event.CommandFailedEvent) {

		},
	}
	// 首先初始化一个mongoDB客户端
	opts := options.Client().ApplyURI("mongodb://root:example@localhost:27017").SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)

	// 创建一个mongoDB的数据库
	mdb := client.Database("we_book")
	// 创建一个mongoDB的集合
	col := mdb.Collection("article")

	//defer func() {
	//	_, err = col.DeleteMany(ctx, bson.D{})
	//}()

	// 插入一条文档
	res, err := col.InsertOne(ctx, Article{
		Id:       1,
		Title:    "我的标题",
		Content:  "我的内容",
		AuthorId: 12,
		Status:   1,
		Ctime:    time.Now().UnixMilli(),
		Utime:    time.Now().UnixMilli(),
	})
	if err != nil {
		panic(err)
	}
	// InsertedID 返回插入的文档id，并非是article的id
	fmt.Printf("插入结果的ID：%d\n", res.InsertedID)

	// 查找 Id = 1 的文章
	filter := bson.D{bson.E{Key: "id", Value: 1}}
	var art Article
	err = col.FindOne(ctx, filter).Decode(&art)
	assert.NoError(t, err)
	fmt.Printf("使用id进行查找结果：%+v\n", art)
	//查找结果：{Id:1 Title:我的标题 Content:我的内容 AuthorId:12 Status:1 Ctime:1722564149250 Utime:1722564149250}

	art = Article{}
	err = col.FindOne(ctx, Article{Title: "我的标题"}).Decode(&art)
	if err == mongo.ErrNoDocuments {
		fmt.Println("没有数据")
	}
	assert.NoError(t, err)
	fmt.Printf("使用title进行查找结果：%+v\n", art)

	// 更新文档
	sets := bson.D{bson.E{Key: "$set", Value: bson.E{Key: "title", Value: "新的标题"}}}
	updateRes, err := col.UpdateOne(ctx, filter, sets)
	assert.NoError(t, err)
	fmt.Printf("更新后的文档数量：%d\n", updateRes.ModifiedCount)
	// 更新后的文档数量：1

	updateRes1, err := col.UpdateOne(ctx, filter, bson.D{
		bson.E{Key: "$set", Value: Article{Title: "新的标题2", AuthorId: 123456}}})
	assert.NoError(t, err)
	fmt.Printf("更新后的文档数量：%d\n", updateRes1.ModifiedCount)
	// 更新后的文档数量：1

	// 删除文档
	//filter = bson.D{bson.E{Key: "id", Value: 1}}
	//delRes, err := col.DeleteOne(ctx, filter)
	//assert.NoError(t, err)
	//fmt.Printf("删除后的文档数量：%d\n", delRes.DeletedCount)
	// 删除后的文档数量：1

	// 使用 Or 查询条件
	or := bson.A{bson.M{"id": 1}, bson.M{"id": 123}}
	orRes, err := col.Find(ctx, bson.D{bson.E{"$or", or}})
	assert.NoError(t, err)
	var arts []Article
	err = orRes.All(ctx, &arts)
	assert.NoError(t, err)
	fmt.Printf("使用or进行查找结果：%+v\n", arts)
	// 使用or进行查找结果：[{Id:1 Title:新的标题2 Content:我的内容 AuthorId:123456 Status:1 Ctime:1722576620900 Utime:1722576620900}]

	// 使用 And 查询条件
	and := bson.A{bson.D{bson.E{"id", 1}},
		bson.D{bson.E{"title", "新的标题2"}}}
	andRes, err := col.Find(ctx, bson.D{bson.E{"$and", and}})
	assert.NoError(t, err)
	ars := []Article{}
	err = andRes.All(ctx, &ars)
	assert.NoError(t, err)
	fmt.Printf("使用and进行查找结果：%+v\n", ars)
	// 使用and进行查找结果：[{Id:1 Title:新的标题2 Content:我的内容 AuthorId:123456 Status:1 Ctime:1722576444390 Utime:1722576444390}]

	// IN 查询
	in := bson.D{bson.E{Key: "id", Value: bson.M{"$in": []any{1, 2, 3}}}}
	inRes, err := col.Find(ctx, in)
	ars = []Article{}
	err = inRes.All(ctx, &ars)
	assert.NoError(t, err)
	fmt.Printf("使用in查找结果：%+v\n", ars)
	// 使用in查找结果：[{Id:1 Title:新的标题2 Content:我的内容 AuthorId:123456 Status:1 Ctime:1722576444390 Utime:1722576444390}]

	inRes, err = col.Find(ctx, in, options.Find().SetProjection(bson.M{
		"id":    1,
		"title": "我的标题",
	}))
	ars = []Article{}
	err = inRes.All(ctx, &ars)
	assert.NoError(t, err)
	fmt.Printf("自定义查询的查找结果：%+v\n", ars)
	// 自定义查询的查找结果：[{Id:1 Title:我的标题 Content: AuthorId:0 Status:0 Ctime:0 Utime:0}]

	idxRes, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"author_id": 1},
		},
	})
	assert.NoError(t, err)
	fmt.Printf("创建索引结果：%+v\n", idxRes)
	// 创建索引结果：[id_1 author_id_1]

	delRes, err := col.DeleteMany(ctx, filter)
	assert.NoError(t, err)
	fmt.Printf("删除后的文档数量：%d\n", delRes.DeletedCount)
	// 删除后的文档数量：1
}

type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}

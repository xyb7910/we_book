package article

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBDAO struct {
	// 制作库
	col *mongo.Collection
	// 线上库
	liveCol *mongo.Collection
	node    *snowflake.Node
	idGen   IDGenerator
}

func (m *MongoDBDAO) UpdateById(ctx context.Context, art Article) error {
	filter := bson.M{"id": art.Id, "author_id": art.AuthorId}
	updates := bson.D{bson.E{"$set", bson.M{
		"title":   art.Title,
		"content": art.Content,
		"status":  art.Status,
		"utime":   time.Now().UnixMilli(),
	}}}
	res, err := m.col.UpdateOne(ctx, filter, updates)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("update fail")
	}
	return nil
}

func (m *MongoDBDAO) Sync(ctx context.Context, art Article) (int64, error) {
	// 第一步，保存到制作库
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = m.UpdateById(ctx, art)
	} else {
		id, err = m.Insert(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	// 操作线上库
	now := time.Now().UnixMilli()
	art.Utime = now
	update := bson.M{
		"$set":         PublishedArticle(art),
		"$setOnInsert": bson.M{"ctime": now},
	}
	filter := bson.M{"id": art.Id}
	_, err = m.liveCol.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return id, err
}

func InitCollection(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	index := []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{bson.E{Key: "author_id", Value: 1},
				bson.E{Key: "ctime", Value: 1},
			},
			Options: options.Index(),
		},
	}
	_, err := db.Collection("article").Indexes().CreateMany(ctx, index)
	if err != nil {
		return err
	}
	_, err = db.Collection("live_articles").Indexes().CreateMany(ctx, index)
	return err
}

func (m *MongoDBDAO) Upsert(ctx context.Context, article Article) error {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBDAO) SyncStatus(ctx context.Context, id int64, author int64, u uint8) error {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBDAO) Transaction(ctx context.Context, bizFunc func(txDAO ArticleDAO) error) error {
	//TODO implement me
	panic("implement me")
}

type IDGenerator int64

func NewMongoDBDAO(db *mongo.Database, node *snowflake.Node) ArticleDAO {
	return &MongoDBDAO{
		col:     db.Collection("articles"),
		liveCol: db.Collection("live_articles"),
		node:    node,
	}
}

func (m *MongoDBDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	// ID 生成使用 snowflake
	id := m.node.Generate().Int64()
	art.Id = id
	_, err := m.col.InsertOne(ctx, art)
	return id, err
}

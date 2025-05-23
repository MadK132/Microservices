package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *MongoStatsRepo) GetAll() (map[string]map[string]int, error) {
	ctx := context.Background()

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "source", Value: "$source"},
				{Key: "action", Value: "$action"},
			}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[string]map[string]int)

	for cursor.Next(ctx) {
		var doc struct {
			ID struct {
				Source string `bson:"source"`
				Action string `bson:"action"`
			} `bson:"_id"`
			Count int `bson:"count"`
		}

		if err := cursor.Decode(&doc); err != nil {
			continue
		}

		if _, ok := result[doc.ID.Source]; !ok {
			result[doc.ID.Source] = make(map[string]int)
		}
		result[doc.ID.Source][doc.ID.Action] = doc.Count
	}

	return result, nil
}


type MongoStatsRepo struct {
	col *mongo.Collection
}


func NewMongoStatsRepo(db *mongo.Database) *MongoStatsRepo {
	return &MongoStatsRepo{
		col: db.Collection("events"),
	}
}

func (r *MongoStatsRepo) Save(source, action, ts string) error {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		t = time.Now()
	}
	event := Event{
		Source:    source,
		Action:    action,
		Timestamp: t,
	}
	_, err = r.col.InsertOne(context.Background(), event)
	return err
}

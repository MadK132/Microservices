package dao

import (
	"context"
	"errors"
	"log"

	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewDAO struct {
	col *mongo.Collection
}

func NewReviewDAO(db *mongo.Database) *ReviewDAO {
	return &ReviewDAO{
		col: db.Collection("review"),
	}
}

func (d *ReviewDAO) CreateReview(ctx context.Context, review model.Review) (primitive.ObjectID, error){
	review.ID = primitive.NewObjectID()

	res, err := d.col.InsertOne(ctx, review)
	if err != nil {
		return primitive.NilObjectID, err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("failed to cast inserted ID to ObjectID")
	}

	log.Printf("Created product with ID: %s", oid.Hex())
	return oid, nil
}

func (d *ReviewDAO) GetReviewByID(ctx context.Context, id primitive.ObjectID) (model.Review, error){
	log.Printf("Attempting to get product with ID: %s", id.Hex())

	var review model.Review

	filter := bson.M{
		"_id":       id,
	}
	log.Printf("MongoDB query filter: %+v", filter)

	err := d.col.FindOne(ctx, filter).Decode(&review)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("Product not found with ID: %s", id.Hex())
			return model.Review{}, model.ErrProductNotFound
		}
		log.Printf("Error retrieving product: %v", err)
		return model.Review{}, err
	}

	log.Printf("Successfully retrieved product with ID: %s", id.Hex())
	return review, nil
}

func (d *ReviewDAO) UpdateReview(ctx context.Context, id primitive.ObjectID, update model.ReviewUpdate) error{
	updateData := bson.M{}
	if update.Rating != nil {
		updateData["rating"] = *update.Rating
	}
	if update.Comment != nil {
		updateData["comment"] = *update.Comment
	}

	_, err := d.col.UpdateOne(ctx, bson.M{
		"_id":       id,
	}, bson.M{
		"$set": updateData,
	})

	return err
} 
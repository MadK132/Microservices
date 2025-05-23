package dao

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductDAO struct {
	col *mongo.Collection
}

func NewProductDAO(db *mongo.Database) *ProductDAO {
	return &ProductDAO{
		col: db.Collection("products"),
	}
}

func (d *ProductDAO) Create(ctx context.Context, product model.Product) (primitive.ObjectID, error) {
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	product.IsDeleted = false

	res, err := d.col.InsertOne(ctx, product)
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

func (d *ProductDAO) GetByID(ctx context.Context, id primitive.ObjectID) (model.Product, error) {
	log.Printf("Attempting to get product with ID: %s", id.Hex())

	var product model.Product

	filter := bson.M{
		"_id":       id,
		"isdeleted": false,
	}
	log.Printf("MongoDB query filter: %+v", filter)

	err := d.col.FindOne(ctx, filter).Decode(&product)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("Product not found with ID: %s", id.Hex())
			return model.Product{}, model.ErrProductNotFound
		}
		log.Printf("Error retrieving product: %v", err)
		return model.Product{}, err
	}

	log.Printf("Successfully retrieved product with ID: %s", id.Hex())
	return product, nil
}

func (d *ProductDAO) Update(ctx context.Context, id primitive.ObjectID, update model.ProductUpdate) error {
	updateData := bson.M{}
	if update.Name != nil {
		updateData["name"] = *update.Name
	}
	if update.Description != nil {
		updateData["description"] = *update.Description
	}
	if update.Price != nil {
		updateData["price"] = *update.Price
	}
	if update.Stock != nil {
		updateData["stock"] = *update.Stock
	}
	if update.Category != nil {
		updateData["category"] = *update.Category
	}
	if update.IsDeleted != nil {
		updateData["isdeleted"] = *update.IsDeleted
	}

	updateData["updatedat"] = time.Now()

	_, err := d.col.UpdateOne(ctx, bson.M{
		"_id":       id,
		"isdeleted": false,
	}, bson.M{
		"$set": updateData,
	})

	return err
}

func (d *ProductDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := d.col.UpdateOne(ctx, bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"isdeleted": true,
			"updatedAt": time.Now(),
		},
	})
	return err
}

func (d *ProductDAO) GetAll(ctx context.Context, filter model.ProductFilter) ([]model.Product, error) {
	query := bson.M{
		"isdeleted": false,
	}

	if filter.Name != nil {
		query["name"] = *filter.Name
	}
	if filter.Category != nil {
		query["category"] = *filter.Category
	}
	if filter.MinPrice != nil || filter.MaxPrice != nil {
		priceFilter := bson.M{}
		if filter.MinPrice != nil {
			priceFilter["$gte"] = *filter.MinPrice
		}
		if filter.MaxPrice != nil {
			priceFilter["$lte"] = *filter.MaxPrice
		}
		query["price"] = priceFilter
	}

	cursor, err := d.col.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []model.Product
	for cursor.Next(ctx) {
		var product model.Product
		if err := cursor.Decode(&product); err == nil {
			products = append(products, product)
		}
	}
	return products, nil
}

func (d *ProductDAO) GetByIDs(ctx context.Context, ids []primitive.ObjectID) ([]model.Product, error) {
	filter := bson.M{
		"_id":       bson.M{"$in": ids},
		"isdeleted": false,
	}

	cursor, err := d.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []model.Product
	for cursor.Next(ctx) {
		var product model.Product
		if err := cursor.Decode(&product); err == nil {
			products = append(products, product)
		}
	}
	return products, nil
}

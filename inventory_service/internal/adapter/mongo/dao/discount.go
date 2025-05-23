package dao

import (
	"context"
	"errors"
	_"fmt"
	"log"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DiscountDAO struct {
	col *mongo.Collection
}

func NewDiscountDAO(db *mongo.Database) *DiscountDAO{
	return &DiscountDAO{
		col: db.Collection("discount"),
	}
}

func (d *DiscountDAO) CreateDiscnout(ctx context.Context, promotion model.Discount) (primitive.ObjectID, error){
	promotion.ID = primitive.NewObjectID()
	promotion.StartDate = time.Now()
	promotion.IsActive = true

	res, err := d.col.InsertOne(ctx, promotion)
	if err != nil {
		return primitive.NilObjectID, err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok{
		return primitive.NilObjectID, errors.New("failed to cast inserted id to ObjectId") 

	}
	log.Printf("Create product with ID: %s", oid.Hex())
	return oid, nil
}

func (d *DiscountDAO) GetAllExistingDiscounts(ctx context.Context, filter model.DiscountFilter) ([]model.Discount, error){
	query := bson.M{
		"isactive" : true,
	}

	if filter.Name != nil {
		query["name"] = *filter.Name
	}
	if filter.Description != nil {
		query["description"] = *filter.Description
	}

	cursor, err := d.col.Find(ctx, query)
	if err != nil{
		return nil, err
	}

	defer cursor.Close(ctx)

	var discounts []model.Discount
	for cursor.Next(ctx){
		var discount model.Discount
		if err := cursor.Decode(&discount); err == nil{
			discounts = append(discounts, discount)
		}
	}
	return discounts, nil 
}

func (d *DiscountDAO) GetAllProductsWithDiscounts(ctx context.Context, filter model.DiscountFilter) ([]model.Product, error) {
	query := bson.M{"isactive": true}

	if filter.ID != nil {
		query["_id"] = *filter.ID
	}

	cursor, err := d.col.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var promotions []model.Discount
	if err := cursor.All(ctx, &promotions); err != nil {
		return nil, err
	}

	productMap := make(map[primitive.ObjectID]struct{})
	for _, promo := range promotions {
		for _, pid := range promo.ApplicableProducts {
			productMap[pid] = struct{}{}
		}
	}

	var productIDs []primitive.ObjectID
	for id := range productMap {
		productIDs = append(productIDs, id)
	}

	productCol := d.col.Database().Collection("products") 

	cursor, err = productCol.Find(ctx, bson.M{
		"_id":       bson.M{"$in": productIDs},
		"isdeleted": false,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []model.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (d *DiscountDAO) DeletePromotion(ctx context.Context, id primitive.ObjectID) error {
	_, err := d.col.DeleteOne(ctx, bson.M{"_id": id}) 
	return err
}

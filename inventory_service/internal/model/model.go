package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string
	Description string
	Price       float64
	Stock       int
	Category    string

	CreatedAt time.Time
	UpdatedAt time.Time

	IsDeleted bool `bson:"isdeleted"`
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return ErrInvalidName
	}
	if p.Price <= 0 {
		return ErrInvalidPrice
	}
	if p.Stock < 0 {
		return ErrInvalidStock
	}
	if p.Description == "" {
		return ErrInvalidDescription
	}
	if !isValidCategory(p.Category) {
		return ErrInvalidCategory
	}
	return nil
}

func isValidCategory(category string) bool {
	validCategories := []string{
		CategoryMilk,
		CategoryDrink,
		CategorySnack,
		CategoryFood,
		CategoryFruit,
	}
	for _, valid := range validCategories {
		if category == valid {
			return true
		}
	}
	return false
}

type ProductFilter struct {
	ID       *primitive.ObjectID
	Name     *string
	Price    *float64
	MinPrice *float64
	MaxPrice *float64
	Category *string

	IsDeleted *bool
}

type ProductUpdate struct {
	ID          *primitive.ObjectID
	Name        *string
	Description *string
	Price       *float64
	Stock       *int
	Category    *string

	UpdatedAt *time.Time

	IsDeleted *bool
}

func (p *ProductUpdate) Validate() error {
	if p.Name != nil && *p.Name == "" {
		return ErrInvalidName
	}
	if p.Price != nil && *p.Price <= 0 {
		return ErrInvalidPrice
	}
	if p.Stock != nil && *p.Stock < 0 {
		return ErrInvalidStock
	}
	if p.Description != nil && *p.Description == "" {
		return ErrInvalidDescription
	}
	if p.Category != nil && !isValidCategory(*p.Category) {
		return ErrInvalidCategory
	}
	return nil
}

type Discount struct {
    ID                primitive.ObjectID  
    Name              string    
    Description       string   
    DiscountPercentage float64  
    ApplicableProducts []primitive.ObjectID 

    StartDate         time.Time 
    EndDate           time.Time 

    IsActive          bool      
}

type DiscountFilter struct {
	ID *primitive.ObjectID
	Name *string
	Description *string
	DiscountPercentage *string
	ApplicableProducts *[]primitive.ObjectID

	StartDate *time.Time
	EndDate *time.Time

	IsActive *bool
}

type Review struct {
    ID         primitive.ObjectID    // or string if uuid
    ProductID  primitive.ObjectID    // or string if uuid
    UserID     primitive.ObjectID    // or string if uuid
    Rating     float64       // 1 to 5, must be calculated 
    Comment    string    
}

type ReviewUpdate struct {
	ID         *primitive.ObjectID    // or string if uuid
    ProductID  *primitive.ObjectID    // or string if uuid
    UserID     *primitive.ObjectID    // or string if uuid
    Rating     *float64       // 1 to 5, must be calculated 
    Comment    *string 
}

type ReviewFilter struct {
	ID         *primitive.ObjectID    // or string if uuid
    ProductID  *primitive.ObjectID    // or string if uuid
    UserID     *primitive.ObjectID    // or string if uuid
    Rating     *float64       // 1 to 5, must be calculated 
    Comment    *string 
}



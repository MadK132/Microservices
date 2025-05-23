package handler

import (
	"net/http"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/http/service/handler/dto"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"

	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	uc ProductUsecase
}

func NewProduct(uc ProductUsecase) *Product {
	return &Product{
		uc: uc,
	}
}

func (h *Product) Create(ctx *gin.Context) {
	product, err := dto.FromProductCreateRequest(ctx)
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	if err := product.Validate(); err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	newProduct, err := h.uc.Create(ctx.Request.Context(), product)
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToProductCreateResponse(newProduct))
}

func (h *Product) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	log.Printf("Received GetByID request with ID string: %s", id)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Failed to convert ID string to ObjectID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid object id"})
		return
	}
	log.Printf("Successfully converted ID string to ObjectID: %s", objectID.Hex())

	product, err := h.uc.GetByID(ctx.Request.Context(), model.ProductFilter{
		ID: &objectID,
	})

	if err != nil {
		errCtx := dto.FromError(err)
		log.Printf("Error retrieving product: %v", errCtx.Message)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	log.Printf("Successfully retrieved product with ID: %s", product.ID.Hex())
	ctx.JSON(http.StatusOK, dto.ToProductCreateResponse(product))
}

func (h *Product) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid object id"})
		return
	}

	updateData, err := dto.FromProductUpdateRequest(ctx)
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	if err := updateData.Validate(); err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	updateData.ID = &objectID

	updated, err := h.uc.Update(ctx.Request.Context(), updateData)
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToProductCreateResponse(updated))
}

func (h *Product) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid object id"})
		return
	}

	err = h.uc.Delete(ctx.Request.Context(), model.ProductFilter{ID: &objectID})
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (h *Product) GetAll(ctx *gin.Context) {
	products, err := h.uc.GetAll(ctx.Request.Context(), model.ProductFilter{})
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	ctx.JSON(http.StatusOK, products)
}
func (h *Product) GetByIDs(ctx *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var objectIDs []primitive.ObjectID
	for _, idStr := range req.IDs {
		objID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid object id: " + idStr})
			return
		}
		objectIDs = append(objectIDs, objID)
	}

	products, err := h.uc.GetByIDs(ctx.Request.Context(), objectIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var res []dto.ProductCreateResponse
	for _, p := range products {
		res = append(res, dto.ToProductCreateResponse(p))
	}

	ctx.JSON(http.StatusOK, res)
}
package handler

import (
	"net/http"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/http/service/handler/dto"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Discount struct {
	uc DiscountUsecase
}

func NewDiscount(uc DiscountUsecase) *Discount {
	return &Discount{
		uc: uc,
	}
}

func (h *Discount) CreateDiscnout(ctx *gin.Context){
	discount, err := dto.FromDiscountCreateRequest(ctx)
	if err != nil{
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	newDiscount, err := h.uc.CreateDiscnout(ctx.Request.Context(), discount)
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToDiscountCreateResponse(newDiscount))
}

func (h *Discount) GetAllExistingDiscounts(ctx *gin.Context) {
	discounts, err := h.uc.GetAllExistingDiscounts(ctx.Request.Context(), model.DiscountFilter{
		IsActive: ptrBool(true),
	})
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	var response []dto.DiscountCreateResponse
	for _,p := range discounts{
		response = append(response, dto.ToDiscountCreateResponse(p))
	}

	ctx.JSON(http.StatusOK, response)
}
func (h *Discount) GetAllProductsWithDiscounts(ctx *gin.Context) {
	products, err := h.uc.GetAllProductsWithDiscounts(ctx.Request.Context(), model.DiscountFilter{
		IsActive: ptrBool(true),
	})
	
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	var response []dto.ProductCreateResponse
	for _, p := range products {
		response = append(response, dto.ToProductCreateResponse(p))
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *Discount) DeleteDiscount(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ObjectID"})
		return
	}

	err = h.uc.DeleteDiscount(ctx.Request.Context(), model.DiscountFilter{ID: &id})
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "promotion deleted"})
}

func ptrBool(b bool) *bool {
	return &b
}

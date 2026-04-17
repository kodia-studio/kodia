package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/pagination"
	"github.com/kodia-studio/kodia/pkg/response"
	"go.uber.org/zap"
)

// ProductHandler handles HTTP requests for Product.
type ProductHandler struct {
	service  ports.ProductService
	validate *validator.Validate
	log      *zap.Logger
}

func NewProductHandler(service ports.ProductService, validate *validator.Validate, log *zap.Logger) *ProductHandler {
	return &ProductHandler{
		service:  service,
		validate: validate,
		log:      log,
	}
}

// GetAll godoc
// @Summary      List all products
// @Tags         products
// @Security     BearerAuth
// @Param        page     query int false "Page number" default(1)
// @Param        per_page query int false "Items per page" default(15)
// @Success      200 {object} response.Response
// @Router       /api/products [get]
func (h *ProductHandler) GetAll(c *gin.Context) {
	params := pagination.FromContext(c)
	
	items, total, err := h.service.GetAll(c.Request.Context(), params)
	if err != nil {
		h.log.Error("Failed to get products", zap.Error(err))
		response.InternalServerError(c, "")
		return
	}

	meta := response.NewMeta(params.Page, params.PerPage, total)
	response.OKWithMeta(c, "Success", items, meta)
}

// GetByID godoc
// @Summary      Get product by ID
// @Tags         products
// @Security     BearerAuth
// @Param        id path string true "ID"
// @Success      200 {object} response.Response
// @Router       /api/products/{id} [get]
func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	item, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "Product not found")
			return
		}
		response.InternalServerError(c, "")
		return
	}
	response.OK(c, "Success", item)
}

// Create godoc
// @Summary      Create product
// @Tags         products
// @Security     BearerAuth
// @Success      201 {object} response.Response
// @Router       /api/products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	// TODO: Define CreateRequest DTO and bind
	// var req dto.CreateProductRequest
	// if err := c.ShouldBindJSON(&req); err != nil { ... }
	
	c.JSON(http.StatusCreated, gin.H{"message": "Not implemented"})
}

// Update godoc
// @Summary      Update product
// @Tags         products
// @Security     BearerAuth
// @Param        id path string true "ID"
// @Success      200 {object} response.Response
// @Router       /api/products/{id} [patch]
func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")
	// TODO: Define UpdateRequest DTO and bind
	_ = id
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented"})
}

// Delete godoc
// @Summary      Delete product
// @Tags         products
// @Security     BearerAuth
// @Param        id path string true "ID"
// @Success      204
// @Router       /api/products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "Product not found")
			return
		}
		response.InternalServerError(c, "")
		return
	}
	response.NoContent(c)
}

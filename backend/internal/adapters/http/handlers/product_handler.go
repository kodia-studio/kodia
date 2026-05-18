package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/internal/adapters/http/transformers"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/services"
	"github.com/kodia-studio/kodia/pkg/response"
)

// ProductHandler handles all product-related HTTP requests.
type ProductHandler struct {
	productService *services.ProductService
}

// NewProductHandler creates a new ProductHandler.
func NewProductHandler(svc *services.ProductService) *ProductHandler {
	return &ProductHandler{productService: svc}
}

// ListPublic godoc
// @Summary List published products
// @Tags products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(15)
// @Success 200 {array} transformers.ProductResponse
// @Router /products [get]
func (h *ProductHandler) ListPublic(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 15
	}

	products, total, err := h.productService.ListPublic(c.Request.Context(), page, perPage)
	if err != nil {
		response.InternalServerError(c, "Failed to list products")
		return
	}

	response.OKWithMeta(c, "Products retrieved", transformers.ProductsToResponse(products),
		response.NewMeta(page, perPage, total))
}

// GetBySlug godoc
// @Summary Get product by slug
// @Tags products
// @Produce json
// @Param slug path string true "Product slug"
// @Success 200 {object} transformers.ProductResponse
// @Router /products/{slug} [get]
func (h *ProductHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	product, err := h.productService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		response.NotFound(c, "Product not found")
		return
	}
	response.OK(c, "Product retrieved", transformers.ProductToResponse(product))
}

// ListAdmin godoc
// @Summary List all products (admin)
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {array} transformers.ProductResponse
// @Router /admin/products [get]
func (h *ProductHandler) ListAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "50"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 200 {
		perPage = 50
	}

	products, total, err := h.productService.ListAdmin(c.Request.Context(), page, perPage)
	if err != nil {
		response.InternalServerError(c, "Failed to list products")
		return
	}

	response.OKWithMeta(c, "Products retrieved", transformers.ProductsToResponse(products),
		response.NewMeta(page, perPage, total))
}

// GetByID godoc
// @Summary Get product by ID (admin)
// @Tags admin
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} transformers.ProductResponse
// @Router /admin/products/{id} [get]
func (h *ProductHandler) GetByID(c *gin.Context) {
	product, err := h.productService.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.NotFound(c, "Product not found")
		return
	}
	response.OK(c, "Product retrieved", transformers.ProductToResponse(product))
}

type createProductRequest struct {
	Name            string   `json:"name" binding:"required"`
	Slug            string   `json:"slug"`
	Tagline         string   `json:"tagline"`
	Description     string   `json:"description" binding:"required"`
	Type            string   `json:"type" binding:"required"`
	ServiceType     string   `json:"service_type"`
	CoverURL        string   `json:"cover_url"`
	Tags            []string `json:"tags"`
	IsPublished     bool     `json:"is_published"`
	MetaTitle       string   `json:"meta_title"`
	MetaDescription string   `json:"meta_description"`
	OGImageURL      string   `json:"og_image_url"`
}

// Create godoc
// @Summary Create a product (admin)
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Param body body createProductRequest true "Product data"
// @Success 201 {object} transformers.ProductResponse
// @Router /admin/products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	var st *domain.ServiceType
	if req.ServiceType != "" {
		s := domain.ServiceType(req.ServiceType)
		st = &s
	}

	product := &domain.Product{
		Slug:            req.Slug,
		Name:            req.Name,
		Tagline:         req.Tagline,
		Description:     req.Description,
		Type:            domain.ProductType(req.Type),
		ServiceType:     st,
		CoverURL:        req.CoverURL,
		Tags:            req.Tags,
		IsPublished:     req.IsPublished,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		OGImageURL:      req.OGImageURL,
	}

	if err := h.productService.Create(c.Request.Context(), product); err != nil {
		response.InternalServerError(c, "Failed to create product")
		return
	}

	full, _ := h.productService.GetByID(c.Request.Context(), product.ID)
	response.Created(c, "Product created", transformers.ProductToResponse(full))
}

// Update godoc
// @Summary Update a product (admin)
// @Tags admin
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Router /admin/products/{id} [patch]
func (h *ProductHandler) Update(c *gin.Context) {
	var fields map[string]any
	if err := c.ShouldBindJSON(&fields); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Cast type fields properly
	if t, ok := fields["type"].(string); ok {
		fields["type"] = domain.ProductType(t)
	}
	if st, ok := fields["service_type"].(string); ok {
		fields["service_type"] = domain.ServiceType(st)
	}

	product, err := h.productService.Update(c.Request.Context(), c.Param("id"), fields)
	if err != nil {
		response.NotFound(c, "Product not found")
		return
	}

	full, _ := h.productService.GetByID(c.Request.Context(), product.ID)
	response.OK(c, "Product updated", transformers.ProductToResponse(full))
}

// Delete godoc
// @Summary Delete a product (admin)
// @Tags admin
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Router /admin/products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	if err := h.productService.Delete(c.Request.Context(), c.Param("id")); err != nil {
		response.NotFound(c, "Product not found")
		return
	}
	response.NoContent(c)
}

type addVariantRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	IsActive    bool   `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
}

// AddVariant godoc
// @Summary Add a variant to a product (admin)
// @Tags admin
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Router /admin/products/{id}/variants [post]
func (h *ProductHandler) AddVariant(c *gin.Context) {
	var req addVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	variant := &domain.ProductVariant{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
	}

	if err := h.productService.AddVariant(c.Request.Context(), c.Param("id"), variant); err != nil {
		response.InternalServerError(c, "Failed to add variant")
		return
	}

	response.Created(c, "Variant added", transformers.VariantToResponse(variant))
}

// UpdateVariant godoc
// @Summary Update a product variant (admin)
// @Tags admin
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param variantId path string true "Variant ID"
// @Router /admin/products/{id}/variants/{variantId} [patch]
func (h *ProductHandler) UpdateVariant(c *gin.Context) {
	var fields map[string]any
	if err := c.ShouldBindJSON(&fields); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Cast numeric fields correctly
	if v, ok := fields["price"].(float64); ok {
		fields["price"] = int64(v)
	}
	if v, ok := fields["promo_price"].(float64); ok {
		fields["promo_price"] = int64(v)
	}
	if v, ok := fields["sort_order"].(float64); ok {
		fields["sort_order"] = int(v)
	}

	variant, err := h.productService.UpdateVariant(c.Request.Context(), c.Param("id"), c.Param("variantId"), fields)
	if err != nil {
		response.NotFound(c, "Variant not found")
		return
	}

	response.OK(c, "Variant updated", transformers.VariantToResponse(variant))
}

// DeleteVariant godoc
// @Summary Delete a product variant (admin)
// @Tags admin
// @Security BearerAuth
// @Router /admin/products/{id}/variants/{variantId} [delete]
func (h *ProductHandler) DeleteVariant(c *gin.Context) {
	if err := h.productService.DeleteVariant(c.Request.Context(), c.Param("id"), c.Param("variantId")); err != nil {
		response.NotFound(c, "Variant not found")
		return
	}
	response.NoContent(c)
}

// UploadVariantFile godoc
// @Summary Upload file for a variant (admin)
// @Tags admin
// @Security BearerAuth
// @Router /admin/products/{id}/variants/{variantId}/upload [post]
func (h *ProductHandler) UploadVariantFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "File is required", nil)
		return
	}

	fileKey, err := h.productService.UploadVariantFile(c.Request.Context(), c.Param("id"), c.Param("variantId"), file)
	if err != nil {
		response.InternalServerError(c, "Failed to upload file")
		return
	}

	response.OK(c, "File uploaded", gin.H{"file_key": fileKey})
}

// UploadCover godoc
// @Summary Upload cover image for a product (admin)
// @Tags admin
// @Security BearerAuth
// @Router /admin/products/{id}/cover [post]
func (h *ProductHandler) UploadCover(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "File is required", nil)
		return
	}

	publicURL, err := h.productService.UploadCover(c.Request.Context(), c.Param("id"), file)
	if err != nil {
		response.InternalServerError(c, "Failed to upload cover")
		return
	}

	response.OK(c, "Cover uploaded", gin.H{"url": publicURL})
}

// UploadPublic godoc
// @Summary Upload a public file (returns URL)
// @Tags admin
// @Security BearerAuth
// @Router /admin/upload/public [post]
func (h *ProductHandler) UploadPublic(c *gin.Context) {
	// Placeholder — reuse cover upload pattern for generic public files
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Use /admin/products/{id}/cover for product covers"})
}

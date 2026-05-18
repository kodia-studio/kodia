package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kodia-studio/kodia/internal/adapters/http/dto"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/services"
	"github.com/kodia-studio/kodia/pkg/response"
)

// ProductDocHandler handles product documentation HTTP requests.
type ProductDocHandler struct {
	productDocService *services.ProductDocService
	productService    *services.ProductService
}

// NewProductDocHandler creates a new ProductDocHandler.
func NewProductDocHandler(docSvc *services.ProductDocService, prodSvc *services.ProductService) *ProductDocHandler {
	return &ProductDocHandler{
		productDocService: docSvc,
		productService:    prodSvc,
	}
}

// ListDocs godoc
// @Summary List documentation for a product (public)
// @Tags product-docs
// @Param slug path string true "Product slug"
// @Param type query string true "Doc type (developer or user_guide)"
// @Router /products/{slug}/docs [get]
func (h *ProductDocHandler) ListDocs(c *gin.Context) {
	slug := c.Param("slug")
	docType := c.DefaultQuery("type", "developer")

	product, err := h.productService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		response.NotFound(c, "Product not found")
		return
	}

	id, err := uuid.Parse(product.ID)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", nil)
		return
	}

	grouped, err := h.productDocService.GetPublicDocs(c.Request.Context(), id, docType)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve documentation")
		return
	}

	var result []dto.ListProductDocResponse
	for sectionSlug, docs := range grouped {
		sectionTitle := ""
		if len(docs) > 0 {
			sectionTitle = docs[0].SectionTitle
		}
		section := dto.ListProductDocResponse{
			SectionSlug:  sectionSlug,
			SectionTitle: sectionTitle,
			Docs:         make([]dto.ProductDocResponse, 0),
		}
		for _, doc := range docs {
			section.Docs = append(section.Docs, dto.ProductDocResponse{
				ID:              doc.ID.String(),
				ProductID:       doc.ProductID.String(),
				DocType:         doc.DocType,
				SectionTitle:    doc.SectionTitle,
				SectionSlug:     doc.SectionSlug,
				Title:           doc.Title,
				Slug:            doc.Slug,
				Content:         doc.Content,
				SortOrder:       doc.SortOrder,
				SectionOrder:    doc.SectionOrder,
				IsPublished:     doc.IsPublished,
				MetaTitle:       doc.MetaTitle,
				MetaDescription: doc.MetaDescription,
			})
		}
		result = append(result, section)
	}

	response.OK(c, "Documentation retrieved", gin.H{
		"sections": result,
	})
}

// GetDoc godoc
// @Summary Get a single documentation page (public)
// @Tags product-docs
// @Param slug path string true "Product slug"
// @Param type query string true "Doc type (developer or user_guide)"
// @Param section path string true "Section slug"
// @Param docSlug path string true "Doc slug"
// @Router /products/{slug}/docs/{section}/{docSlug} [get]
func (h *ProductDocHandler) GetDoc(c *gin.Context) {
	slug := c.Param("slug")
	docType := c.DefaultQuery("type", "developer")
	sectionSlug := c.Param("section")
	docSlug := c.Param("docSlug")

	product, err := h.productService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		response.NotFound(c, "Product not found")
		return
	}

	id, err := uuid.Parse(product.ID)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", nil)
		return
	}

	doc, err := h.productDocService.GetDoc(c.Request.Context(), id, docType, sectionSlug, docSlug)
	if err != nil {
		response.NotFound(c, "Documentation not found")
		return
	}

	response.OK(c, "Documentation retrieved", dto.ProductDocResponse{
		ID:              doc.ID.String(),
		ProductID:       doc.ProductID.String(),
		DocType:         doc.DocType,
		SectionTitle:    doc.SectionTitle,
		SectionSlug:     doc.SectionSlug,
		Title:           doc.Title,
		Slug:            doc.Slug,
		Content:         doc.Content,
		SortOrder:       doc.SortOrder,
		SectionOrder:    doc.SectionOrder,
		IsPublished:     doc.IsPublished,
		MetaTitle:       doc.MetaTitle,
		MetaDescription: doc.MetaDescription,
	})
}

// ListAdminDocs godoc
// @Summary List product documentation (admin)
// @Tags admin
// @Security BearerAuth
// @Param productId path string true "Product ID"
// @Param type query string false "Filter by doc type (developer or user_guide)"
// @Router /admin/products/{productId}/docs [get]
func (h *ProductDocHandler) ListAdminDocs(c *gin.Context) {
	productID := c.Param("id")
	docType := c.Query("type")

	id, err := uuid.Parse(productID)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", nil)
		return
	}

	docs, err := h.productDocService.GetByProduct(c.Request.Context(), id, docType)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve documentation")
		return
	}

	var result []dto.ProductDocResponse
	for _, doc := range docs {
		result = append(result, dto.ProductDocResponse{
			ID:              doc.ID.String(),
			ProductID:       doc.ProductID.String(),
			DocType:         doc.DocType,
			SectionTitle:    doc.SectionTitle,
			SectionSlug:     doc.SectionSlug,
			Title:           doc.Title,
			Slug:            doc.Slug,
			Content:         doc.Content,
			SortOrder:       doc.SortOrder,
			SectionOrder:    doc.SectionOrder,
			IsPublished:     doc.IsPublished,
			MetaTitle:       doc.MetaTitle,
			MetaDescription: doc.MetaDescription,
		})
	}

	response.OK(c, "Documentation retrieved", gin.H{
		"docs": result,
	})
}

// CreateDoc godoc
// @Summary Create product documentation (admin)
// @Tags admin
// @Security BearerAuth
// @Param productId path string true "Product ID"
// @Router /admin/products/{productId}/docs [post]
func (h *ProductDocHandler) CreateDoc(c *gin.Context) {
	productID := c.Param("id")

	id, err := uuid.Parse(productID)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", nil)
		return
	}

	var req dto.CreateProductDocRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	doc := &domain.ProductDoc{
		ID:              uuid.New(),
		ProductID:       id,
		DocType:         req.DocType,
		SectionTitle:    req.SectionTitle,
		SectionSlug:     req.SectionSlug,
		Title:           req.Title,
		Slug:            req.Slug,
		Content:         req.Content,
		SortOrder:       req.SortOrder,
		SectionOrder:    req.SectionOrder,
		IsPublished:     req.IsPublished,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if req.IsPublished {
		now := time.Now()
		doc.PublishedAt = &now
	}

	if err := h.productDocService.Create(c.Request.Context(), doc); err != nil {
		response.InternalServerError(c, "Failed to create documentation")
		return
	}

	response.Created(c, "Documentation created", dto.ProductDocResponse{
		ID:              doc.ID.String(),
		ProductID:       doc.ProductID.String(),
		DocType:         doc.DocType,
		SectionTitle:    doc.SectionTitle,
		SectionSlug:     doc.SectionSlug,
		Title:           doc.Title,
		Slug:            doc.Slug,
		Content:         doc.Content,
		SortOrder:       doc.SortOrder,
		SectionOrder:    doc.SectionOrder,
		IsPublished:     doc.IsPublished,
		MetaTitle:       doc.MetaTitle,
		MetaDescription: doc.MetaDescription,
	})
}

// UpdateDoc godoc
// @Summary Update product documentation (admin)
// @Tags admin
// @Security BearerAuth
// @Param productId path string true "Product ID"
// @Param id path string true "Doc ID"
// @Router /admin/products/{productId}/docs/{id} [put]
func (h *ProductDocHandler) UpdateDoc(c *gin.Context) {
	docID := c.Param("docId")

	id, err := uuid.Parse(docID)
	if err != nil {
		response.BadRequest(c, "Invalid doc ID", nil)
		return
	}

	var req dto.UpdateProductDocRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	doc, err := h.productDocService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Documentation not found")
		return
	}

	// Update fields
	if req.DocType != "" {
		doc.DocType = req.DocType
	}
	if req.SectionTitle != "" {
		doc.SectionTitle = req.SectionTitle
	}
	if req.SectionSlug != "" {
		doc.SectionSlug = req.SectionSlug
	}
	if req.Title != "" {
		doc.Title = req.Title
	}
	if req.Slug != "" {
		doc.Slug = req.Slug
	}
	if req.Content != "" {
		doc.Content = req.Content
	}
	doc.SortOrder = req.SortOrder
	doc.SectionOrder = req.SectionOrder
	doc.MetaTitle = req.MetaTitle
	doc.MetaDescription = req.MetaDescription

	// Handle publish status
	wasPublished := doc.IsPublished
	doc.IsPublished = req.IsPublished
	if req.IsPublished && !wasPublished {
		now := time.Now()
		doc.PublishedAt = &now
	}

	doc.UpdatedAt = time.Now()

	if err := h.productDocService.Update(c.Request.Context(), doc); err != nil {
		response.InternalServerError(c, "Failed to update documentation")
		return
	}

	response.OK(c, "Documentation updated", dto.ProductDocResponse{
		ID:              doc.ID.String(),
		ProductID:       doc.ProductID.String(),
		DocType:         doc.DocType,
		SectionTitle:    doc.SectionTitle,
		SectionSlug:     doc.SectionSlug,
		Title:           doc.Title,
		Slug:            doc.Slug,
		Content:         doc.Content,
		SortOrder:       doc.SortOrder,
		SectionOrder:    doc.SectionOrder,
		IsPublished:     doc.IsPublished,
		MetaTitle:       doc.MetaTitle,
		MetaDescription: doc.MetaDescription,
	})
}

// DeleteDoc godoc
// @Summary Delete product documentation (admin)
// @Tags admin
// @Security BearerAuth
// @Param productId path string true "Product ID"
// @Param id path string true "Doc ID"
// @Router /admin/products/{productId}/docs/{id} [delete]
func (h *ProductDocHandler) DeleteDoc(c *gin.Context) {
	docID := c.Param("docId")

	id, err := uuid.Parse(docID)
	if err != nil {
		response.BadRequest(c, "Invalid doc ID", nil)
		return
	}

	if err := h.productDocService.Delete(c.Request.Context(), id); err != nil {
		response.NotFound(c, "Documentation not found")
		return
	}

	response.NoContent(c)
}

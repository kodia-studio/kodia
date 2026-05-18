package dto

type CreateProductDocRequest struct {
	DocType         string `json:"doc_type" binding:"required,oneof=developer user_guide"`
	SectionTitle    string `json:"section_title" binding:"required"`
	SectionSlug     string `json:"section_slug" binding:"required"`
	Title           string `json:"title" binding:"required"`
	Slug            string `json:"slug" binding:"required"`
	Content         string `json:"content" binding:"required"`
	SortOrder       int    `json:"sort_order"`
	SectionOrder    int    `json:"section_order"`
	IsPublished     bool   `json:"is_published"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
}

type UpdateProductDocRequest struct {
	DocType         string `json:"doc_type"`
	SectionTitle    string `json:"section_title"`
	SectionSlug     string `json:"section_slug"`
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	Content         string `json:"content"`
	SortOrder       int    `json:"sort_order"`
	SectionOrder    int    `json:"section_order"`
	IsPublished     bool   `json:"is_published"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
}

type ProductDocResponse struct {
	ID              string `json:"id"`
	ProductID       string `json:"product_id"`
	DocType         string `json:"doc_type"`
	SectionTitle    string `json:"section_title"`
	SectionSlug     string `json:"section_slug"`
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	Content         string `json:"content"`
	SortOrder       int    `json:"sort_order"`
	SectionOrder    int    `json:"section_order"`
	IsPublished     bool   `json:"is_published"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
}

type ListProductDocResponse struct {
	SectionSlug string                  `json:"section_slug"`
	SectionTitle string                 `json:"section_title"`
	Docs        []ProductDocResponse    `json:"docs"`
}

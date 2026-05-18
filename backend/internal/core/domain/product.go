package domain

import "time"

type ProductType string

const (
	ProductTypeSourceCode ProductType = "source_code"
	ProductTypeService    ProductType = "service"
)

type ServiceType string

const (
	ServiceTypeSupport    ServiceType = "support"
	ServiceTypeCustomApp  ServiceType = "custom_app"
)

type Product struct {
	ID              string
	Slug            string
	Name            string
	Tagline         string
	Description     string
	Type            ProductType
	ServiceType     *ServiceType
	CoverURL        string
	Tags            []string
	IsPublished     bool
	SortOrder       int
	MetaTitle       string
	MetaDescription string
	OGImageURL      string
	Variants        []ProductVariant
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

type ProductVariant struct {
	ID            string
	ProductID     string
	Name          string
	Description   string
	Price         int64
	PromoPrice    *int64
	PromoStartAt  *time.Time
	PromoEndAt    *time.Time
	FileKey       string
	IsActive      bool
	SortOrder     int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

func (v *ProductVariant) EffectivePrice() int64 {
	if v.PromoPrice != nil && v.PromoStartAt != nil && v.PromoEndAt != nil {
		now := time.Now()
		if now.After(*v.PromoStartAt) && now.Before(*v.PromoEndAt) {
			return *v.PromoPrice
		}
	}
	return v.Price
}

func (v *ProductVariant) IsPromoActive() bool {
	if v.PromoPrice == nil || v.PromoStartAt == nil || v.PromoEndAt == nil {
		return false
	}
	now := time.Now()
	return now.After(*v.PromoStartAt) && now.Before(*v.PromoEndAt)
}

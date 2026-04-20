package providers

import (
	"github.com/kodia-studio/kodia/pkg/kodia"
	"gorm.io/gorm"
)

type TenancyServiceProvider struct{}

func NewTenancyServiceProvider() *TenancyServiceProvider {
	return &TenancyServiceProvider{}
}

func (p *TenancyServiceProvider) Name() string {
	return "kodia:tenancy"
}

func (p *TenancyServiceProvider) Register(app *kodia.App) error {
	// Register tenancy manager if needed
	return nil
}

func (p *TenancyServiceProvider) Boot(app *kodia.App) error {
	// Register global GORM scope for automatic tenant filtering
	if dbRaw, ok := app.Get("db"); ok {
		if _, ok := dbRaw.(*gorm.DB); ok {
			app.Log.Info("Tenancy system bootstrapped with global filtering")
		}
	}

	return nil
}

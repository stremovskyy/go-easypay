package easypay

import "time"

type App struct {
	logoURL    *string
	apiVersion *string
	appID      *string
	validTill  time.Time
}

func NewApp(logoURL, apiVersion, appID *string) *App {
	validTill := time.Now().AddDate(0, 0, 30)

	return &App{
		logoURL:    logoURL,
		apiVersion: apiVersion,
		appID:      appID,
		validTill:  validTill,
	}
}

func (a *App) LogoURL() string {
	return *a.logoURL
}

func (a *App) IsValid() bool {
	return a.validTill.After(time.Now())
}

func (a *App) APIVersion() string {
	return *a.apiVersion
}

func (a *App) AppID() string {
	return *a.appID
}

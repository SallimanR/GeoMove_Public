package auth

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/yandex"
)

type oAuthProviders struct {
	vkConfig     *oauth2.Config
	yandexConfig *oauth2.Config
	googleConfig *oauth2.Config
}

type Config struct {
	VKClientID     string
	VKClientSecret string
	VKRedirectURL  string

	YandexClientID     string
	YandexClientSecret string
	YandexRedirectURL  string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	CookieDomain string
	CookieSecure bool
	FrontendURL  string
	StaticDir    string
}

func LoadConfig() *Config {
	return &Config{
		VKClientID:     os.Getenv("VK_CLIENT_ID"),
		VKClientSecret: os.Getenv("VK_CLIENT_SECRET"),
		VKRedirectURL:  os.Getenv("VK_REDIRECT_URL"),

		YandexClientID:     os.Getenv("YANDEX_CLIENT_ID"),
		YandexClientSecret: os.Getenv("YANDEX_CLIENT_SECRET"),
		YandexRedirectURL:  os.Getenv("YANDEX_REDIRECT_URL"),

		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),

		CookieDomain: os.Getenv("COOKIE_DOMAIN"),
		CookieSecure: os.Getenv("COOKIE_SECURE") == "true",
		FrontendURL:  os.Getenv("FRONTEND_URL"),

		StaticDir: os.Getenv("STATIC_DIR"),
	}
}

func GetOAuthProviders(cfg *Config) oAuthProviders {
	vkConfig := &oauth2.Config{
		ClientID:     cfg.VKClientID,
		ClientSecret: cfg.VKClientSecret,
		RedirectURL:  cfg.VKRedirectURL,
		Scopes:       []string{"vkid.personal_info", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://id.vk.ru/authorize",
			TokenURL:  "https://id.vk.com/oauth2/auth",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	yandexConfig := &oauth2.Config{
		ClientID:     cfg.YandexClientID,
		ClientSecret: cfg.YandexClientSecret,
		RedirectURL:  cfg.YandexRedirectURL,
		Scopes:       []string{"login:email", "login:info"},
		Endpoint:     yandex.Endpoint,
	}

	googleConfig := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	return oAuthProviders{
		vkConfig:     vkConfig,
		yandexConfig: yandexConfig,
		googleConfig: googleConfig,
	}
}

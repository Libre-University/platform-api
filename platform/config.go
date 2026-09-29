package platform

import "os"

// Config uygulama yapılandırmasıdır. Ortam değişkenleri LU_ önekiyle okunur.
type Config struct {
	// HTTPAddr HTTP dinleme adresi (LU_HTTP_ADDR, varsayılan ":8080").
	HTTPAddr string
	// DatabaseURL PostgreSQL bağlantı adresi (LU_DATABASE_URL). Boşsa
	// uygulama veritabanısız çalışır; /ready bunu bildirir.
	DatabaseURL string
	// Version derleme sırasında verilen uygulama sürümüdür.
	Version string
}

// ConfigFromEnv ortam değişkenlerinden Config üretir.
func ConfigFromEnv() Config {
	cfg := Config{
		HTTPAddr:    os.Getenv("LU_HTTP_ADDR"),
		DatabaseURL: os.Getenv("LU_DATABASE_URL"),
		Version:     "dev",
	}
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = ":8080"
	}
	return cfg
}

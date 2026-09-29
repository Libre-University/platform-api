package platform

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Libre-University/platform-api/core"
	"github.com/Libre-University/platform-api/platform/internal/audit"
	"github.com/Libre-University/platform-api/platform/internal/db"
	"github.com/Libre-University/platform-api/platform/internal/httpx"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// schemaName çekirdeğin PostgreSQL şemasıdır.
const schemaName = "platform"

// App çalışan uygulamadır.
type App struct {
	cfg     Config
	log     *slog.Logger
	modules []core.Module
	router  chi.Router
	db      *db.DB
}

// New modülleri başlatır ve HTTP yönlendiricisini kurar. DatabaseURL doluysa
// bağlantı havuzu açılır; migration'lar ayrıca Migrate ile çalıştırılır.
func New(ctx context.Context, cfg Config, log *slog.Logger, modules []core.Module) (*App, error) {
	if log == nil {
		log = slog.Default()
	}
	seen := map[string]bool{}
	for _, m := range modules {
		if !core.ValidName(m.Name()) {
			return nil, fmt.Errorf("platform: geçersiz modül adı %q", m.Name())
		}
		if seen[m.Name()] {
			return nil, fmt.Errorf("platform: modül %q birden fazla kez verildi", m.Name())
		}
		seen[m.Name()] = true
	}

	a := &App{cfg: cfg, log: log, modules: modules}

	if cfg.DatabaseURL != "" {
		d, err := db.Open(ctx, cfg.DatabaseURL)
		if err != nil {
			return nil, fmt.Errorf("platform: veritabanı: %w", err)
		}
		a.db = d
	}

	svc := &core.Services{
		Logger: log,
		Audit:  audit.NewLogAuditor(log),
	}
	for _, m := range modules {
		if err := m.Init(ctx, svc); err != nil {
			a.Close()
			return nil, fmt.Errorf("platform: modül %s başlatılamadı: %w", m.Name(), err)
		}
		log.Info("modül başlatıldı", "module", m.Name(), "version", m.Version())
	}

	a.router = a.buildRouter()
	return a, nil
}

func (a *App) buildRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	// Gerçek istemci IP'si (audit için) Faz 1'de ters vekil sunucuya güvenilen
	// başlıkla ele alınır; chi'nin RealIP ara katmanı sahteciliğe açık olduğu
	// için kullanılmaz.
	r.Use(httpx.Logger(a.log))
	r.Use(middleware.Recoverer)

	var readiness httpx.ReadinessFunc
	if a.db != nil {
		readiness = a.db.Ping
	}
	r.Get("/health", httpx.Health(a.cfg.Version))
	r.Get("/ready", httpx.Ready(readiness))
	r.Get("/api/v1/modules", httpx.ModuleList(a.moduleInfos()))

	for _, m := range a.modules {
		r.Route("/api/v1/"+m.Name(), m.Routes)
	}
	return r
}

func (a *App) moduleInfos() []httpx.ModuleInfo {
	out := make([]httpx.ModuleInfo, 0, len(a.modules))
	for _, m := range a.modules {
		perms := make([]string, 0, len(m.Permissions()))
		for _, p := range m.Permissions() {
			perms = append(perms, p.Code)
		}
		out = append(out, httpx.ModuleInfo{Name: m.Name(), Version: m.Version(), Permissions: perms})
	}
	return out
}

// Handler HTTP işleyicisini döndürür.
func (a *App) Handler() http.Handler { return a.router }

// Modules başlatılmış modülleri döndürür.
func (a *App) Modules() []core.Module { return a.modules }

// Migrate çekirdek ve modül migration'larını çalıştırır. Her modül kendi
// şemasında ve kendi goose sürüm tablosuyla izlenir.
func (a *App) Migrate(ctx context.Context) error {
	if a.db == nil {
		return fmt.Errorf("platform: migration için LU_DATABASE_URL gerekli")
	}
	coreFS, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		return err
	}
	if err := a.db.Migrate(ctx, schemaName, coreFS); err != nil {
		return fmt.Errorf("platform: çekirdek migration: %w", err)
	}
	a.log.Info("migration tamamlandı", "schema", schemaName)
	for _, m := range a.modules {
		fsys := m.Migrations()
		if fsys == nil {
			continue
		}
		if err := a.db.Migrate(ctx, m.Name(), fsys); err != nil {
			return fmt.Errorf("platform: modül %s migration: %w", m.Name(), err)
		}
		a.log.Info("migration tamamlandı", "schema", m.Name())
	}
	return nil
}

// Close kaynakları serbest bırakır.
func (a *App) Close() {
	if a.db != nil {
		a.db.Close()
	}
}

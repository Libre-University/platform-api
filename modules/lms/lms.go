// Package lms, öğrenme yönetim sistemi modülüdür. ADR-0013 kabul edilirse bu
// modül LMS yazmaz; mevcut özgür bir LMS'e (ilk hedef Moodle) entegrasyonu
// yürütür: şube-ders eşleştirmesi, OBS→LMS katılımcı senkronizasyonu ve
// LMS→OBS not aktarımı. Kabul edilmezse ders sayfası, materyal ve ödev
// varlıkları burada geliştirilir (MODULES.md §3, FR-03).
package lms

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Libre-University/platform-api/core"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Version modül sürümüdür.
const Version = "0.0.1"

// Module LMS modülüdür.
type Module struct {
	log *slog.Logger
}

func init() { core.Register(&Module{}) }

// Name modül adıdır; /api/v1/lms ve lms şeması.
func (m *Module) Name() string { return "lms" }

// Version modül sürümüdür.
func (m *Module) Version() string { return Version }

// Migrations lms şemasının goose migration'larını döndürür.
func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		panic(err)
	}
	return sub
}

// Permissions LMS izinleridir.
func (m *Module) Permissions() []core.Permission {
	return []core.Permission{
		{Code: "lms.course_link.read", Description: "Şube-LMS ders eşleştirmesini görüntüle"},
		{Code: "lms.course_link.manage", Description: "Şubeyi LMS dersiyle eşleştir ve senkronizasyonu yönet"},
		{Code: "lms.grade_import.run", Description: "LMS'ten not aktarımını başlat"},
	}
}

// Init hizmetleri saklar.
func (m *Module) Init(_ context.Context, s *core.Services) error {
	m.log = s.Logger.With("module", m.Name())
	return nil
}

// Routes modül uçlarını bağlar. Faz 0'da yalnızca bilgi ucu vardır.
func (m *Module) Routes(r chi.Router) {
	r.Get("/", m.info)
}

func (m *Module) info(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write([]byte(`{"module":"lms","version":"` + Version + `","phase":0}` + "\n"))
}

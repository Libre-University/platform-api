// Package obs, Öğrenci Bilgi Sistemi modülüdür: öğrencilik kaydı, müfredat,
// ders açma, dönemlik kayıt ve danışman onayı, not girişi ve ders sonucu
// (MODULES.md §2, FR-02, DATA_MODEL.md).
//
// Bu paket çekirdeğe yalnızca core üzerinden erişir. İç kod internal/ altında
// tutulur; başka modüller erişemez (ADR-0012).
package obs

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

// Module OBS modülüdür.
type Module struct {
	log *slog.Logger
}

func init() { core.Register(&Module{}) }

// Name modül adıdır; /api/v1/obs ve obs şeması.
func (m *Module) Name() string { return "obs" }

// Version modül sürümüdür.
func (m *Module) Version() string { return Version }

// Migrations obs şemasının goose migration'larını döndürür.
func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		panic(err)
	}
	return sub
}

// Permissions OBS izinleridir (DATA_MODEL.md › Temel İş Kuralları).
func (m *Module) Permissions() []core.Permission {
	return []core.Permission{
		{Code: "obs.student_record.read", Description: "Öğrencilik kaydını görüntüle"},
		{Code: "obs.student_record.manage", Description: "Öğrencilik kaydı oluştur ve güncelle"},
		{Code: "obs.curriculum.manage", Description: "Müfredat ve müfredat derslerini yönet"},
		{Code: "obs.section.read", Description: "Ders şubelerini görüntüle"},
		{Code: "obs.section.manage", Description: "Ders şubesi aç, kapat ve düzenle"},
		{Code: "obs.registration.submit", Description: "Dönemlik ders kaydını hazırla ve gönder"},
		{Code: "obs.registration.approve", Description: "Danışan öğrencinin dönemlik kaydını onayla veya geri gönder"},
		{Code: "obs.grade.enter", Description: "Sorumlu olunan şubede not gir"},
		{Code: "obs.grade.override", Description: "Dönem kapandıktan sonra gerekçeli not değişikliği yap"},
		{Code: "obs.grading_policy.manage", Description: "Değerlendirme politikası tanımla ve sürümle"},
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
	_, _ = w.Write([]byte(`{"module":"obs","version":"` + Version + `","phase":0}` + "\n"))
}

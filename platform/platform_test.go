package platform

import (
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/Libre-University/platform-api/core"
)

type fakeModule struct {
	name  string
	inits int
	svc   *core.Services
}

func (f *fakeModule) Name() string      { return f.name }
func (f *fakeModule) Version() string   { return "0.1.0" }
func (f *fakeModule) Migrations() fs.FS { return nil }
func (f *fakeModule) Permissions() []core.Permission {
	return []core.Permission{{Code: f.name + ".thing.read"}}
}
func (f *fakeModule) Init(_ context.Context, s *core.Services) error {
	f.inits++
	f.svc = s
	return nil
}
func (f *fakeModule) Routes(r chi.Router) {
	r.Get("/hello", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("hello from " + f.name))
	})
}

func newTestApp(t *testing.T, mods ...core.Module) *App {
	t.Helper()
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	app, err := New(context.Background(), Config{Version: "test"}, log, mods)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(app.Close)
	return app
}

func get(t *testing.T, h http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestHealthAndReadyWithoutDatabase(t *testing.T) {
	app := newTestApp(t)
	if rec := get(t, app.Handler(), "/health"); rec.Code != http.StatusOK {
		t.Fatalf("/health %d", rec.Code)
	}
	rec := get(t, app.Handler(), "/ready")
	if rec.Code != http.StatusOK {
		t.Fatalf("/ready %d", rec.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["database"] != "not_configured" {
		t.Fatalf("database=%q", body["database"])
	}
}

func TestModulesAreInitialisedAndMounted(t *testing.T) {
	m := &fakeModule{name: "demo"}
	app := newTestApp(t, m)

	if m.inits != 1 || m.svc == nil || m.svc.Audit == nil || m.svc.Logger == nil {
		t.Fatalf("modül düzgün başlatılmadı: inits=%d svc=%+v", m.inits, m.svc)
	}
	rec := get(t, app.Handler(), "/api/v1/demo/hello")
	if rec.Code != http.StatusOK || rec.Body.String() != "hello from demo" {
		t.Fatalf("rota bağlanmadı: %d %q", rec.Code, rec.Body.String())
	}

	rec = get(t, app.Handler(), "/api/v1/modules")
	var body struct {
		Modules []struct {
			Name        string   `json:"name"`
			Permissions []string `json:"permissions"`
		} `json:"modules"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body.Modules) != 1 || body.Modules[0].Name != "demo" || body.Modules[0].Permissions[0] != "demo.thing.read" {
		t.Fatalf("modül listesi: %+v", body.Modules)
	}
}

func TestNewRejectsDuplicateModuleNames(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	_, err := New(context.Background(), Config{}, log, []core.Module{&fakeModule{name: "x"}, &fakeModule{name: "x"}})
	if err == nil {
		t.Fatal("tekrar eden modül adı hata vermeliydi")
	}
}

func TestMigrateRequiresDatabase(t *testing.T) {
	app := newTestApp(t)
	if err := app.Migrate(context.Background()); err == nil {
		t.Fatal("veritabanısız Migrate hata vermeliydi")
	}
}

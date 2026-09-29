package obs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Libre-University/platform-api/core"
	"github.com/Libre-University/platform-api/core/coretest"
)

func TestModuleContract(t *testing.T) {
	m := &Module{}
	if !core.ValidName(m.Name()) {
		t.Fatalf("geçersiz modül adı %q", m.Name())
	}
	if m.Migrations() == nil {
		t.Fatal("migration dosya sistemi nil")
	}
	for _, p := range m.Permissions() {
		if !strings.HasPrefix(p.Code, m.Name()+".") {
			t.Errorf("izin kodu modül adıyla başlamalı: %q", p.Code)
		}
	}
}

func TestInfoEndpoint(t *testing.T) {
	h, _ := coretest.Mount(t, &Module{})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/obs/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("durum %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"module":"obs"`) {
		t.Fatalf("gövde: %s", rec.Body.String())
	}
}

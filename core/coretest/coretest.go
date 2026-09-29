// Package coretest, modülleri gerçek çekirdek olmadan test etmek için sahte
// hizmetler ve yardımcılar sağlar (ADR-0012). Modül testleri platform paketine
// bağımlı olmamalıdır.
package coretest

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/Libre-University/platform-api/core"
)

// AuditRecorder, olayları bellekte tutan sahte Auditor'dır.
type AuditRecorder struct {
	mu     sync.Mutex
	Events []core.AuditEvent
}

// Record olayı listeye ekler.
func (a *AuditRecorder) Record(_ context.Context, e core.AuditEvent) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Events = append(a.Events, e)
	return nil
}

// Len kaydedilen olay sayısını döndürür.
func (a *AuditRecorder) Len() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.Events)
}

// NewServices, sessiz bir logger ve bellek içi audit kaydedici ile Services
// üretir.
func NewServices() (*core.Services, *AuditRecorder) {
	rec := &AuditRecorder{}
	return &core.Services{
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Audit:  rec,
	}, rec
}

// Mount modülü başlatır ve platformun yaptığı gibi /api/v1/<name> altına
// bağlanmış bir HTTP işleyici döndürür.
func Mount(tb testing.TB, m core.Module) (http.Handler, *AuditRecorder) {
	tb.Helper()
	svc, rec := NewServices()
	if err := m.Init(context.Background(), svc); err != nil {
		tb.Fatalf("modül %s başlatılamadı: %v", m.Name(), err)
	}
	r := chi.NewRouter()
	r.Route("/api/v1/"+m.Name(), m.Routes)
	return r, rec
}

package core

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"regexp"
	"sort"
	"sync"

	"github.com/go-chi/chi/v5"
)

// Permission, bir modülün tanımladığı izindir. Kod, "<modül>.<varlık>.<eylem>"
// biçimindedir (ör. obs.registration.approve).
type Permission struct {
	Code        string
	Description string
}

// Services, çekirdeğin modüllere sağladığı hizmetlerdir. Faz 0'da yalnızca
// log ve audit vardır; yetki, bildirim, dosya ve akademik sorgular Faz 1'de
// eklenir (core v0.1).
type Services struct {
	Logger *slog.Logger
	Audit  Auditor
}

// Module, platforma derleme zamanında eklenen bir iş modülüdür.
type Module interface {
	// Name modülün kısa adıdır; URL öneki (/api/v1/<name>) ve PostgreSQL
	// şema adı olarak kullanılır. ValidName kuralına uymalıdır.
	Name() string
	// Version modülün sürümüdür (semver).
	Version() string
	// Migrations goose biçimindeki SQL migration dosyalarını içeren dosya
	// sistemidir. Migration'ı olmayan modül nil döndürür.
	Migrations() fs.FS
	// Permissions modülün tanımladığı izinlerdir.
	Permissions() []Permission
	// Init, HTTP sunucusu başlamadan önce bir kez çağrılır.
	Init(ctx context.Context, s *Services) error
	// Routes modülün HTTP uçlarını verilen yönlendiriciye bağlar. Yönlendirici
	// zaten /api/v1/<name> altına bağlanmıştır.
	Routes(r chi.Router)
}

var nameRe = regexp.MustCompile(`^[a-z][a-z0-9_]{1,31}$`)

// ValidName, modül adının URL ve şema adı olarak güvenli olup olmadığını
// söyler: küçük harfle başlar, en fazla 32 karakter, yalnızca [a-z0-9_].
func ValidName(name string) bool {
	return nameRe.MatchString(name)
}

var registry = struct {
	mu   sync.Mutex
	mods map[string]Module
}{mods: map[string]Module{}}

// Register bir modülü kaydeder. Genellikle modül paketinin init() işlevinden
// çağrılır. Geçersiz veya tekrar eden ad programlama hatasıdır ve panic
// üretir.
func Register(m Module) {
	if m == nil {
		panic("core: nil modül kaydedilemez")
	}
	name := m.Name()
	if !ValidName(name) {
		panic(fmt.Sprintf("core: geçersiz modül adı %q", name))
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	if _, dup := registry.mods[name]; dup {
		panic(fmt.Sprintf("core: modül %q zaten kayıtlı", name))
	}
	registry.mods[name] = m
}

// Registered kayıtlı modülleri ada göre sıralı döndürür.
func Registered() []Module {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	out := make([]Module, 0, len(registry.mods))
	for _, m := range registry.mods {
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name() < out[j].Name() })
	return out
}

// resetRegistry yalnızca testler içindir.
func resetRegistry() {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.mods = map[string]Module{}
}

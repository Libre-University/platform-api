# platform-api

LibreUniversity'nin **çekirdek ve dağıtım** backend uygulaması. Go ile yazılır ve tek bir binary olarak dağıtılır ([ADR-0004](https://github.com/Libre-University/docs/blob/main/docs/adr/0004-start-with-modular-monolith-for-mvp.md), [ADR-0011](https://github.com/Libre-University/docs/blob/main/docs/adr/0011-use-go-for-backend.md)). İş modülleri MVP boyunca bu repoda `modules/` altında, her biri kendi `go.mod` dosyasına sahip ayrı Go modülü olarak geliştirilir ve derleme sırasında tek uygulamada birleşir ([ADR-0012](https://github.com/Libre-University/docs/blob/main/docs/adr/0012-develop-modules-inside-platform-api.md)).

## Ne İş Yapar?

- Tüm modüllerin ortak kullandığı çekirdek servisleri sağlar: kimlik ve yetki, akademik çekirdek veri, audit, bildirim, dosya (Faz 1).
- Modüllerin kullandığı genel arayüz paketini (`core`) yayımlar; çekirdeğin geri kalanı `platform/internal/` altında derleyici tarafından kapalıdır.
- MVP modülleriyle (`modules/obs`, `modules/lms`) varsayılan dağıtımı derler (`cmd/libre-university`). Farklı modül seti isteyen kurumlar bu dizini kopyalayıp içe aktarılan modülleri değiştirir.
- Modülleri çekirdek olmadan test etmek için `core/coretest` paketini sağlar.
- Tek REST API (OpenAPI 3, Faz 0 sonunda) sunar; modüllerin uçları `/api/v1/<modül>` altında toplanır.

## Dizin Yapısı

```
go.work                    # çalışma alanı: aşağıdaki beş Go modülü
core/                      # modüllerin gördüğü genel arayüz: Module, Register, Services, Auditor
core/coretest/             # sahte hizmetler ve Mount yardımcısı (modül testleri için)
platform/                  # çekirdek uygulama; dışarıya Config, New, App açık
platform/internal/httpx/   # /health, /ready, /api/v1/modules, istek logu
platform/internal/db/      # pgx havuzu, şema bazlı goose migration
platform/internal/audit/   # Faz 0: log tabanlı Auditor
platform/migrations/       # platform şeması migration'ları
modules/obs/               # Öğrenci Bilgi Sistemi modülü (obs şeması, /api/v1/obs)
modules/lms/               # LMS modülü (lms şeması, /api/v1/lms); kapsamı ADR-0013'e bağlı
cmd/libre-university/      # varsayılan dağıtımın main paketi: serve, migrate, modules, version
```

Her modül `core.Module` arayüzünü uygular ve `init()` içinde `core.Register` ile kendini kaydeder. Modül sınırları üç katmanla korunur: Go `internal` kuralı (`platform/internal`, `modules/*/internal`), `golangci-lint` `depguard` (modüller `platform` paketini içe aktaramaz, çekirdek modülleri tanımaz) ve şema bazlı migration (her modülün kendi PostgreSQL şeması ve kendi `goose_db_version` tablosu).

## Teknoloji

- Go 1.26+ (desteklenen son iki sürüm), `net/http` + chi
- PostgreSQL 16+, pgx; goose migration'ları (modül başına şema ve sürüm tablosu)
- `log/slog` ile JSON log; test için `testing` + `httptest`
- Kalite: gofmt, go vet, golangci-lint (depguard dahil), `go test -race`
- Faz 0 sonunda: OpenAPI + oapi-codegen, sqlc, River, go-oidc ([ADR-0011](https://github.com/Libre-University/docs/blob/main/docs/adr/0011-use-go-for-backend.md))

## Hızlı Başlangıç

```sh
make build                       # bin/libre-university
./bin/libre-university modules   # derlenmiş modüller ve izin sayıları
./bin/libre-university serve     # :8080; LU_DATABASE_URL yoksa veritabanısız çalışır
curl localhost:8080/health
curl localhost:8080/api/v1/modules
```

PostgreSQL ile:

```sh
export LU_DATABASE_URL='postgres://libre:libre@localhost:5432/libre?sslmode=disable'
make migrate                     # platform, obs, lms şemaları ve migration'ları
make run
curl localhost:8080/ready        # {"database":"ok","status":"ok"}
```

Go yerelde yoksa Docker ile: `make GO="docker run --rm -v $PWD:/src -w /src golang:1.26 go" test`.

Geliştirme adımları, test ve lint için [CONTRIBUTING.md](CONTRIBUTING.md).

## Yeni Modül Eklemek

1. `modules/<ad>/` dizini ve `go.mod` (`github.com/Libre-University/platform-api/modules/<ad>`, `replace` ile `../../core`).
2. `core.Module` uygulaması; `init()` içinde `core.Register`.
3. `migrations/00001_init.sql`; şemayı platform oluşturur, modül tablolarını yazar.
4. `go.work` dosyasına `use ./modules/<ad>`; `cmd/libre-university/main.go` içine blank import ve `replace`.
5. `.github/CODEOWNERS` içine `/modules/<ad>/` satırı.
6. Testler `core/coretest` ile çekirdek olmadan yazılır.

Modül adı `^[a-z][a-z0-9_]{1,31}$` desenine uymalıdır; URL öneki ve şema adı olarak kullanılır.

## Fazlara Göre İşler

| Faz | Bu repoda yapılacaklar |
| --- | --- |
| Faz 0 | Go iskeleti, çekirdek paketler, modül kayıt mekanizması, `core` arayüz taslağı, `coretest`, CI, sağlık uçları, OpenAPI, yapay örnek veri |
| Faz 1 | OIDC girişi, rol/izin ve birim kapsamlı yetki, akademik çekirdek veri, audit log tablosu, e-posta bildirimi, dosya depolama, `core` v0.1, metrikler |
| Faz 2 | OBS modülünün ihtiyaç duyduğu çekirdek arayüz genişletmeleri, olay (domain event) altyapısı, ders kayıt yoğunluğu için kuyruk ve performans |
| Faz 3 | Bildirim merkezi ve push kanalı, mobil uygulama API ihtiyaçları |
| Faz 4+ | KVKK açık rıza, saklama politikası ve anonimleştirme; raporlama dışa aktarımı |

Ayrıntılı ve işaretlenebilir liste: [ROADMAP.md](ROADMAP.md). Fazlar [ana yol haritası](https://github.com/Libre-University/docs/blob/main/ROADMAP.md) ile hizalıdır. Açık işler için `phase:*` etiketlerine bakın.

## Katkı

Katkı rehberi, davranış kuralları ve güvenlik politikası organizasyon genelinde [`.github`](https://github.com/Libre-University/.github) reposundadır. Mimari kararlar [`docs`](https://github.com/Libre-University/docs) reposundaki ADR'lerle alınır.

## Lisans

Bu proje [GNU Affero Genel Kamu Lisansı v3.0 veya sonrası](LICENSE) (AGPL-3.0-or-later) ile lisanslanmıştır. Ağ üzerinden hizmet olarak sunulan değiştirilmiş sürümlerin kaynak kodu da kullanıcılarla paylaşılmalıdır ([ADR-0002](https://github.com/Libre-University/docs/blob/main/docs/adr/0002-prefer-agpl-3-or-later-license.md)).

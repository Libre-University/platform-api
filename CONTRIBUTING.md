# platform-api Katkı Rehberi

Genel katkı ilkeleri, davranış kuralları ve güvenlik politikası organizasyon genelinde [`.github`](https://github.com/Libre-University/.github/blob/main/CONTRIBUTING.md) reposundadır. Bu belge yalnızca bu repoya özgü geliştirme adımlarını anlatır.

## Gereksinimler

- Go 1.26 veya üstü (desteklenen son iki sürüm). Yerelde Go yoksa tüm `make` hedefleri Docker ile çalıştırılabilir: `make GO="docker run --rm -v $PWD:/src -w /src golang:1.26 go" test`.
- PostgreSQL 16+ (yalnızca `migrate`, `/ready` ve ileride veritabanı testleri için). Hızlı yol: `docker run --rm -e POSTGRES_USER=libre -e POSTGRES_PASSWORD=libre -e POSTGRES_DB=libre -p 5432:5432 postgres:17`.
- golangci-lint v2 (`make lint` için).

## Günlük Komutlar

| Komut | Ne yapar |
| --- | --- |
| `make build` | `bin/libre-university` üretir |
| `make test` | Çalışma alanındaki tüm modüllerde `go test -race` |
| `make vet` / `make fmt` / `make lint` | CI'daki kalite kontrolleri |
| `make tidy` | Her modülde `go mod tidy`; CI temiz ağaç bekler |
| `make run` | Sunucuyu `:8080` üzerinde başlatır |
| `make migrate` | Migration'ları uygular (`LU_DATABASE_URL` gerekli) |

Bu repo bir Go çalışma alanıdır (`go.work`); kökte `go.mod` yoktur. Bu yüzden `go test ./...` kökte çalışmaz; `make test` veya `go test github.com/Libre-University/platform-api/...` kullanın.

## Modül Sınırları

[ADR-0012](https://github.com/Libre-University/docs/blob/main/docs/adr/0012-develop-modules-inside-platform-api.md) gereği:

- Modüller (`modules/*`) yalnızca `core` paketini içe aktarır. `platform` paketini içe aktaran PR lint'te reddedilir (`depguard`).
- Çekirdek (`core`, `platform`) hiçbir modülü tanımaz; modüller `cmd/` içinde bir araya gelir.
- Her modülün tabloları kendi şemasındadır; başka modülün tablosuna doğrudan sorgu yazılmaz. Veri `core` arayüzü veya olaylar üzerinden alınır.
- Modül testleri `core/coretest` ile, çekirdek ve veritabanı olmadan yazılır.

## Migration Kuralları

- goose biçimi: `-- +goose Up` / `-- +goose Down`. Dosya adı `NNNNN_aciklama.sql`.
- Şemayı platform oluşturur; modül migration'ı yalnızca kendi şemasındaki nesneleri değiştirir (`obs.student_record` gibi şema nitelikli adlar).
- Uygulanmış bir migration değiştirilmez; düzeltme yeni migration ile yapılır.
- CI, migration'ları boş bir PostgreSQL üzerinde iki kez çalıştırır; ikinci çalıştırma boş geçmelidir.

## Commit ve PR

- Commit'ler `git commit -s` ile imzalanır (DCO, [ADR-0014](https://github.com/Libre-University/docs/blob/main/docs/adr/0014-use-dco-for-contributions.md)).
- Branch adı: `feat/…`, `fix/…`, `docs/…`, `chore/…`.
- PR açıklaması amaç, kapsam, test ve KVKK/güvenlik etkisi bölümlerini içerir (şablon `.github` reposunda).
- Yeni bağımlılık eklerken lisansı ve gerekçesi PR'da yazılır ([OPEN_SOURCE_POLICY.md](https://github.com/Libre-University/docs/blob/main/OPEN_SOURCE_POLICY.md)).

## Go'ya Yeni Başlayanlar İçin

- [A Tour of Go](https://go.dev/tour/) ve [Effective Go](https://go.dev/doc/effective_go).
- Çalışma alanları: [Go workspaces](https://go.dev/doc/tutorial/workspaces).
- HTTP yönlendirici: [chi](https://github.com/go-chi/chi). Migration: [goose](https://github.com/pressly/goose).
- İlk katkı için `modules/obs/obs_test.go` ve `core/coretest` iyi bir başlangıç noktasıdır.

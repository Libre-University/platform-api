# platform-api

LibreUniversity'nin **çekirdek ve dağıtım** backend uygulaması. Go ile yazılır ve tek bir binary olarak dağıtılır ([ADR-0004](https://github.com/Libre-University/docs/blob/main/docs/adr/0004-start-with-modular-monolith-for-mvp.md), [ADR-0011](https://github.com/Libre-University/docs/blob/main/docs/adr/0011-use-go-for-backend.md)). İş modülleri (`module-obs`, `module-lms`, …) ayrı repolarda Go modülü olarak geliştirilir ve derleme sırasında bu uygulamaya eklenir ([ADR-0010](https://github.com/Libre-University/docs/blob/main/docs/adr/0010-develop-modules-as-separate-packages.md)).

## Ne İş Yapar?

- Tüm modüllerin ortak kullandığı çekirdek servisleri sağlar: kimlik ve yetki, akademik çekirdek veri, audit, bildirim, dosya.
- Modüllerin kullandığı genel arayüz paketini (`core`) yayımlar; çekirdeğin geri kalanı `internal/` altında kapalıdır.
- MVP modülleriyle (`module-obs`, `module-lms`) varsayılan dağıtımı derler; farklı modül setleri için `libre-build` aracını sağlar.
- Yeni modül repoları için modül şablonu (`gonew`) ve modül test paketi (`core/coretest`) sağlar.
- Tek REST API (OpenAPI 3) sunar; modüllerin uçları da bu API altında toplanır.

## Teknoloji

- Go (desteklenen son iki sürüm), `net/http` + chi
- OpenAPI önce yazılır; sunucu kodu oapi-codegen ile üretilir
- PostgreSQL 16+, pgx + sqlc, goose migration'ları
- River ile arka plan işleri ve kuyruk (Redis gerekmez)
- OIDC ile kimlik (Keycloak, go-oidc); parola uygulamada tutulmaz
- `log/slog`, Prometheus; test için `testing` + testcontainers-go; kalite için golangci-lint

## Dizin Yapısı (hedef)

```
cmd/libre-university/   # varsayılan dağıtımın main paketi
cmd/libre-build/        # özel modül setiyle derleme aracı
core/                   # modüllerin kullandığı genel arayüz (Module, Register, yetki, audit, bildirim, dosya)
core/coretest/          # modül test paketi
internal/identity/      # kullanıcı, rol, izin, birim kapsamlı yetki, OIDC oturumu
internal/academic/      # organizasyon birimi, program, dönem, ders kataloğu, akademisyen
internal/audit/         # değiştirilemez denetim kaydı
internal/notifications/ # bildirim kayıtları ve kanal gönderimi
internal/files/         # dosya metadata, S3 uyumlu depolama
internal/modules/       # modül kaydı, migration ve rota birleştirme, sürüm uyumluluğu
api/openapi.yaml        # çekirdek API sözleşmesi
template/module/        # yeni modül şablonu
```

## Fazlara Göre İşler

| Faz | Bu repoda yapılacaklar |
| --- | --- |
| Faz 0 | Go iskeleti, çekirdek paketler, modül kayıt mekanizması, `core` arayüz taslağı, modül şablonu ve test paketi, CI, sağlık uçları, OpenAPI, yapay örnek veri |
| Faz 1 | OIDC girişi, rol/izin ve birim kapsamlı yetki, akademik çekirdek veri, audit log, e-posta bildirimi, dosya depolama, `core` v0.1, metrikler |
| Faz 2 | OBS modülünün ihtiyaç duyduğu çekirdek arayüz genişletmeleri, olay (domain event) altyapısı, ders kayıt yoğunluğu için kuyruk ve performans |
| Faz 3 | Bildirim merkezi ve push kanalı, mobil uygulama API ihtiyaçları |
| Faz 4+ | KVKK açık rıza, saklama politikası ve anonimleştirme; raporlama dışa aktarımı |

Ayrıntılı ve işaretlenebilir liste: [ROADMAP.md](ROADMAP.md). Fazlar [ana yol haritası](https://github.com/Libre-University/docs/blob/main/ROADMAP.md) ile hizalıdır. Açık işler için `phase:*` etiketlerine bakın.

## Katkı

Katkı rehberi, davranış kuralları ve güvenlik politikası organizasyon genelinde [`.github`](https://github.com/Libre-University/.github) reposundadır. Mimari kararlar [`docs`](https://github.com/Libre-University/docs) reposundaki ADR'lerle alınır.

## Lisans

Bu proje [GNU Affero Genel Kamu Lisansı v3.0 veya sonrası](LICENSE) (AGPL-3.0-or-later) ile lisanslanmıştır. Ağ üzerinden hizmet olarak sunulan değiştirilmiş sürümlerin kaynak kodu da kullanıcılarla paylaşılmalıdır ([ADR-0002](https://github.com/Libre-University/docs/blob/main/docs/adr/0002-prefer-agpl-3-or-later-license.md)).

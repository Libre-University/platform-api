# platform-api Yol Haritası

Fazlar [ana yol haritası](https://github.com/Libre-University/docs/blob/main/ROADMAP.md) ile hizalıdır. Veri modeli: [DATA_MODEL.md](https://github.com/Libre-University/docs/blob/main/DATA_MODEL.md). Teknoloji: [ADR-0011](https://github.com/Libre-University/docs/blob/main/docs/adr/0011-use-go-for-backend.md). Modül yapısı: [ADR-0010](https://github.com/Libre-University/docs/blob/main/docs/adr/0010-develop-modules-as-separate-packages.md).

## Faz 0: İskelet (davet öncesi)

- [ ] Go modülü (`github.com/Libre-University/platform-api`), dizin yapısı, `cmd/libre-university`
- [ ] Çekirdek paket iskeletleri: `internal/identity`, `internal/academic`, `internal/audit`, `internal/notifications`, `internal/files`, `internal/modules`
- [ ] Modül kaydı: `core.Module` arayüzü, `core.Register`, migration/rota/iş birleştirme, sürüm uyumluluk kontrolü
- [ ] `core` genel arayüz paketinin ilk taslağı
- [ ] Modül şablonu (`gonew`) ve modül test paketi (`core/coretest`)
- [ ] `golangci-lint` (depguard dahil) ile modül sınırı kuralları
- [ ] CI: gofmt, go vet, golangci-lint, `go test -race`, sqlc ve oapi-codegen üretim kontrolü, bağımlılık lisans taraması
- [ ] `/health` ve `/ready` uçları
- [ ] `api/openapi.yaml` ve CI'da birleşik şemanın artefakt olarak yayımlanması
- [ ] Yapay (anonim) örnek veri üreten `seed` alt komutu
- [ ] `CONTRIBUTING.md` içinde yerel geliştirme adımları ve Go'ya giriş kaynakları

## Faz 1: Çekirdek Platform

- [ ] OIDC ile giriş/çıkış (Keycloak, `adapters` içindeki OIDC adapteri üzerinden)
- [ ] Kullanıcı, rol, izin, kullanıcı-rol ve birim kapsamlı yetki kontrolü
- [ ] Organizasyon birimi hiyerarşisi, akademik dönem, program, ders, akademisyen
- [ ] Audit log: kritik işlemler, giriş/çıkış, yetki değişiklikleri; güncelleme/silme veritabanı düzeyinde engelli
- [ ] Bildirim modeli ve e-posta kanalı (SMTP adapteri, River ile gönderim)
- [ ] Dosya metadata ve S3/MinIO adapteri
- [ ] `core` v0.1: yetki, audit, bildirim, dosya, akademik sorgu arayüzleri
- [ ] Prometheus metrikleri ve yapılandırılmış (JSON) loglar

## Faz 2: Akademik İş Akışlarına Destek

- [ ] Olay (domain event) altyapısı: modüller arası gevşek bağlı iletişim
- [ ] `module-obs` ihtiyaçlarına göre `core` arayüz genişletmeleri
- [ ] Ders kayıt yoğunluğu için River kuyruğu, bağlantı havuzu ayarları ve yük testi

## Faz 3: Self-Servis Desteği

- [ ] Bildirim merkezi API'si, push kanalı
- [ ] Mobil uygulama için API ihtiyaçları

## Faz 4+: Kurumsal Süreçler

- [ ] KVKK: açık rıza kaydı, saklama politikası, anonimleştirme işleri
- [ ] Raporlama için veri dışa aktarımı

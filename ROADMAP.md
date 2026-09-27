# platform-api Yol Haritası

Fazlar [ana yol haritası](https://github.com/Libre-University/docs/blob/main/ROADMAP.md) ile hizalıdır. Veri modeli: [DATA_MODEL.md](https://github.com/Libre-University/docs/blob/main/DATA_MODEL.md). Modül yapısı: [ADR-0010](https://github.com/Libre-University/docs/blob/main/docs/adr/0010-develop-modules-as-separate-packages.md).

## Faz 0: İskelet (davet öncesi)

- [ ] Django projesi, `uv` ile bağımlılık yönetimi, `src/` düzeni
- [ ] Çekirdek modül iskeletleri: `identity`, `academic`, `audit`, `notifications`, `files`, `modules`
- [ ] Modül yükleme: `libre_university.modules` entry point'inden `INSTALLED_APPS` ve URL kaydı
- [ ] `libre-university-core` genel arayüz paketinin ilk taslağı
- [ ] Modül şablonu (cookiecutter) ve modül test düzeneği (pytest eklentisi)
- [ ] `import-linter` ile çekirdek iç kodu erişim kuralları
- [ ] CI: ruff, mypy, pytest, migration kontrolü, bağımlılık lisans taraması
- [ ] `/health` ve `/ready` uçları
- [ ] OpenAPI şemasının CI'da üretilip artefakt olarak yayımlanması
- [ ] Yapay (anonim) örnek veri üreten `seed` komutu
- [ ] `CONTRIBUTING.md` içinde yerel geliştirme adımları

## Faz 1: Çekirdek Platform

- [ ] OIDC ile giriş/çıkış (Keycloak, `adapters` içindeki OIDC adapteri üzerinden)
- [ ] `User`, `Role`, `Permission`, `UserRole` ve birim kapsamlı yetki kontrolü
- [ ] `OrganizationUnit` hiyerarşisi, `AcademicTerm`, `Program`, `Course`, `Instructor`
- [ ] `AuditLog`: kritik işlemler, giriş/çıkış, yetki değişiklikleri; güncelleme/silme engelli
- [ ] `Notification` modeli ve e-posta kanalı (SMTP adapteri)
- [ ] `StoredFile` ve S3/MinIO adapteri
- [ ] `libre-university-core` 0.1: yetki, audit, bildirim, dosya, akademik sorgu arayüzleri
- [ ] Prometheus metrikleri ve yapılandırılmış (JSON) loglar

## Faz 2: Akademik İş Akışlarına Destek

- [ ] Olay (domain event) altyapısı: modüller arası gevşek bağlı iletişim
- [ ] `module-obs` ihtiyaçlarına göre çekirdek arayüz genişletmeleri
- [ ] Ders kayıt yoğunluğu için önbellek ve kuyruk altyapısı

## Faz 3: Self-Servis Desteği

- [ ] Bildirim merkezi API'si, push kanalı
- [ ] Mobil uygulama için API ihtiyaçları

## Faz 4+: Kurumsal Süreçler

- [ ] KVKK: açık rıza kaydı, saklama politikası, anonimleştirme işleri
- [ ] Raporlama için veri dışa aktarımı

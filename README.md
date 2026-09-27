# platform-api

LibreUniversity'nin **çekirdek ve modül barındırıcı (host)** backend uygulaması. Python/Django ile yazılır, tek uygulama olarak dağıtılır ([ADR-0004](https://github.com/Libre-University/docs/blob/main/docs/adr/0004-start-with-modular-monolith-for-mvp.md), [ADR-0008](https://github.com/Libre-University/docs/blob/main/docs/adr/0008-use-python-django-and-react.md)). İş modülleri (`module-obs`, `module-lms`, …) ayrı repolarda paket olarak geliştirilir ve bu uygulamaya kurulur ([ADR-0010](https://github.com/Libre-University/docs/blob/main/docs/adr/0010-develop-modules-as-separate-packages.md)).

## Ne İş Yapar?

- Tüm modüllerin ortak kullandığı çekirdek servisleri sağlar: kimlik ve yetki, akademik çekirdek veri, audit, bildirim, dosya.
- Kurulu modül paketlerini `libre_university.modules` entry point'i üzerinden keşfeder ve yükler.
- Modüllerin kullandığı genel arayüz paketini (`libre-university-core`) yayımlar.
- Yeni modül repoları için modül şablonu ve modül test düzeneğini sağlar.
- Tek REST API (OpenAPI 3) sunar; modüllerin uçları da bu API altında toplanır.

## Teknoloji

- Python 3.12+, Django 5.2 LTS, Django REST Framework, drf-spectacular
- PostgreSQL 16+, Celery + Redis/Valkey
- OIDC ile kimlik (Keycloak), parola uygulamada tutulmaz
- pytest, ruff, mypy, import-linter, uv

## Çekirdek Modüller

| Modül | Sorumluluk |
| --- | --- |
| `identity` | Kullanıcı, rol, izin, birim kapsamlı yetki, OIDC oturumu |
| `academic` | Organizasyon birimi, program, dönem, ders kataloğu, akademisyen |
| `audit` | Değiştirilemez denetim kaydı |
| `notifications` | Bildirim kayıtları ve kanal gönderimi |
| `files` | Dosya metadata, S3 uyumlu depolama |
| `modules` | Modül keşfi, kaydı ve sürüm uyumluluk kontrolü |

## Fazlara Göre İşler

| Faz | Bu repoda yapılacaklar |
| --- | --- |
| Faz 0 | Django iskeleti, çekirdek modül klasörleri, modül yükleme mekanizması, modül şablonu, CI, sağlık uçları, OpenAPI üretimi, yapay örnek veri |
| Faz 1 | OIDC girişi, rol/izin ve birim kapsamlı yetki, akademik çekirdek veri, audit log, e-posta bildirimi, dosya depolama, `libre-university-core` 0.1, metrikler |
| Faz 2 | OBS modülünün ihtiyaç duyduğu çekirdek arayüz genişletmeleri, olay (domain event) altyapısı, performans iyileştirmeleri |
| Faz 3 | Bildirim merkezi ve push kanalı, mobil uygulama API ihtiyaçları |
| Faz 4+ | KVKK açık rıza, saklama politikası ve anonimleştirme; raporlama dışa aktarımı |

Ayrıntılı ve işaretlenebilir liste: [ROADMAP.md](ROADMAP.md). Fazlar [ana yol haritası](https://github.com/Libre-University/docs/blob/main/ROADMAP.md) ile hizalıdır. Açık işler için `phase:*` etiketlerine bakın.

## Katkı

Katkı rehberi, davranış kuralları ve güvenlik politikası organizasyon genelinde [`.github`](https://github.com/Libre-University/.github) reposundadır. Mimari kararlar [`docs`](https://github.com/Libre-University/docs) reposundaki ADR'lerle alınır.

## Lisans

Lisans kararı [ADR-0002](https://github.com/Libre-University/docs/blob/main/docs/adr/0002-prefer-agpl-3-or-later-license.md) ile kesinleştirilecektir (öneri: AGPL-3.0-or-later).

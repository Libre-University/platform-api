// Package platform, LibreUniversity çekirdek uygulamasıdır: modülleri
// başlatır, HTTP yönlendiricisini kurar, veritabanı bağlantısını ve
// migration'ları yönetir. Dışarıya yalnızca Config, New ve App açıktır;
// kimlik, akademik çekirdek, audit, bildirim ve dosya bileşenleri
// platform/internal altında kapalıdır (ADR-0012).
package platform

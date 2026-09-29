// Package core, iş modüllerinin platforma bağlanmak için kullandığı genel
// arayüzdür. Modüller çekirdeğe yalnızca bu paket üzerinden erişir; çekirdeğin
// geri kalanı platform/internal altında derleyici tarafından kapalıdır
// (ADR-0012).
//
// Bir modül core.Module arayüzünü uygular ve init() içinde core.Register ile
// kendini kaydeder. Dağıtım uygulaması (cmd/libre-university) modülleri içe
// aktarır ve core.Registered() ile platforma verir.
package core

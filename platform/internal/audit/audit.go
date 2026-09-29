// Package audit, çekirdeğin denetim kaydı uygulamasıdır. Faz 0'da olaylar
// yapılandırılmış log olarak yazılır; Faz 1'de değiştirilemez audit tablosu
// gelir (DATA_MODEL.md › AuditLog).
package audit

import (
	"context"
	"log/slog"

	"github.com/Libre-University/platform-api/core"
)

// LogAuditor olayları slog'a yazar.
type LogAuditor struct{ log *slog.Logger }

// NewLogAuditor verilen logger ile LogAuditor üretir.
func NewLogAuditor(log *slog.Logger) *LogAuditor { return &LogAuditor{log: log} }

// Record olayı "audit" seviyesinde loglar.
func (l *LogAuditor) Record(ctx context.Context, e core.AuditEvent) error {
	l.log.InfoContext(ctx, "audit",
		"actor", e.Actor,
		"action", e.Action,
		"entity_type", e.EntityType,
		"entity_id", e.EntityID,
		"reason", e.Reason,
	)
	return nil
}

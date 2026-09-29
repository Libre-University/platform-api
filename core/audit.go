package core

import "context"

// AuditEvent, kritik bir işlemin denetim kaydıdır (DATA_MODEL.md › AuditLog).
type AuditEvent struct {
	// Actor işlemi yapan kullanıcının kimliğidir; sistem işlerinde boş olabilir.
	Actor string
	// Action "<modül>.<varlık>.<eylem>" biçimindedir.
	Action     string
	EntityType string
	EntityID   string
	// Reason not değişikliği, kayıt iptali gibi işlemlerde zorunludur.
	Reason string
	Before map[string]any
	After  map[string]any
}

// Auditor denetim kaydı yazar. Çekirdek tarafından sağlanır; modüller kendi
// audit tablosu tutmaz.
type Auditor interface {
	Record(ctx context.Context, e AuditEvent) error
}

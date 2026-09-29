-- +goose Up
-- Çekirdek şeması platform tarafından oluşturulur; Faz 1'de kimlik, akademik
-- çekirdek, audit, bildirim ve dosya tabloları buraya gelir (DATA_MODEL.md).
COMMENT ON SCHEMA platform IS 'LibreUniversity çekirdek: kimlik, akademik çekirdek, audit, bildirim, dosya';

-- +goose Down
SELECT 1;

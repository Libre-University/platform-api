-- +goose Up
-- lms şeması platform tarafından oluşturulur. Tablolar ADR-0013 kararına göre
-- Faz 2-3'te gelir: entegrasyon modelinde CourseLink ve GradeImport; yerli LMS
-- modelinde CoursePage, LearningMaterial, Assignment, AssignmentSubmission.
COMMENT ON SCHEMA lms IS 'Öğrenme yönetim sistemi modülü (module: lms)';

-- +goose Down
SELECT 1;

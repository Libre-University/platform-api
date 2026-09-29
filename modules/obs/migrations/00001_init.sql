-- +goose Up
-- obs şeması platform tarafından oluşturulur. Faz 1'de StudentRecord,
-- Curriculum, CurriculumCourse; Faz 2'de CourseSection, SectionMeeting,
-- TermRegistration, Enrollment, GradeItem, Grade, GradingPolicy, CourseResult
-- tabloları gelir (DATA_MODEL.md).
COMMENT ON SCHEMA obs IS 'Öğrenci Bilgi Sistemi modülü (module: obs)';

-- +goose Down
SELECT 1;

ALTER TABLE IF EXISTS "student_group_stats" DROP CONSTRAINT IF EXISTS "student_group_stats_pkey";
ALTER TABLE "student_group_stats" DROP CONSTRAINT IF EXISTS "student_group_stats_student_id_test_group_id_idx";
ALTER TABLE IF EXISTS "student_group_stats" DROP CONSTRAINT IF EXISTS "student_id";
ALTER TABLE IF EXISTS "tests" DROP CONSTRAINT IF EXISTS "group_id";
ALTER TABLE IF EXISTS "test_groups" DROP CONSTRAINT IF EXISTS "teacher_id";

DROP TABLE IF EXISTS "student_answers";
DROP TABLE IF EXISTS "students";
DROP TABLE IF EXISTS "tests";
DROP TABLE IF EXISTS "test_groups";
DROP TABLE IF EXISTS "teachers";
DROP TABLE IF EXISTS "users";
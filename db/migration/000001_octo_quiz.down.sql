-- Down-migration script to drop the Telegram Quizlet bot database schema
-- Target: PostgreSQL
-- Note: Drops tables in reverse dependency order to respect foreign key constraints

-- Drop test_sessions (depends on users, test_sets)
DROP TABLE IF EXISTS test_sessions;

-- Drop test_progress (depends on users, test_sets, words)
DROP TABLE IF EXISTS test_progress;

-- Drop student_progress (depends on users, test_sets, words)
DROP TABLE IF EXISTS student_progress;

-- Drop words (depends on test_sets)
DROP TABLE IF EXISTS words;

-- Drop class_test_sets (depends on classes, test_sets)
DROP TABLE IF EXISTS class_test_sets;

-- Drop test_sets (depends on users)
DROP TABLE IF EXISTS test_sets;

-- Drop class_students (depends on classes, users)
DROP TABLE IF EXISTS class_students;

-- Drop classes (depends on users)
DROP TABLE IF EXISTS classes;

-- Drop users (no dependencies)
DROP TABLE IF EXISTS users;
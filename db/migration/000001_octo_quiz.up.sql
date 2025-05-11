CREATE TABLE "users" (
  "user_id" bigint PRIMARY KEY,
  "telegram_id" bigint NOT NULL,
  "full_name" varchar NOT NULL,
  "username" varchar NOT NULL,
  "role" varchar NOT NULL,
  "phone" bigint NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "classes" (
  "class_id" bigserial PRIMARY KEY,
  "class_name" varchar NOT NULL,
  "teacher_id" bigint NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "class_students" (
  "class_id" bigint NOT NULL,
  "student_id" bigint NOT NULL,
  "added_at" timestamptz DEFAULT (now()),
  PRIMARY KEY ("class_id", "student_id")
);

CREATE TABLE "test_sets" (
  "test_set_id" bigserial PRIMARY KEY,
  "test_set_name" varchar NOT NULL,
  "creator_id" bigint NOT NULL,
  "is_public" boolean NOT NULL DEFAULT false,
  "time_limit" integer DEFAULT 0
);

CREATE TABLE "class_test_sets" (
  "class_id" bigint NOT NULL,
  "test_set_id" bigint NOT NULL,
  PRIMARY KEY ("class_id", "test_set_id")
);

CREATE TABLE "words" (
  "words_id" bigserial PRIMARY KEY,
  "test_set_id" bigint NOT NULL,
  "english_word" varchar NOT NULL,
  "uzbek_word" varchar NOT NULL
);

CREATE TABLE "student_progress" (
  "progress_id" bigserial PRIMARY KEY,
  "student_id" bigint NOT NULL,
  "test_set_id" bigint NOT NULL,
  "words_id" bigint NOT NULL,
  "correct_count" integer DEFAULT 0,
  "incorrect_count" integer DEFAULT 0
);

CREATE TABLE "test_progress" (
  "progress_id" bigserial PRIMARY KEY,
  "student_id" bigint NOT NULL,
  "test_set_id" bigint NOT NULL,
  "words_id" bigint NOT NULL,
  "completed" boolean NOT NULL DEFAULT false
);

CREATE TABLE "test_sessions" (
  "session_id" bigserial PRIMARY KEY,
  "student_id" bigint NOT NULL,
  "test_set_id" bigint NOT NULL,
  "start_time" timestamptz NOT NULL DEFAULT (now()),
  "correct_count" integer NOT NULL DEFAULT 0,
  "incorrect_count" integer NOT NULL DEFAULT 0,
  "completed" boolean NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX ON "users" ("telegram_id", "username", "phone");

CREATE UNIQUE INDEX ON "words" ("words_id", "test_set_id");

CREATE UNIQUE INDEX ON "student_progress" ("student_id", "test_set_id", "words_id");

CREATE UNIQUE INDEX ON "test_progress" ("student_id", "test_set_id", "words_id");

COMMENT ON TABLE "users" IS 'CHECK (role IN ("student", "teacher"))';

COMMENT ON COLUMN "users"."role" IS 'Must be student or teacher';

COMMENT ON COLUMN "test_sets"."time_limit" IS 'Seconds, NULL if no limit';

COMMENT ON COLUMN "test_sessions"."start_time" IS 'Unix timestamp';

COMMENT ON COLUMN "test_sessions"."completed" IS 'True if all words completed';

ALTER TABLE "classes" ADD FOREIGN KEY ("teacher_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;

ALTER TABLE "class_students" ADD FOREIGN KEY ("class_id") REFERENCES "classes" ("class_id") ON DELETE CASCADE;

ALTER TABLE "class_students" ADD FOREIGN KEY ("student_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;

ALTER TABLE "test_sets" ADD FOREIGN KEY ("creator_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;

ALTER TABLE "class_test_sets" ADD FOREIGN KEY ("class_id") REFERENCES "classes" ("class_id") ON DELETE CASCADE;

ALTER TABLE "class_test_sets" ADD FOREIGN KEY ("test_set_id") REFERENCES "test_sets" ("test_set_id") ON DELETE CASCADE;

ALTER TABLE "words" ADD FOREIGN KEY ("test_set_id") REFERENCES "test_sets" ("test_set_id") ON DELETE CASCADE;

ALTER TABLE "student_progress" ADD FOREIGN KEY ("student_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;

ALTER TABLE "student_progress" ADD FOREIGN KEY ("test_set_id") REFERENCES "test_sets" ("test_set_id") ON DELETE CASCADE;

ALTER TABLE "student_progress" ADD FOREIGN KEY ("words_id", "test_set_id") REFERENCES "words" ("words_id", "test_set_id") ON DELETE CASCADE;

ALTER TABLE "test_progress" ADD FOREIGN KEY ("student_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;

ALTER TABLE "test_progress" ADD FOREIGN KEY ("test_set_id") REFERENCES "test_sets" ("test_set_id") ON DELETE CASCADE;

ALTER TABLE "test_progress" ADD FOREIGN KEY ("words_id", "test_set_id") REFERENCES "words" ("words_id", "test_set_id") ON DELETE CASCADE;

ALTER TABLE "test_sessions" ADD FOREIGN KEY ("student_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;

ALTER TABLE "test_sessions" ADD FOREIGN KEY ("test_set_id") REFERENCES "test_sets" ("test_set_id") ON DELETE CASCADE;

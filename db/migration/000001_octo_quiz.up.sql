CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "telegram_id" bigint UNIQUE NOT NULL,
  "full_name" varchar NOT NULL,
  "username" varchar UNIQUE NOT NULL,
  "phone_number" varchar UNIQUE NOT NULL
);

CREATE TABLE "teachers" (
  "id" bigserial PRIMARY KEY,
  "telegram_id" bigint UNIQUE NOT NULL,
  "full_name" varchar NOT NULL
);

CREATE TABLE "test_groups" (
  "id" bigserial PRIMARY KEY,
  "teacher_id" bigint NOT NULL,
  "group_name" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "tests" (
  "id" bigserial PRIMARY KEY,
  "group_id" bigint NOT NULL,
  "english_word" varchar NOT NULL,
  "uzbek_word" varchar NOT NULL
);

CREATE TABLE "students" (
  "id" bigserial PRIMARY KEY,
  "telegram_id" bigint UNIQUE NOT NULL,
  "full_name" varchar NOT NULL
);

CREATE TABLE "student_group_stats" (
  "id" bigserial PRIMARY KEY,
  "student_id" bigint NOT NULL,
  "test_group_id" bigint NOT NULL,
  "total_correct" integer DEFAULT 0
);

CREATE UNIQUE INDEX ON "student_group_stats" ("student_id", "test_group_id");

ALTER TABLE "teachers" ADD FOREIGN KEY ("telegram_id") REFERENCES "users" ("id");

ALTER TABLE "test_groups" ADD FOREIGN KEY ("teacher_id") REFERENCES "teachers" ("id");

ALTER TABLE "tests" ADD FOREIGN KEY ("group_id") REFERENCES "test_groups" ("id");

ALTER TABLE "students" ADD FOREIGN KEY ("telegram_id") REFERENCES "users" ("id");

ALTER TABLE "student_group_stats" ADD FOREIGN KEY ("student_id") REFERENCES "students" ("id");

ALTER TABLE "student_group_stats" ADD FOREIGN KEY ("test_group_id") REFERENCES "test_groups" ("id");

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "avatar" varchar,
  "username" varchar NOT NULL UNIQUE,
  "name" varchar,
  "email" varchar NOT NULL UNIQUE,
  "password" varchar NOT NULL,
  "phone" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "last_login" timestamptz NOT NULL
);
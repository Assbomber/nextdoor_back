CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "avatar" varchar,
  "name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "password" varchar NOT NULL,
  "phone" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

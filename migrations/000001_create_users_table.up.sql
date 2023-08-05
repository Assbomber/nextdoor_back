CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TYPE "genders" AS ENUM ('male','female','others');

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "avatar" varchar,
  "username" varchar NOT NULL UNIQUE,
  "name" varchar,
  "email" varchar NOT NULL UNIQUE,
  "password" varchar NOT NULL,
  "phone" varchar,
  "gender" genders,
  "birth_date" date,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "last_login" timestamptz NOT NULL,
  "is_deleted" boolean NOT NULL DEFAULT false
);

CREATE TABLE "users_locations" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL REFERENCES users(id),
  "location" geometry(POINT, 4326) NOT NULL, -- 4326 is the SRID (Spatial Reference ID) for WGS 84, a commonly used GPS coordinate system
  "active" boolean NOT NULL DEFAULT false 
);

CREATE UNIQUE INDEX idx_unique_active_address
ON "users_locations" (user_id)
WHERE active;
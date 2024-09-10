-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2024-09-09T09:48:39.179Z

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE,
  "hashed_password" varchar NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "phone_number" varchar NOT NULL,
  "profile_photo" bigint NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "is_email_verified" bool NOT NULL DEFAULT false,
  "isAdmin" boolean DEFAULT false,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "verify_emails" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "email" varchar NOT NULL,
  "secret_code" varchar NOT NULL,
  "is_used" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "expired_at" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);

CREATE TABLE "media" (
  "id" bigserial PRIMARY KEY,
  "media_ref" varchar UNIQUE NOT NULL,
  "url" varchar UNIQUE NOT NULL,
  "aws_id" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "product_ref_no" varchar UNIQUE NOT NULL,
  "product_name" varchar UNIQUE NOT NULL,
  "price" numeric NOT NULL,
  "product_images" varchar NOT NULL,
  "collection" bigint NOT NULL,
  "quantity" int NOT NULL,
  "status" enum(available,out_of_stock,discontinued) NOT NULL,
  "last_updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "collections" (
  "id" bigserial PRIMARY KEY,
  "collection_name" varchar UNIQUE NOT NULL,
  "last_updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "ref_no" varchar,
  "username" varchar NOT NULL,
  "amount" bigint NOT NULL,
  "payment_method" varchar NOT NULL,
  "product" varchar NOT NULL,
  "order_status" enum(pending,shipped,delivered,cancelled) NOT NULL,
  "last_updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "product_media" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "product_media_ref" varchar UNIQUE NOT NULL,
  "product_id" varchar NOT NULL,
  "media_id" varchar NOT NULL
);

CREATE INDEX ON "users" ("profile_photo");

CREATE INDEX ON "products" ("collection");

CREATE INDEX ON "products" ("product_images");

CREATE INDEX ON "products" ("collection", "product_images");

CREATE INDEX ON "orders" ("username");

CREATE INDEX ON "orders" ("product");

CREATE INDEX ON "orders" ("username", "product");

CREATE INDEX ON "product_media" ("product_id");

CREATE INDEX ON "product_media" ("media_id");

CREATE INDEX ON "product_media" ("product_id", "media_id");

COMMENT ON COLUMN "orders"."amount" IS 'it must be positive';

ALTER TABLE "users" ADD FOREIGN KEY ("profile_photo") REFERENCES "media" ("id");

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "products" ADD FOREIGN KEY ("product_images") REFERENCES "product_media" ("product_media_ref");

ALTER TABLE "products" ADD FOREIGN KEY ("collection") REFERENCES "collections" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "orders" ADD FOREIGN KEY ("product") REFERENCES "products" ("product_ref_no");

ALTER TABLE "product_media" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("product_ref_no");

ALTER TABLE "product_media" ADD FOREIGN KEY ("media_id") REFERENCES "media" ("media_ref");

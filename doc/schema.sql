-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2024-10-13T18:38:42.580Z

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE,
  "hashed_password" varchar NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "phone_number" varchar NOT NULL,
  "profile_photo" varchar,
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
  "product_description" varchar NOT NULL,
  "product_code" varchar NOT NULL,
  "price" numeric NOT NULL,
  "sale_price" varchar NOT NULL,
  "product_image_main" varchar NOT NULL,
  "product_image_other_1" varchar NOT NULL,
  "product_image_other_2" varchar NOT NULL,
  "product_image_other_3" varchar NOT NULL,
  "collection" bigint NOT NULL,
  "quantity" int NOT NULL,
  "color" varchar NOT NULL,
  "size" varchar NOT NULL,
  "status" enum(active,out_of_stock,archived) NOT NULL,
  "last_updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "collections" (
  "id" bigserial PRIMARY KEY,
  "collection_name" varchar UNIQUE NOT NULL,
  "collection_description" varchar NOT NULL,
  "product_count" int,
  "thumbnail_image" varchar NOT NULL,
  "header_image" varchar NOT NULL,
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
  "is_main_image" boolean NOT NULL DEFAULT false,
  "media_id" varchar NOT NULL
);

CREATE TABLE "cart" (
  "id" bigserial PRIMARY KEY,
  "user_ref_id" bigint NOT NULL,
  "product_id" int NOT NULL,
  "product_name" varchar NOT NULL,
  "product_price" varchar NOT NULL,
  "product_quantity" bigint NOT NULL,
  "product_image" varchar NOT NULL,
  "product_color" varchar NOT NULL,
  "product_size" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "products" ("collection");

CREATE INDEX ON "products" ("product_image_main");

CREATE INDEX ON "products" ("product_image_other_1");

CREATE INDEX ON "products" ("product_image_other_2");

CREATE INDEX ON "products" ("product_image_other_3");

CREATE INDEX ON "products" ("collection", "product_image_main", "product_image_other_1", "product_image_other_2", "product_image_other_3");

CREATE INDEX ON "orders" ("username");

CREATE INDEX ON "orders" ("product");

CREATE INDEX ON "orders" ("username", "product");

CREATE INDEX ON "product_media" ("product_id");

CREATE INDEX ON "product_media" ("media_id");

CREATE INDEX ON "product_media" ("product_id", "media_id");

COMMENT ON COLUMN "orders"."amount" IS 'it must be positive';

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "products" ADD FOREIGN KEY ("product_image_main") REFERENCES "product_media" ("product_media_ref");

ALTER TABLE "products" ADD FOREIGN KEY ("product_image_other_1") REFERENCES "product_media" ("product_media_ref");

ALTER TABLE "products" ADD FOREIGN KEY ("product_image_other_2") REFERENCES "product_media" ("product_media_ref");

ALTER TABLE "products" ADD FOREIGN KEY ("product_image_other_3") REFERENCES "product_media" ("product_media_ref");

ALTER TABLE "products" ADD FOREIGN KEY ("collection") REFERENCES "collections" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "orders" ADD FOREIGN KEY ("product") REFERENCES "products" ("product_ref_no");

ALTER TABLE "product_media" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("product_ref_no");

ALTER TABLE "product_media" ADD FOREIGN KEY ("media_id") REFERENCES "media" ("media_ref");

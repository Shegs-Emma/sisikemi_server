CREATE TYPE product_status AS ENUM ('available', 'out_of_stock', 'discontinued');
CREATE TYPE order_status AS ENUM ('pending', 'shipped', 'delivered', 'cancelled');


CREATE TABLE "users" (
  "username" varchar UNIQUE NOT NULL PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "phone_number" varchar NOT NULL,
  "profile_photo" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "is_admin" boolean DEFAULT false NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
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
  "product_code" varchar UNIQUE NOT NULL,
  "price" bigint NOT NULL,
  "sale_price" varchar NOT NULL,
  "product_image_main" varchar,
  "product_image_other_1" varchar,
  "product_image_other_2" varchar,
  "product_image_other_3" varchar,
  "collection" bigint NOT NULL,
  "quantity" int NOT NULL,
  "color" varchar NOT NULL,
  "size" varchar NOT NULL,
  "status" product_status NOT NULL,
  "last_updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "collections" (
  "id" bigserial PRIMARY KEY,
  "collection_name" varchar UNIQUE NOT NULL,
  "collection_description" varchar NOT NULL,
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
  "order_status" order_status NOT NULL,
  "last_updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "product_media" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "product_media_ref" varchar UNIQUE NOT NULL,
  "product_id" varchar NOT NULL,
  "is_main_image" boolean DEFAULT false NOT NULL,
  "media_id" varchar NOT NULL
);

CREATE INDEX ON "users" ("profile_photo");

CREATE INDEX ON "products" ("collection");

CREATE INDEX ON "products" ("collection");

CREATE INDEX ON "orders" ("username");

CREATE INDEX ON "orders" ("product");

CREATE INDEX ON "orders" ("username", "product");

CREATE INDEX ON "product_media" ("product_id");

CREATE INDEX ON "product_media" ("media_id");

CREATE INDEX ON "product_media" ("product_id", "media_id");

COMMENT ON COLUMN "orders"."amount" IS 'it must be positive';

ALTER TABLE "users" ADD FOREIGN KEY ("profile_photo") REFERENCES "media" ("media_ref");

ALTER TABLE "products" ADD FOREIGN KEY ("collection") REFERENCES "collections" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "orders" ADD FOREIGN KEY ("product") REFERENCES "products" ("product_ref_no");

ALTER TABLE "product_media" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("product_ref_no");

ALTER TABLE "product_media" ADD FOREIGN KEY ("media_id") REFERENCES "media" ("media_ref");
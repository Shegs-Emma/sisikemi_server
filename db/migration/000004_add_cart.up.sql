CREATE TABLE "cart" (
  "id" bigserial PRIMARY KEY,
  "product_id" int NOT NULL,
  "product_name" varchar NOT NULL,
  "product_price" bigint NOT NULL,
  "product_quantity" bigint NOT NULL,
  "product_image" varchar NOT NULL,
  "product_color" varchar NOT NULL,
  "product_size" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
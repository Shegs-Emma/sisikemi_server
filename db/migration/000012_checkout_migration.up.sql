CREATE TABLE "shipping_address" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "country" varchar NOT NULL,
  "address" varchar NOT NULL,
  "town" varchar NOT NULL,
  "postal_code" varchar,
  "landmark" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "expired_at" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);

ALTER TABLE "shipping_address" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "orders" ADD COLUMN "shipping_method" varchar NOT NULL;
ALTER TABLE "orders" ADD COLUMN "shipping_address_id" bigint NOT NULL;
ALTER TABLE "orders" ADD FOREIGN KEY ("shipping_address_id") REFERENCES "shipping_address" ("id");
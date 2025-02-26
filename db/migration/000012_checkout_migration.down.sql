
ALTER TABLE "orders" DROP COLUMN "shipping_address_id";
ALTER TABLE "orders" DROP COLUMN "shipping_method";

DROP TABLE IF EXISTS "shipping_address" CASCADE;
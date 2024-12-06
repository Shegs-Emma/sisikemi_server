ALTER TABLE "products" DROP COLUMN "color";
ALTER TABLE "products" ADD COLUMN "color" varchar[] NOT NULL;
ALTER TABLE "products" DROP COLUMN "size";
ALTER TABLE "products" ADD COLUMN "size" varchar[] NOT NULL;
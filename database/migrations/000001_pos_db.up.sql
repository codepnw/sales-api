BEGIN;

SET TIME ZONE 'Asia/Bangkok';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CREATE SEQUENCE
CREATE SEQUENCE seq_user_id START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE seq_product_id START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE seq_customer_id START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE seq_order_id START WITH 1 INCREMENT BY 1;

-- CREATE TYPE ENUM
CREATE TYPE "enum_order_status" AS ENUM ('WAITING', 'COMPLETED', 'CANCEL');
CREATE TYPE "enum_payment_method" AS ENUM ('CASH', 'TRANSFER', 'ETC');

-- CREATE TABLE
CREATE TABLE "users" (
  "user_id" VARCHAR(7) PRIMARY KEY DEFAULT CONCAT('U', LPAD(NEXTVAL('seq_user_id')::TEXT, 6, '0')),
  "email" VARCHAR UNIQUE NOT NULL,
  "username" VARCHAR UNIQUE NOT NULL,
  "password" VARCHAR NOT NULL,
  "role_id" INT NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE "user_roles" (
  "role_id" SERIAL PRIMARY KEY,
  "name" VARCHAR UNIQUE NOT NULL,
  "desc" VARCHAR
);

CREATE TABLE "products" (
  "product_id" VARCHAR(7) PRIMARY KEY DEFAULT CONCAT('P', LPAD(NEXTVAL('seq_product_id')::TEXT, 6, '0')),
  "name" VARCHAR NOT NULL,
  "desc" VARCHAR,
  "price" FLOAT DEFAULT 0,
  "discount" FLOAT DEFAULT 0,
  "stock" INT DEFAULT 0,
  "category_id" INT NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE "categories" (
  "category_id" SERIAL PRIMARY KEY,
  "title" VARCHAR NOT NULL,
  "desc" VARCHAR
);

CREATE TABLE "customers" (
  "customer_id" VARCHAR(7) PRIMARY KEY DEFAULT CONCAT('C', LPAD(NEXTVAL('seq_customer_id')::TEXT, 6, '0')),
  "first_name" VARCHAR NOT NULL,
  "last_name" VARCHAR NOT NULL,
  "phone" VARCHAR NOT NULL UNIQUE,
  "email" VARCHAR NOT NULL UNIQUE,
  "address" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE "orders" (
  "order_id" VARCHAR(7) PRIMARY KEY DEFAULT CONCAT('O', LPAD(NEXTVAL('seq_order_id')::TEXT, 6, '0')),
  "customer_id" VARCHAR NOT NULL,
  "total_amount" FLOAT NOT NULL,
  "payment_method" enum_payment_method NOT NULL,
  "status" enum_order_status NOT NULL,
  "order_date" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE "order_items" (
  "order_item_id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "order_id" VARCHAR NOT NULL,
  "product_id" VARCHAR NOT NULL,
  "quantity" INT NOT NULL DEFAULT 1,
  "price" FLOAT NOT NULL DEFAULT 0,
  "discount" INT DEFAULT 0
);

CREATE TABLE "payments" (
  "payment_id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "order_id" VARCHAR NOT NULL,
  "amount" FLOAT NOT NULL DEFAULT 0,
  "payment_method" enum_payment_method NOT NULL,
  "payment_date" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE "inventory_logs" (
  "inventory_log_id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "product_id" VARCHAR NOT NULL,
  "change" VARCHAR NOT NULL,
  "description" VARCHAR,
  "date" TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "user_roles" ("role_id");

ALTER TABLE "products" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("category_id");

ALTER TABLE "orders" ADD FOREIGN KEY ("customer_id") REFERENCES "customers" ("customer_id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("order_id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("product_id");

ALTER TABLE "payments" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("order_id");

ALTER TABLE "inventory_logs" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("product_id");

COMMIT;
CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint,
  "amount" bigint NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint,
  "to_account_id" bigint,
  "amount" bigint,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

-- 1. 删除现有的外键约束
ALTER TABLE "entries" DROP CONSTRAINT IF EXISTS entries_account_id_fkey;
ALTER TABLE "transfers" DROP CONSTRAINT IF EXISTS transfers_from_account_id_fkey;
ALTER TABLE "transfers" DROP CONSTRAINT IF EXISTS transfers_to_account_id_fkey;

-- 2. 添加新的外键约束，启用级联删除
ALTER TABLE "entries" ADD CONSTRAINT entries_account_id_fkey FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;
ALTER TABLE "transfers" ADD CONSTRAINT transfers_from_account_id_fkey FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;
ALTER TABLE "transfers" ADD CONSTRAINT transfers_to_account_id_fkey FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;


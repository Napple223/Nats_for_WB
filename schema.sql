/*схема БД*/

DROP TABLE IF EXISTS orders, delivery, items, payment CASCADE;

CREATE TABLE IF NOT EXISTS orders (
    order_uid TEXT PRIMARY KEY UNIQUE,
    track_number TEXT NOT NULL,
    "entry" TEXT NOT NULL,
    locale TEXT,
    internal_signature TEXT,
    customer_id TEXT,
    delivery_service TEXT,
    shardkey TEXT,
    sm_id INTEGER,
    date_created TEXT,
    oof_shard TEXT 
);

CREATE TABLE IF NOT EXISTS delivery (
    order_ref TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
    "name" TEXT NOT NULL,
    phone TEXT,
    zip TEXT NOT NULL,
    city TEXT NOT NULL,
    "address" TEXT NOT NULL,
    region TEXT NOT NULL,
    email TEXT
);

CREATE TABLE IF NOT EXISTS items (
    order_ref TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id INTEGER NOT NULL,
    track_number TEXT NOT NULL,
    price INTEGER NOT NULL,
    rid TEXT NOT NULL,
    "name" TEXT NOT NULL,
    sale INTEGER DEFAULT 0,
    size TEXT NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id INTEGER NOT NULL,
    brand TEXT NOT NULL,
    status INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS payment (
    order_ref TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
    "transaction" TEXT NOT NULL,
    request_id TEXT,
    currency TEXT NOT NULL,
    "provider" TEXT NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt INTEGER NOT NULL,
    bank TEXT NOT NULL,
    delivery_cost INTEGER DEFAULT 0,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER
)
-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
order_uid VARCHAR(20) UNIQUE NOT NULL PRIMARY KEY,
track_number VARCHAR(15) UNIQUE NOT NULL,
entry VARCHAR(4) NOT NULL,
locale VARCHAR(2) NOT NULL,
internal_signature VARCHAR(20) NOT NULL,
customer_id VARCHAR(20) NOT NULL,
delivery_service VARCHAR(20) NOT NULL,
shardkey VARCHAR(10) NOT NULL,
sm_id INT NOT NULL,
date_created TIMESTAMP NOT NULL DEFAULT NOW(),
oof_shard VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS deliveries (
order_uid VARCHAR(19) UNIQUE NOT NULL REFERENCES orders(order_uid),
name VARCHAR(255) NOT NULL,
phone VARCHAR(15) NOT NULL,
zip VARCHAR(10) NOT NULL,
city VARCHAR(255) NOT NULL,
address VARCHAR(255) NOT NULL,
region VARCHAR(255) NOT NULL,
email VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS payments (
transaction VARCHAR(20) UNIQUE NOT NULL REFERENCES orders(order_uid),
request_id VARCHAR(19) NOT NULL,
currency VARCHAR(3) NOT NULL,
provider VARCHAR(100) NOT NULL,
amount INT NOT NULL,
payment_dt INT NOT NULL,
bank VARCHAR(100) NOT NULL,
delivery_cost INT NOT NULL,
goods_total INT NOT NULL,
custom_fee INT NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
chrt_id INT NOT NULL PRIMARY KEY,
track_number VARCHAR(15),
price INT NOT NULL,
rid VARCHAR(21) NOT NULL,
name VARCHAR(255) NOT NULL,
sale INT NOT NULL,
size VARCHAR(255) NOT NULL,
total_price INT NOT NULL,
nm_id INT NOT NULL,
brand VARCHAR(100) NOT NULL,
status INT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS payments, items, deliveries, orders;

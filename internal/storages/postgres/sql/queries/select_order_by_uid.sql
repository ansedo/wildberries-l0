SELECT
    orders.order_uid,
    orders.track_number,
    orders.entry,
    orders.locale,
    orders.internal_signature,
    orders.customer_id,
    orders.delivery_service,
    orders.shardkey,
    orders.sm_id,
    orders.date_created,
    orders.oof_shard,
    deliveries.name "delivery.name",
    deliveries.phone "delivery.phone",
    deliveries.zip "delivery.zip",
    deliveries.city "delivery.city",
    deliveries.address "delivery.address",
    deliveries.region "delivery.region",
    deliveries.email "delivery.email",
    payments.transaction "payment.transaction",
    payments.request_id "payment.request_id",
    payments.currency "payment.currency",
    payments.provider "payment.provider",
    payments.amount "payment.amount",
    payments.payment_dt "payment.payment_dt",
    payments.bank "payment.bank",
    payments.delivery_cost "payment.delivery_cost",
    payments.goods_total "payment.goods_total",
    payments.custom_fee "payment.custom_fee"
FROM
    orders
INNER JOIN
    deliveries ON deliveries.order_uid = orders.order_uid
INNER JOIN
    payments ON payments.transaction = orders.order_uid
WHERE
    orders.order_uid = $1

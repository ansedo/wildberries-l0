INSERT INTO orders(
    order_uid,
    track_number,
    entry,
    locale,
    internal_signature,
    customer_id,
    delivery_service,
    shardkey,
    sm_id,
    date_created,
    oof_shard
) VALUES (
    :order_uid,
    :track_number,
    :entry,
    :locale,
    :internal_signature,
    :customer_id,
    :delivery_service,
    :shardkey,
    :sm_id,
    :date_created,
    :oof_shard
)

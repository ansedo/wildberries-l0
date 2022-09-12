INSERT INTO payments(
    transaction,
    request_id,
    currency,
    provider,
    amount,
    payment_dt,
    bank,
    delivery_cost,
    goods_total,
    custom_fee
) VALUES (
    :payment.transaction,
    :payment.request_id,
    :payment.currency,
    :payment.provider,
    :payment.amount,
    :payment.payment_dt,
    :payment.bank,
    :payment.delivery_cost,
    :payment.goods_total,
    :payment.custom_fee
)

INSERT INTO deliveries(
    order_uid,
    name,
    phone,
    zip,
    city,
    address,
    region,
    email
) VALUES (
    :order_uid,
    :delivery.name,
    :delivery.phone,
    :delivery.zip,
    :delivery.city,
    :delivery.address,
    :delivery.region,
    :delivery.email
)

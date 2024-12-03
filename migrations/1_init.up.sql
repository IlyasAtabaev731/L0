CREATE TABLE orders
(
    order_uid          VARCHAR(255) PRIMARY KEY,
    track_number       VARCHAR(255),
    entry              VARCHAR(255),
    locale             VARCHAR(255),
    internal_signature VARCHAR(255),
    customer_id        VARCHAR(255),
    delivery_service   VARCHAR(255),
    shardkey           VARCHAR(255),
    sm_id              INTEGER,
    date_created       TIMESTAMP WITH TIME ZONE,
    oof_shard          VARCHAR(255)
);

CREATE TABLE deliveries
(
    order_uid VARCHAR(255),
    name      VARCHAR(255),
    phone     VARCHAR(255),
    zip       VARCHAR(255),
    city      VARCHAR(255),
    address   VARCHAR(200),
    region    VARCHAR(255),
    email     VARCHAR(255),
    PRIMARY KEY (order_uid),
    CONSTRAINT fk_orders_delivery
        FOREIGN KEY (order_uid)
            REFERENCES orders(order_uid)
            ON DELETE CASCADE
);

CREATE TABLE payments
(
    order_uid      VARCHAR(255),
    transaction_id VARCHAR(255),
    request_id     VARCHAR(255),
    currency       VARCHAR(255),
    provider       VARCHAR(255),
    amount         INTEGER,
    payment_dt     INTEGER,
    bank           VARCHAR(255),
    delivery_cost  INTEGER,
    goods_total    INTEGER,
    custom_fee     INTEGER,
    PRIMARY KEY (order_uid),
    CONSTRAINT fk_orders_payment
        FOREIGN KEY (order_uid)
            REFERENCES orders(order_uid)
            ON DELETE CASCADE
);

CREATE TABLE items
(
    item_id      SERIAL PRIMARY KEY,
    order_uid    VARCHAR(255),
    chrt_id      INTEGER,
    track_number VARCHAR(255),
    price        INTEGER,
    rid          VARCHAR(255),
    name         VARCHAR(255),
    sale         INTEGER,
    size         VARCHAR(255),
    total_price  INTEGER,
    nm_id        INTEGER,
    brand        VARCHAR(255),
    status       INTEGER,
    CONSTRAINT fk_orders_items
        FOREIGN KEY (order_uid)
            REFERENCES orders(order_uid)
            ON DELETE CASCADE
);
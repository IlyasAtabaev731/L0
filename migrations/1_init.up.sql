CREATE TABLE orders
(
    order_uid          VARCHAR(50)              NOT NULL PRIMARY KEY,
    track_number       VARCHAR(50)              NOT NULL,
    entry              VARCHAR(10)              NOT NULL,
    locale             VARCHAR(5)               NOT NULL,
    internal_signature VARCHAR(100),
    customer_id        VARCHAR(50)              NOT NULL,
    delivery_service   VARCHAR(20)              NOT NULL,
    shardkey           VARCHAR(10)              NOT NULL,
    sm_id              INTEGER                  NOT NULL,
    date_created       TIMESTAMP WITH TIME ZONE NOT NULL,
    oof_shard          VARCHAR(10)              NOT NULL
);

CREATE TABLE deliveries
(
    order_uid VARCHAR(50)  NOT NULL,
    name      VARCHAR(100) NOT NULL,
    phone     VARCHAR(20)  NOT NULL,
    zip       VARCHAR(20)  NOT NULL,
    city      VARCHAR(100) NOT NULL,
    address   VARCHAR(200) NOT NULL,
    region    VARCHAR(100) NOT NULL,
    email     VARCHAR(100) NOT NULL,
    PRIMARY KEY (order_uid),
    CONSTRAINT fk_orders_delivery
        FOREIGN KEY (order_uid)
            REFERENCES orders (order_uid)
            ON DELETE CASCADE
);

CREATE TABLE payments
(
    order_uid      VARCHAR(50)              NOT NULL,
    transaction_id VARCHAR(50)              NOT NULL,
    request_id     VARCHAR(50),
    currency       VARCHAR(10)              NOT NULL,
    provider       VARCHAR(50)              NOT NULL,
    amount         INTEGER                  NOT NULL,
    payment_dt     TIMESTAMP WITH TIME ZONE NOT NULL,
    bank           VARCHAR(50)              NOT NULL,
    delivery_cost  INTEGER                  NOT NULL,
    goods_total    INTEGER                  NOT NULL,
    custom_fee     INTEGER                  NOT NULL,
    PRIMARY KEY (order_uid),
    CONSTRAINT fk_orders_payment
        FOREIGN KEY (order_uid)
            REFERENCES orders (order_uid)
            ON DELETE CASCADE
);

CREATE TABLE items
(
    item_id      SERIAL PRIMARY KEY,
    order_uid    VARCHAR(50)  NOT NULL,
    chrt_id      INTEGER      NOT NULL,
    track_number VARCHAR(50)  NOT NULL,
    price        INTEGER      NOT NULL,
    rid          VARCHAR(50)  NOT NULL,
    name         VARCHAR(100) NOT NULL,
    sale         INTEGER      NOT NULL,
    size         VARCHAR(10)  NOT NULL,
    total_price  INTEGER      NOT NULL,
    nm_id        INTEGER      NOT NULL,
    brand        VARCHAR(100) NOT NULL,
    status       INTEGER      NOT NULL,
    CONSTRAINT fk_orders_items
        FOREIGN KEY (order_uid)
            REFERENCES orders (order_uid)
            ON DELETE CASCADE
);
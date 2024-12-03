CREATE TABLE orders
(
    order_uid          VARCHAR(50)               PRIMARY KEY,
    track_number       VARCHAR(50)              ,
    entry              VARCHAR(10)              ,
    locale             VARCHAR(5)               ,
    internal_signature VARCHAR(100),
    customer_id        VARCHAR(50)              ,
    delivery_service   VARCHAR(20)              ,
    shardkey           VARCHAR(10)              ,
    sm_id              INTEGER                  ,
    date_created       TIMESTAMP WITH TIME ZONE ,
    oof_shard          VARCHAR(10)              
);

CREATE TABLE deliveries
(
    order_uid VARCHAR(50)  ,
    name      VARCHAR(100) ,
    phone     VARCHAR(20)  ,
    zip       VARCHAR(20)  ,
    city      VARCHAR(100) ,
    address   VARCHAR(200) ,
    region    VARCHAR(100) ,
    email     VARCHAR(100) ,
    PRIMARY KEY (order_uid),
    CONSTRAINT fk_orders_delivery
        FOREIGN KEY (order_uid)
            REFERENCES orders (order_uid)
            ON DELETE CASCADE
);

CREATE TABLE payments
(
    order_uid      VARCHAR(50)              ,
    transaction_id VARCHAR(50)              ,
    request_id     VARCHAR(50),
    currency       VARCHAR(10)              ,
    provider       VARCHAR(50)              ,
    amount         INTEGER                  ,
    payment_dt     TIMESTAMP WITH TIME ZONE ,
    bank           VARCHAR(50)              ,
    delivery_cost  INTEGER                  ,
    goods_total    INTEGER                  ,
    custom_fee     INTEGER                  ,
    PRIMARY KEY (order_uid),
    CONSTRAINT fk_orders_payment
        FOREIGN KEY (order_uid)
            REFERENCES orders (order_uid)
            ON DELETE CASCADE
);

CREATE TABLE items
(
    item_id      SERIAL PRIMARY KEY,
    order_uid    VARCHAR(50)  ,
    chrt_id      INTEGER      ,
    track_number VARCHAR(50)  ,
    price        INTEGER      ,
    rid          VARCHAR(50)  ,
    name         VARCHAR(100) ,
    sale         INTEGER      ,
    size         VARCHAR(10)  ,
    total_price  INTEGER      ,
    nm_id        INTEGER      ,
    brand        VARCHAR(100) ,
    status       INTEGER      ,
    CONSTRAINT fk_orders_items
        FOREIGN KEY (order_uid)
            REFERENCES orders (order_uid)
            ON DELETE CASCADE
);
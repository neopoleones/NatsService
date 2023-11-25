CREATE TABLE wb_delivery (
    delivery_id SERIAL PRIMARY KEY,
    zip         INT NOT NULL,
    name        TEXT NOT NULL,
    email       TEXT NOT NULL,
    phone       VARCHAR(32) NOT NULL,
    region      VARCHAR(64) NOT NULL,
    city        VARCHAR(128) NOT NULL,
    address     VARCHAR(128) NOT NULL
);

CREATE TABLE wb_payment (
    payment_id    SERIAL PRIMARY KEY,
    request_id    TEXT,

    payment_dt    INT NOT NULL,
    amount        INT NOT NULL,
    delivery_cost INT NOT NULL,
    goods_total   INT NOT NULL,
    custom_fee    INT NOT NULL,

    bank          VARCHAR(16) NOT NULL,
    provider      VARCHAR(24) NOT NULL,
    currency      VARCHAR(8)  NOT NULL,
    transaction   VARCHAR(24) NOT NULL
);

CREATE TABLE wb_order
(
    order_id     SERIAL PRIMARY KEY,
    delivery_id  SERIAL,
    payment_id   SERIAL,
    shard_key    INT NOT NULL,
    sm_id        INT NOT NULL,
    oof_shard    INT NOT NULL,
    date_created TIMESTAMP NOT NULL,

    internal_signature  TEXT,
    raw_order_id        VARCHAR(32) NOT NULL UNIQUE,
    track_number        VARCHAR(24) NOT NULL,
    customer_id         VARCHAR(16) NOT NULL,
    delivery_service    VARCHAR(16) NOT NULL,
    entry               VARCHAR(8)  NOT NULL,
    locale              VARCHAR(3)  NOT NULL,

    FOREIGN KEY (delivery_id) REFERENCES wb_delivery (delivery_id) ON DELETE CASCADE,
    FOREIGN KEY (payment_id) REFERENCES wb_payment (payment_id) ON DELETE CASCADE
);

CREATE TABLE wb_item (
    item_id      SERIAL PRIMARY KEY,
    order_id     SERIAL,

    status       INT NOT NULL,
    chrt_id      INT NOT NULL,
    price        INT NOT NULL,
    sale         INT NOT NULL,
    total_price  INT NOT NULL,
    nm_id        INT NOT NULL,

    size         VARCHAR(4),
    rid          VARCHAR(24) NOT NULL,
    track_number VARCHAR(32) NOT NULL,

    name         TEXT NOT NULL ,
    brand        TEXT NOT NULL,

    FOREIGN KEY (order_id) REFERENCES wb_order (order_id)
);
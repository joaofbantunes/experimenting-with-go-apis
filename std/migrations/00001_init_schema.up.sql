CREATE TABLE dishes
(
    id          bigint GENERATED ALWAYS AS IDENTITY,
    external_id uuid NOT NULL,
    name        text NOT NULL,
    CONSTRAINT pk_dishes PRIMARY KEY (id)
);

CREATE UNIQUE INDEX ix_dishes_external_id ON dishes (external_id);

CREATE TABLE orders
(
    id            bigint GENERATED ALWAYS AS IDENTITY,
    external_id   uuid                     NOT NULL,
    status        integer                  NOT NULL,
    registered_at timestamp with time zone NOT NULL,
    CONSTRAINT pk_orders PRIMARY KEY (id)
);

CREATE UNIQUE INDEX ix_orders_external_id ON orders (external_id);

CREATE TABLE order_items
(
    order_id bigint   NOT NULL,
    dish_id  bigint   NOT NULL,
    quantity smallint NOT NULL,
    CONSTRAINT pk_order_items PRIMARY KEY (order_id, dish_id),
    CONSTRAINT fk_order_items_dishes_dish_id FOREIGN KEY (dish_id) REFERENCES dishes (id) ON DELETE CASCADE,
    CONSTRAINT fk_order_items_orders_order_id FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);

CREATE INDEX ix_order_items_dish_id ON order_items (dish_id);

CREATE TABLE outbox_messages
(
    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    type          varchar(128) not null,
    payload       bytea        not null,
    created_at    timestamp(6) not null,
    trace_context bytea        null
);
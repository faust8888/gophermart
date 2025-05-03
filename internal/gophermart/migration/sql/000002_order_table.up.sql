CREATE TABLE "order" (
    id SERIAL PRIMARY KEY,
    user_login VARCHAR(255) NOT NULL,
    order_number BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_order_number_index UNIQUE (order_number)
);
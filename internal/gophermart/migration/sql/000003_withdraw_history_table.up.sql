CREATE TABLE "withdraw_history" (
    id SERIAL PRIMARY KEY,
    user_login VARCHAR(255) NOT NULL,
    order_number BIGINT NOT NULL,
    sum NUMERIC(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
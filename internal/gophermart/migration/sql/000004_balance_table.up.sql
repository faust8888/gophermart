CREATE TABLE "balance" (
    id SERIAL PRIMARY KEY,
    user_login VARCHAR(255) NOT NULL,
    withdrawn_sum NUMERIC(10, 2) NOT NULL,
    current_sum NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
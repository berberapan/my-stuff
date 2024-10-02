CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users(email, password) VALUES (
    'test@example.com',
    '$2a$12$emzpVkBb8ykQwGjnetHlc.m1H3xoNMCWXwefO7.K7WdQp6Xdp6jYO' --bcrypt encryption for 'password'
)

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    accessories TEXT,
    place VARCHAR(255),
    additional_notes TEXT,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO users(email, password) VALUES (
    'test@example.com',
    '$2a$12$emzpVkBb8ykQwGjnetHlc.m1H3xoNMCWXwefO7.K7WdQp6Xdp6jYO' --bcrypt encryption for 'password'
);

INSERT INTO items (name, description, accessories, place, additional_notes, user_id)
VALUES (
    'Laptop',
    '15-inch screen',
    'Charger',
    'Office',
    'yay -Syu',
    1
);

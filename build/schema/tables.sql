

CREATE TABLE IF NOT EXISTS app_db.users (
    id UUID PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX ON app_db.users (username);

CREATE TABLE IF NOT EXISTS app_db.assets (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    symbol TEXT NOT NULL,
    amount NUMERIC NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

INSERT INTO app_db.users (username, password_hash) VALUES (
    "bob",
    "$2a$10$g6efNvur33Ya9lV1Fo3ClekXylbUgv2JEl.9.21vCFRPEoBopk6k."
);

INSERT INTO app_db.assets (user_id, symbol, amount) VALUES (
    (SELECT id FROM users WHERE username = "bob"),
    "EUR",
    10000.00
);

INSERT INTO app_db.assets (user_id, symbol, amount) VALUES (
    (SELECT id FROM users WHERE username = "bob"),
    "USD",
    10000.00
);

INSERT INTO app_db.users (username, password_hash) VALUES (
    "tracy",
    "$2a$10$rK2Ia/SRslB0c7GTTsHKqewhFoccS/Q/189UHCkiT6HIMKg.lHgHi"
);

INSERT INTO app_db.assets (user_id, symbol, amount) VALUES (
    (SELECT id FROM users WHERE username = "tracy"),
    "EUR",
    10000.00
);

INSERT INTO app_db.assets (user_id, symbol, amount) VALUES (
    (SELECT id FROM users WHERE username = "tracy"),
    "USD",
    10000.00
);


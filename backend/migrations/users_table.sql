CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    bing_email TEXT UNIQUE NOT NULL,
    profile_picture TEXT DEFAULT '',
    joined_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    communities JSONB DEFAULT '[]',
    friends JSONB DEFAULT '[]',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(bing_email);
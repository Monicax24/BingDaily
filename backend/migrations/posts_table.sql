CREATE TABLE IF NOT EXISTS posts (
    post_id SERIAL PRIMARY KEY,
    community_id INTEGER REFERENCES communities(community_id) ON DELETE CASCADE,
    picture TEXT NOT NULL,
    caption TEXT,
    author INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    time_posted TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    likes JSONB DEFAULT '[]'
);
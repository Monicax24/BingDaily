CREATE TABLE IF NOT EXISTS communities (
    community_id SERIAL PRIMARY KEY,
    picture TEXT DEFAULT '',
    description TEXT,
    members JSONB DEFAULT '[]',
    moderators JSONB DEFAULT '[]',
    posts JSONB DEFAULT '[]',
    post_time TIME DEFAULT '09:00:00',
    default_prompt TEXT DEFAULT 'What did you do today?'
);
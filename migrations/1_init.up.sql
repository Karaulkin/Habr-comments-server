CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username TEXT UNIQUE NOT NULL
);

CREATE TABLE posts (
                       id SERIAL PRIMARY KEY,
                       author_id INT REFERENCES users(id) ON DELETE CASCADE,
                       title TEXT NOT NULL,
                       content TEXT NOT NULL,
                       allow_comments BOOLEAN DEFAULT TRUE,
                       created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE comments (
                          id SERIAL PRIMARY KEY,
                          post_id INT REFERENCES posts(id) ON DELETE CASCADE,
                          author_id INT REFERENCES users(id) ON DELETE CASCADE,
                          parent_id INT REFERENCES comments(id) ON DELETE CASCADE, -- для вложенности
                          content TEXT CHECK (LENGTH(content) <= 2000) NOT NULL,
                          created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_comments_post ON comments(post_id);

CREATE INDEX idx_comments_parent ON comments(parent_id);
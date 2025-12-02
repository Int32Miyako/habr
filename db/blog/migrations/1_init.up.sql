CREATE TABLE IF NOT EXISTS blogs (
    id SERIAL PRIMARY KEY,
    name TEXT);

CREATE INDEX IF NOT EXISTS idx_blog_name ON blogs(name);

-- +goose Up 
CREATE TABLE posts(
    id UUID Primary Key,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    body TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE  posts;

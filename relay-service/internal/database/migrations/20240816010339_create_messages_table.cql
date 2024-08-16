create table if not exists messages (
    id BIGINT primary key,
    created_at timestamp,
    message text,
    senderId text,
    receiverId text
);

CREATE INDEX ON messages (receiverId);
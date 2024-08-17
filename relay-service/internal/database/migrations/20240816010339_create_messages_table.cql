create table if not exists messages (
    id UUID primary key,
    created_at timestamp,
    message text,
    sender_id text,
    receiver_id text
);

CREATE INDEX ON messages (receiver_id);
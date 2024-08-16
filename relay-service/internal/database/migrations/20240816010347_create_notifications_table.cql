create table if not exists notifications (
    id BIGINT primary key,
    created_at timestamp,
    messageId BIGINT
)
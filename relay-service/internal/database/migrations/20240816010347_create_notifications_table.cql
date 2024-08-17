create table if not exists notifications (
    id UUID primary key,
    created_at timestamp,
    message_id UUID
)
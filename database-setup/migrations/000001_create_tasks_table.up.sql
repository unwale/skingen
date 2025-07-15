CREATE TABLE IF NOT EXISTS task_service.tasks (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,

    prompt     TEXT NOT NULL,
    status     VARCHAR(20) NOT NULL DEFAULT 'pending',
    result_url TEXT
);

CREATE INDEX idx_tasks_deleted_at ON task_service.tasks(deleted_at);
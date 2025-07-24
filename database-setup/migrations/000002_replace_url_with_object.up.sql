ALTER TABLE task_service.tasks
    DROP COLUMN IF EXISTS result_url,
    ADD COLUMN IF NOT EXISTS object_id TEXT;
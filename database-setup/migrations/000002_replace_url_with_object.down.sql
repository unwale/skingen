ALTER TABLE task_service.tasks
    DROP COLUMN IF EXISTS object_id,
    ADD COLUMN IF NOT EXISTS result_url TEXT;
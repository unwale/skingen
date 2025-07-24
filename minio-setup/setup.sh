#!/bin/sh

set -e

until mc alias set myminio http://minio:9000 ${MINIO_ROOT_USER} ${MINIO_ROOT_PASSWORD}; do
    echo "Waiting for MinIO server..."
    sleep 2
done

echo "MinIO server is ready. Starting setup."

echo "Creating bucket 'skin-images'..."
mc mb myminio/${SKIN_IMAGES_BUCKET} --ignore-existing

echo "Creating policies..."
mc admin policy create myminio readonly-policy /app/policies/readonly.json
mc admin policy create myminio writeonly-policy /app/policies/writeonly.json

echo "Creating users..."
mc admin user add myminio ${MINIO_TASK_SERVICE_USER} ${MINIO_TASK_SERVICE_PASSWORD}
mc admin user add myminio ${MINIO_WORKER_USER} ${MINIO_WORKER_PASSWORD}

echo "Attaching policies to users..."
mc admin policy attach myminio readonly-policy --user ${MINIO_TASK_SERVICE_USER}
mc admin policy attach myminio writeonly-policy --user ${MINIO_WORKER_USER}

echo "Setup complete! Users and policies are configured."

exit 0
FROM python:3.11-slim

COPY ./services/model-server/requirements.txt /app/requirements.txt
WORKDIR /app
RUN pip install --no-cache-dir -r requirements.txt

COPY ./services/model-server /app

EXPOSE ${MODEL_SERVER_PORT}

CMD ["python", "main.py"]
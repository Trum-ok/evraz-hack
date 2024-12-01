FROM python:3.12.5-alpine:3.20

WORKDIR /app
COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY . .
RUN pip install python-dotenv

CMD ["python", "run.py"]

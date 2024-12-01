FROM python:3.12.5-slim

WORKDIR /app
COPY requirements.txt .

RUN pip cache purge
RUN pip install --upgrade pip
RUN pip install --no-cache-dir -r requirements.txt

COPY . .
RUN pip install python-dotenv

CMD ["python", "run.py"]

FROM python:3.12.5-alpine

WORKDIR /app
COPY requirements.txt .

RUN pip install --upgrade pip
RUN pip install -r requirements.txt

COPY . .
RUN pip install python-dotenv

CMD ["python", "run.py"]

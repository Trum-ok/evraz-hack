import os
os.environ["TRANSFORMERS_NO_PYTORCH_WARNING"] = "1"

import requests
from decouple import config
from transformers import AutoTokenizer

API_KEY=config('API_KEY')
URL = config("URL")
HEADERS = {
    "Authorization": API_KEY,
    "Content-Type": "application/json"
}


def count_tokens(prompt: str) -> int:
    tokenizer = AutoTokenizer.from_pretrained("mistralai/Mistral-Nemo-Instruct-2407")
    tokens = tokenizer.encode(prompt)
    return len(tokens)


def check_for_file(prompt: str) -> bool:
    tokens = count_tokens(prompt)
    return (tokens < 6144 - 1024 - 200)


def check_for_archive(prompt: str) -> bool:
    tokens = count_tokens(prompt)
    return (tokens <= 6144 - 1024 + 50)


def send_request(prompt: str) -> requests.Response:
    """
    Отправка запроса к LLM по API

    :args:
    - prompt `str`: код пользователя 
    
    :returns:
    - `requests.Response`: ответ от LLM
    """
    data = {
        "model": "mistral-nemo-instruct-2407",
        "messages": [
            {"role": "system", "content": "Язык: русский(Россия)"},
            # {"role": "system", "content": ""},
            {"role": "system", "content": "Ты пишешь code-review отчеты по C# проектам на русском языке в MarkDown-файлы"},
            {"role": "user", "content": prompt}
        ],
        "max_tokens": 1024,
        "temperature": 0.3
    }
    return requests.post(url=URL, headers=HEADERS, json=data)


def finally_generation(text: str, date: str, project_name: str) -> requests.Response:
    data = {
        "model": "mistral-nemo-instruct-2407",
        "messages": [
            {"role": "system", "content": "Язык: Русский(Россия)"},
            {"role": "system", "content": f"Сегоднящняя дата: {date}. Название проекта: {project_name}"},
            {"role": "system", "content": "Ты очень хорошо исправляешь MARKDOWN-файлы с code-review по C#. Мне не нужно указание ЯП, автора, и обзоршика."},
            {"role": "user", "content": f"Вот CodeReview.md: \n{text}\n\nТвоя задача вернуть мне оформленный, красивый MARKDOWN-файл с заголовком формата 'Code-review проекта `название_проекта` от `сегодняшняя_дата`'. Постарайся пожалуйста!"}
        ],
        "max_tokens": 1024,
        "temperature": 0.1
    }
    return requests.post(url=URL, headers=HEADERS, json=data)


if __name__ == "__main__":
    prompt = ""
    res = send_request(prompt)
    print(res)

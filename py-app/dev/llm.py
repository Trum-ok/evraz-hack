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


def check_tokens(prompt: str) -> bool:
    tokens = count_tokens(prompt)
    return (tokens <= 6144 - 1024 - 50)


def send_request(sys_prompt: str, prompt: str, file: bool = False) -> requests.Response:
    """
    Отправка запроса к LLM по API

    :args:
    - prompt `str`: код пользователя 
    
    :returns:
    - `requests.Response`: ответ от LLM
    """
    while not check_tokens(prompt):  # проверка количества токенов
        prompt = prompt[:-100]    
    
    data = {
        "model": "mistral-nemo-instruct-2407",
        "messages": [
            {"role": "system", "content": "Язык: русский(Россия)"},
            {"role": "system", "content": sys_prompt},
            {"role": "user", "content": prompt}
        ],
        "max_tokens": 1024,
        "temperature": 0.3
    }
    return requests.post(url=URL, headers=HEADERS, json=data)


def finally_generation(text: str, date: str, project_name: str) -> requests.Response:
    t = f"Твоя задача вернуть мне оформленный, красивый MARKDOWN-файл с заголовком формата 'Code-review проекта `название_проекта` от `сегодняшняя_дата`'. Постарайся пожалуйста! Вот CodeReview.md: \n\n{text}"
    
    while not check_tokens(t):
        if len(text) <= 100:
            raise ValueError("Текст слишком короткий для генерации.")
        text = text[:-100]
        t = f"Твоя задача вернуть мне оформленный, красивый MARKDOWN-файл с заголовком формата 'Code-review проекта `название_проекта` от `сегодняшняя_дата`'. Постарайся пожалуйста! Вот CodeReview.md: \n\n{text}"

    data = {
        "model": "mistral-nemo-instruct-2407",
        "messages": [
            {"role": "system", "content": "Язык: Русский(Россия)"},
            {"role": "system", "content": f"Сегоднящняя дата: {date}. Название проекта: {project_name}"},
            {"role": "system", "content": "Ты очень хорошо исправляешь MARKDOWN-файлы с code-review по. Мне НЕ НУЖНО указание ЯП, автора, и обзоршика."},
            {"role": "user", "content": t}
        ],
        "max_tokens": 1024,
        "temperature": 0.2
    }
    return requests.post(url=URL, headers=HEADERS, json=data)


if __name__ == "__main__":
    prompt = ""
    res = send_request(prompt)
    print(res)

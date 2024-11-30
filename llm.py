import requests
from decouple import config

API_KEY=config('API_KEY')
URL = config("URL")
print(API_KEY, "|", URL)


HEADERS = {
    "Authorization": API_KEY,
    "Content-Type": "application/json"
}


def send_request():
    data = {
        "model": "mistral-nemo-instruct-2407",
        "messages": [
            {"role": "system", "content": "Ты пишешь code-review отчеты по C# проектам на русском языке в MarkDown-файлы"},
            {"role": "user", "content": "Log.Error(ex, ex.Message);"}
        ],
        "max_tokens": 1024,
        "temperature": 0.3
    }

    r = requests.post(url=URL, headers=HEADERS, json=data)
    return r.json(), r.status_code

    # # Обработка ответа
    # if response.status_code == 200:
    #     print("Ответ модели:")
    #     print(response.json())
    # else:
    #     print(f"Ошибка: {response.status_code}")
    #     print(response.text)


def finally_generation(text: str) -> requests.Response:
    data = {
        "model": "mistral-nemo-instruct-2407",
        "messages": [
            {"role": "system", "content": "Язык: Русский(Россия)"},
            {"role": "system", "content": "Ты очень хорошо исправляешь MARKDOWN-файлы с code-review по C#. Мне не нужны заголовки 'code review', указание языка, автора, даты, обзоршика."},
            {"role": "user", "content": f"Вот CodeReview.md: \n{text}\n\nТвоя задача вернуть мне оформленный, красивый MARKDOWN-файл без заголовок, дат и авторов. Постарайся пожалуйста!"}
        ],
        "max_tokens": 1024,
        "temperature": 0.1
    }
    return requests.post(url=URL, headers=HEADERS, json=data)


res, status_code = send_request()
# print(res['choices'][0]['message']['content'])
if res and status_code == 200:
    print(finally_generation(res['choices'][0]['message']['content']).json()['choices'][0]['message']['content'])

# ХАКАТОН ЕВРАЗА 3.0
![евраз](/_assets/evraz.jpg)

## Инструмент для code-review C#\Python\TypeScript кода

- [Телеграм-бот](https://t.me/EVRAZ_CR_bot)
- [Веб-сайт]() *(in dev)* 
- [Сайт хакатона](https://xn--80aaaairqt2ajzt9a.xn--p1ai/)

## Как это работает?
Мы разработали **несколько пайплайнов** анализа кода (для C#, Python и TypeScript) на основе руководств по code review компании Eвраз. При помощи RAG-алгоритмов LLM может проводить code review так, как это делают в компании, только в несколько раз быстрее, удобнее и дешевле.

> [!WARNING]  
> Не тестируйте одновременно тг-бота и веб-сайт! \
> Если нарушить это правило, то **что-то ляжет** из-за огранниченого количества клиентов к API LLM :c \
> Не убивайте горутины и ветки процессов 🥺

## Локальный запуск
```bash
git clone https://github.com/Trum-ok/evraz-hack
cd evraz-hack
pip install -r requirements.txt
python main.py 
```
или

```bash
git clone https://github.com/Trum-ok/evraz-hack
cd evraz-hack
docker build -t evraz-cr-bot .
docker run --env-file .env -p 8000:8000 evraz-cr-bot
```

## Лицензия
Кажется это проект под MIT лицензией *(а что это)*

## Авторы
Создано командой Invalid Syntax с большой любовью и огромными усилиями 💗

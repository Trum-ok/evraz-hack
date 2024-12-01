# ХАКАТОН ЕВРАЗА 3.0
![евраз](/_assets/evraz.jpg)

## Инструмент для code-review C#\Python\JS\TS\TSX\Go-lang кода

- **[Телеграм-бот](https://t.me/EVRAZ_CR_bot)**
- **[Веб-сайт](https://evrasa-dq2203.amvera.io/)**
- **[Сайт хакатона](https://xn--80aaaairqt2ajzt9a.xn--p1ai/)**
- **[Презентация](https://www.figma.com/slides/GzNeGO8u0tvLKLr9V58SVf/%D0%95%D0%B2%D1%80%D0%B0%D0%B7?node-id=1-66&t=eREjkzmXZVj3xbhI-1)**

## Проблематика
В крупных компаниях для разработки элементов различных систем обычно привлекают разработчиков разного уровня и специализации, а также подрядные организации. Несмотря на наличие единых правил оформления, стиля, наименования и структуры, участие большого числа людей существенно увеличивает время, необходимое для проверки и исправления создаваемых merge requests. Часто на внесение исправлений уходит даже больше времени, чем на написание самого кода в этом MR. \
*(дубликат с презентации)*

## Как это работает?
Мы разработали **несколько пайплайнов** анализа кода *(для C#, Python, TypeScript/JavaScript & React)* на основе руководств по code review компании Eвраз. При помощи RAG-алгоритмов LLM может проводить code review так, как это делают в компании, только в несколько раз быстрее, удобнее и дешевле. Для анализа файловой структуры проекта используется **[prft](https://github.com/Trum-ok/project-file-tree)** *(наша библиотека)*.

> [!WARNING]  
> Не тестируйте одновременно тг-бота и веб-сайт! \
> Если нарушить это правило, то **что-то ляжет** из-за огранниченого количества клиентов к API LLM у одного ключа **:c**    
> Не убивайте горутины и ветки процессов 🥺 

> [!WARNING]  
> Если вы видите индикатор загрузки на сайте, то все хорошо, если пропадает и появляется черный экран смерти, то перезагрузите страницу.

К сожалению за столь короткое время очень тяжело реализовать сильный RAG + есть некий нехваток данных для создания баз знаний под каждый ЯП, но в данном проекте отражена сама суть этой идеи - 
**дать модели контекст (мы преобразовали ваши pdf для код-ревьюверов и развернули их в мини-базу-знаний), максимально приближенный к контексту код-ревьюверов**: опыт (кол-во параметров и данные, на которых обучали модель, может сильно влиять), плохие хорошие практики в коде, актуальные стандарты, *etc.*

> [!TIP]
> **Контекст - наше все!**


## Локальный запуск
### py-app
```bash
git clone https://github.com/Trum-ok/evraz-hack
cd evraz-hack/py-app
pip install -r requirements.txt
python main.py 
```
или

```bash
git clone https://github.com/Trum-ok/evraz-hack
cd evraz-hack/py-app
docker build -t evraz-cr-bot .
docker run --env-file .env evraz-cr-bot
```
### go-app
```bash
git clone https://github.com/Trum-ok/evraz-hack
cd evraz-hack
cd go-app
cd app
go run main.go
```
или
```bash
git clone https://github.com/Trum-ok/evraz-hack
cd evraz-hack
cd go-app
docker-compose up --build
```

## Файловая структура проекта
```bash
 ├── go-app     - веб сайт
 │   ├── app
 │   │   ├── handlers
 │   │   │   ├── download.go
 │   │   │   ├── error.go
 │   │   │   └── upload.go
 │   │   ├── pkg
 │   │   │   ├── file.go
 │   │   │   ├── llm.go
 │   │   │   ├── pdf.go
 │   │   │   └── zip.go
 │   │   ├── servers
 │   │   │   └── sever.go
 │   │   ├── static
 │   │   │   └── шрифты
 │   │   ├── templates
 │   │   │   └── лэндинги
 │   │   ├── uploads
 │   │   │   └── icon.jpg
 │   │   └── main.go
 │   ├── Dockerfile
 │   ├── README.md
 │   ├── docker-compose.yml
 │   ├── example1.md
 │   ├── go.mod
 │   └── go.sum
 ├── py-app     - телеграм бот
 │   ├── dev
 │   │   ├── __init__.py
 │   │   ├── llm.py
 │   │   ├── process.py
 │   │   ├── result.py
 │   │   └── text.py
 │   ├── prompts
 │   │   └── файлы с промптами
 │   ├── .dockerignore
 │   ├── Dockerfile
 │   ├── requirements.txt
 │   └── run.py
 ├── .gitignore
 ├── LICENSE
 └── README.md
```

## Лицензия
Кажется это проект под MIT лицензией *(а что это)*

## Авторы
Создано командой Invalid Syntax с большой любовью и огромными усилиями 💗

import time
from datetime import datetime

from decouple import config
from telebot import TeleBot
from threading import Thread
from process import process_archive, process_file

TOKEN = config('TELEGRAM_TOKEN')
bot = TeleBot(TOKEN)


def update_status(chat_id, message_id, initial_text, while_processing_flag):
    seconds_elapsed = 0
    while while_processing_flag():
        time.sleep(1)
        seconds_elapsed += 1
        try:
            bot.edit_message_text(
                chat_id=chat_id,
                message_id=message_id,
                text=f"{initial_text} ({seconds_elapsed} секунд)"
            )
        except Exception as e:
            print(f"Ошибка обновления сообщения: {e}")


@bot.message_handler(content_types=['document'])
def handle_document(message):
    file_info = bot.get_file(message.document.file_id)
    downloaded_file = bot.download_file(file_info.file_path)

    if message.document.file_name.endswith('.zip'):
        result_report = process_archive(downloaded_file)
        r_type = "архив"
    else:
        result_report = process_file(downloaded_file)
        r_type = "файл"

    processing_message = bot.reply_to(message, f"Ваш {r_type} обрабатывается...")
    chat_id = message.chat.id
    message_id = processing_message.message_id
    processing_flag = [True]

    Thread(target=update_status, args=(chat_id, message_id, f"Ваш {r_type} обрабатывается", lambda: processing_flag[0])).start()
    result_report = process_archive(downloaded_file) if r_type == "архив" else process_file(downloaded_file)
    processing_flag[0] = False

    bot.edit_message_text(
        chat_id=chat_id,
        message_id=message_id,
        text="Формируем ответ... (0 секунд)"
    )

    forming_flag = [True]

    Thread(target=update_status, args=(chat_id, message_id, "Формируем ответ", lambda: forming_flag[0])).start()

    # Формируем ответ (симуляция)
    time.sleep(3)  # Замените на реальную логику формирования ответа
    forming_flag[0] = False

    # Обновляем сообщение на финальное и отправляем файл
    bot.edit_message_text(
        chat_id=chat_id,
        message_id=message_id,
        text=f"Ваш {r_type} был обработан, результаты прикреплены к сообщению.\nCode Review от {datetime.now()}"
    )
    with open(result_report, "rb") as report_file:
        bot.send_document(chat_id=chat_id, document=report_file)


@bot.message_handler(commands=['start'])
def start_message(message):
    bot.reply_to(message, "Привет! Я бот для проверки проектов. Отправьте мне файл или архив для обработки.")


@bot.message_handler(content_types=["sticker"])
def handle_sticker(message):
    bot.reply_to(message, "Классный стикер!")


@bot.message_handler(content_types=["photo"])
def handle_photo(message):
    bot.reply_to(message, "Классная картинка!")


@bot.message_handler(func=lambda message: True)
def unknown_command(message):
    bot.reply_to(message, "Я не знаю, что делать с этим. Пожалуйста, отправьте мне файл или архив для обработки.")


if __name__ == '__main__':
    print("Bot started")
    bot.infinity_polling()

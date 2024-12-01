import os
import sys
import json
import asyncio
import logging
import random

from aiogram import F
from aiogram import Bot, Dispatcher
from aiogram.client.default import DefaultBotProperties
from aiogram.enums import ParseMode
from aiogram.filters import CommandStart
from aiogram.types import Message, InputMediaDocument, FSInputFile
from decouple import config

from dev.text import CANT, BLESS, START
from dev.process import process

TOKEN = config("TELEGRAM_TOKEN")

dp = Dispatcher()


@dp.message(CommandStart())
async def start_message(message: Message) -> None:
    await message.answer(START)


@dp.message(F.document)
async def document_message(message: Message) -> None:
    sec = 0
    form_sec = 0
    process_flag = True
    formating_flag = False

    msg = await message.answer(f"Обрабатываю... ({sec} сек.)")
    doc_type =message.document.file_name.split('.')[-1]
    if doc_type not in ['zip', 'cs', 'py', 'ts', 'js']:
        await msg.edit_text(f"Формат *{doc_type}* не поддерживается. Пожалуйста, поробуйте другой файл")
        return

    async def increment_timer(timer_type):
        nonlocal sec, form_sec
        while process_flag or formating_flag:
            await asyncio.sleep(1)
            if timer_type == "processing" and process_flag:
                sec += 1
                await msg.edit_text(f"{"🏃‍♂️" if sec % 2 else "🏃‍♂️‍➡️"} Обрабатываю... ({sec} сек.)")
            elif timer_type == "formatting" and formating_flag:
                form_sec += 1
                await msg.edit_text(f"Формирую ответ... ({form_sec} сек.)")
            # if sec >= 60:
            #     await msg.edit_text("При обработке файла произошла ошибка, попробуйте позже...")
            if form_sec >= 30:
                await msg.edit_text("При форматировании файла произошла ошибка, попробуйте позже...")

    process_task = asyncio.create_task(increment_timer("processing"))
    file_id = message.document.file_id
    file = await bot.get_file(file_id)
    path = rf"./_downloads/{file_id[1:10]}"
    await bot.download_file(file.file_path, path)
    fp_output = "./_results/" + "".join([str(random.randrange(1, 100)) for _ in range(10)])+'.md'


    try:
        d = await asyncio.to_thread(process, path, fp_output)  # Выполняем процесс в отдельном потоке
    finally:
        process_flag = False  # Устанавливаем флаг завершения
        formating_flag = False
        process_task.cancel()  # Завершаем таймер
        try:
            await process_task  # Ожидаем завершение `increment_timer`
        except asyncio.CancelledError:
            pass

    d = d if isinstance(d, str) else json.dumps(d)
    if 'psycopg2.OperationalError' in d:
        await msg.edit_text("Слишком много запросов к API, поробуйте позже...")
        return

    t = json.loads(d if isinstance(d, str) else json.dumps(d))
    print(json.loads(d if isinstance(d, str) else json.dumps(d)))
    process_flag = False

    # formating_flag = True
    # msg = await msg.edit_text(f"Формирую ответ... ({form_sec} сек.)")
    # formatting_task = asyncio.create_task(increment_timer("formatting"))
    # await asyncio.sleep(2.5) # логика форматинга файла
    # formating_flag = False  # Завершение форматирования
    # await formatting_task

    if t['path']:
        file_to_attach = FSInputFile(t['path'], filename="CodeReview.md")
        media = InputMediaDocument(media=file_to_attach, caption="**Code review готов! 🔥**")

        await msg.edit_media(media)
        await message.answer_sticker(random.choice(BLESS))

        os.remove(t['path'])
    else:
        await msg.edit_text("Произошла ошибка, попробуйте еще раз")


@dp.message(F.sticker)
async def handle_sticker(message: Message) -> None:
    await message.answer("Классный стикер!")


@dp.message(F.photo)
async def handle_photo(message: Message) -> None:
    await message.answer("Классная картинка!")


@dp.message(F.text)
async def unknown_command(message: Message) -> None:
    await message.answer(CANT)


async def main() -> None:
    global bot
    by_default = DefaultBotProperties(parse_mode=ParseMode.MARKDOWN)
    bot = Bot(token=TOKEN, default=by_default)
    await dp.start_polling(bot)


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO, stream=sys.stdout)
    asyncio.run(main())

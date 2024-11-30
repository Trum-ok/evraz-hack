import os
from typing import Optional
from datetime import datetime


def find_latest_modification(directory) -> Optional[datetime]:
    """
    Ищет время последнего редактирования
    
    :args:
    - directory - путь к директории

    :returns:
    - latest_date `datetime` - время последнего редактирования

    """
    latest_time = None

    for root, dirs, files in os.walk(directory):
        for file in files:
            file_path = os.path.join(root, file)
            modification_time = os.path.getmtime(file_path)  # время последнего изменения файла

            if latest_time is None or modification_time > latest_time:
                latest_time = modification_time

    if latest_time:
        latest_date = datetime.fromtimestamp(latest_time)
        return latest_date
    else:
        return None


def write_to_md_file(file_name: str, text: str) -> None:
    with open(file_name, "w", encoding="utf-8") as file:
        file.write(text)


if __name__ == "__main__":
    directory = "D:\\OKBOOMER\\3агрузки\\Готовые проекты"
    latest_date = find_latest_modification(directory)

    if latest_date:
        print(f"Самое позднее время редактирования: {latest_date}")
    else:
        print("Файлы не найдены.")

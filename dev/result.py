import os
import zipfile
from typing import Optional, Union
from datetime import datetime


def find_latest_modification(path: Union[str, os.PathLike]) -> Optional[datetime]:
    """
    Ищет время последнего редактирования файла/директории или внутри архива ZIP.

    :args:
    - path: путь к директории, файлу или ZIP-архиву

    :returns:
    - latest_date: `datetime` - время последнего редактирования, если найдено, иначе None
    """
    latest_time = None

    # Если путь указывает на файл
    if os.path.isfile(path):
        if zipfile.is_zipfile(path):  # Если это ZIP-архив
            with zipfile.ZipFile(path, 'r') as zip_file:
                for file_info in zip_file.infolist():
                    # Конвертируем время из ZIP (формат: (YYYY, MM, DD, HH, MM, SS))
                    modification_time = datetime(*file_info.date_time).timestamp()
                    if latest_time is None or modification_time > latest_time:
                        latest_time = modification_time
        else:
            latest_time = os.path.getmtime(path)

    elif os.path.isdir(path):
        for root, dirs, files in os.walk(path):
            for file in files:
                file_path = os.path.join(root, file)
                modification_time = os.path.getmtime(file_path)
                if latest_time is None or modification_time > latest_time:
                    latest_time = modification_time

    if latest_time:
        latest_date = datetime.fromtimestamp(latest_time)
        return latest_date
    return None


if __name__ == "__main__":
    directory = "D:\\OKBOOMER\\3агрузки\\Готовые проекты"
    latest_date = find_latest_modification(directory)

    if latest_date:
        print(f"Самое позднее время редактирования: {latest_date}")
    else:
        print("Файлы не найдены.")

import io
from zipfile import ZipFile


def process_file(input_fp, output_fp):
    """_summary_

    Args:
        input_fp (str): путь к исходному файлу
        output_fp (str): путь к конечному файлу
    """
    # 1. Считать содержимое из .txt
    with open(input_fp, "r") as txt_file:
        content = txt_file.read()

    print("Содержимое файла .txt:", content)

    # 2. Очистить содержимое .txt
    with open(input_fp, "w") as txt_file:
        txt_file.write("")  # Записываем пустую строку для очистки

    print(f"Файл {input_fp} очищен.")

    # 3. Изменить содержимое (например, добавим новый текст)
    modified_content = content.upper() + "\nДобавлено: Новый текст!"

    with open(output_fp, "w") as out_file:
        out_file.write(modified_content)

    print(f"Измененное содержимое записано в {output_fp}.")


def create_report(report_path, contents):
    # Для создания файла репорта, можно править для ваших потребностей
    with open(report_path, "w") as file:
        file.write(contents)
    return report_path


# # Функция для обработки файлов и создания репортов
# def process_file(file) -> str:
#     # Здесь должна быть логика обработки файла
#     print("Processing file:", file)
#     report = create_report("report.txt", "Hello world")
#     return report


# Функция для обработки архивов
def process_archive(zip_file):
    with ZipFile(io.BytesIO(zip_file), 'r') as archive:
        for file in archive.namelist():
            with archive.open(file) as nested_file:
                file_contents = nested_file.readlines()
                # Здесь должна быть логика обработки архива

    # report = create_report("report.txt", file_contents)
    # return report
    return "null.txt"

# file: io.FileIO = None
# file.write(f'## Анализ проекта {name} от {datetime.now}')
# file.write(f"Дата последнего изменения проекта: {last_redaction}")



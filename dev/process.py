import io
from zipfile import ZipFile
from dev.llm import finally_generation, send_request
from dev.result import find_latest_modification


def process_file(input_fp: str, output_fp: str):
    """
    Обработка файла

    :args:
    - input_fp `str`: путь к исходному файлу
    - output_fp `str`: путь к конечному файлу
    """
    with open(input_fp, "r") as txt_file:
        content = txt_file.read()
    response = send_request(content)

    if response.status_code != 200:
        return {"error (g)": response.text}

    # project_name = input_fp.split("/")[-1].split(".")[0]
    # date = find_latest_modification(input_fp)
    # response = finally_generation(response.text, date, project_name)
    # response.encoding = 'utf-8'
    if response.status_code != 200:
        return {"error (f)": response.text}
    with open(output_fp, 'a') as md:
        md.write(response.text)
    return {"path": output_fp}


def process_archive(input_fp: str, output_fp: str) -> dict[str, str]:
    """
    Обработка архива

    :args:
    - input_fp (str): путь к исходному файлу
    - output_fp (str): путь к конечному файлу

    :returns:
    - error/path
    """
    with ZipFile(io.BytesIO(input_fp), 'r') as archive:
        for file in archive.namelist():
            with archive.open(file) as nested_file:
                file_contents = nested_file.readlines()
                # Здесь должна быть логика обработки архива


    project_name = input_fp.split("/")[-1].split(".")[0]
    date = find_latest_modification(input_fp)
    response = finally_generation(file_contents, date, project_name)
    if response.status_code != 200:
        return {"error": response.text}
    with open('output_fp', 'a') as md:
        md.write(response.text)
    return {"path": output_fp}


def process(input_fp: str, output_fp: str) -> dict[str, str]:
    if input_fp.endswith(".zip"):
        return process_archive(input_fp, output_fp)
    else:
        return process_file(input_fp, output_fp)


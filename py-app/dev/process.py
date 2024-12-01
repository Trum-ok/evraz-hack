import re
import json

from enum import Enum
from zipfile import ZipFile
from dev.llm import finally_generation, send_request
from dev.result import find_latest_modification
from prompts import py_prompt, cs_prompt, react_prompt, tsjs_prompt, go_prompt


class SysPrompts(Enum):
    cs = "Ты пишешь code-review отчеты по C# проектам на русском языке в MarkDown-файлы"
    py = "Ты пишешь code-review отчеты по Python проектам на русском языке в MarkDown-файлы"
    ts = "Ты пишешь code-review отчеты по TypeScript проектам на русском языке в MarkDown-файлы"
    tsx = "Ты пишешь code-review отчеты по TypeScript (React) проектам на русском языке в MarkDown-файлы"
    js = "Ты пишешь code-review отчеты по JavaScript проектам на русском языке в MarkDown-файлы"
    go = "Ты пишешь code-review отчеты по Go-lang проектам на русском языке в MarkDown-файлы"


class Prompts(Enum):
    cs = cs_prompt
    py = py_prompt
    ts = tsjs_prompt
    tsx = react_prompt
    js = tsjs_prompt
    go = go_prompt


def get_sysprompt(file_extension: str) -> str:
    """Получить системный промпт на основе расширения файла."""
    return SysPrompts[file_extension].value if file_extension in SysPrompts.__members__ else ""


def get_prompt(file_extension: str) -> str:
    """Получить промпт на основе расширения файла."""
    return Prompts[file_extension].value if file_extension in Prompts.__members__ else ""


def get_file_struct() -> str:
    pass


def process_file(input_fp: str, output_fp: str):
    """
    Обработка файла

    :args:
    - input_fp `str`: путь к исходному файлу
    - output_fp `str`: путь к конечному файлу
    """
    file_extension = input_fp.split(".")[-1]
    sys_prompt = get_sysprompt(file_extension)
    prompt = get_prompt(file_extension).format(files=f"└── {input_fp.split("/")[-1]}")

    with open(input_fp, "r", encoding='utf-8') as txt_file:
        content = txt_file.read()
    response = send_request(sys_prompt, prompt+content)

    if response.status_code != 200:
        return {"error (g)": response.text.encode().decode('unicode_escape')}

    project_name = input_fp.split("/")[-1].split(".")[0]
    date = find_latest_modification(input_fp)
    response = finally_generation(json.loads(response.text)["choices"][0]["message"]["content"], date, project_name)
    if response.status_code != 200:
        return {"error (f)": response.text.encode().decode('unicode_escape')}
    with open(output_fp, 'a', encoding='utf-8') as md:
        rt = json.loads(response.text)
        md.write(rt["choices"][0]["message"]["content"])
    return {"path": output_fp}


def process_archive(input_fp: str, output_fp: str) -> dict[str, str]:
    """
    Обработка архива

    :args:
    - input_fp (str): путь к архиву
    - output_fp (str): путь к конечному файлу

    :returns:
    - error/path (dict): результат выполнения
    """
    with open(input_fp, 'rb') as archive_file:
        with ZipFile(archive_file, 'r') as archive:
            for file in archive.namelist():
                file_extension = file.split(".")[-1]
                sys_prompt = get_sysprompt(file_extension)
                prompt = get_prompt(file_extension).format(files=" ")

                if not sys_prompt:
                    continue

                with archive.open(file) as nested_file:
                    file_contents = nested_file.read().decode('utf-8')

                response = send_request(sys_prompt, prompt+file_contents)
                if response.status_code != 200:
                    return {"error (g)": response.text.encode().decode('unicode_escape')}

                project_name = file.split("/")[-1].split(".")[0]
                date = find_latest_modification(input_fp)

                final_response = finally_generation(
                    json.loads(response.text)["choices"][0]["message"]["content"],
                    date,
                    project_name
                )
                if final_response.status_code != 200:
                    return {"error (f)": final_response.text.encode().decode('unicode_escape')}

                with open(output_fp, 'a', encoding='utf-8') as md:
                    rt = json.loads(final_response.text)
                    md.write(rt["choices"][0]["message"]["content"])

    return {"path": output_fp}


def process(input_fp: str, output_fp: str) -> dict[str, str]:
    if input_fp.endswith(".zip"):
        return process_archive(input_fp, output_fp)
    else:
        return process_file(input_fp, output_fp)


def after_process(input_fp: str, output_fp: str, pr_name: str, date: str) -> None:
    with open(input_fp, 'r', encoding='utf-8') as file:
        lines = file.readlines()
        if lines and lines[0].lower().startswith("конечно"):
            lines.pop(0)

    lines = [line for line in lines if not line.lower().startswith("конечно")]

    project_name_pattern = re.compile(r'название_проекта|project_name|другое_название|ваше_название', re.IGNORECASE)

    lines = [project_name_pattern.sub(pr_name, line) for line in lines]

    with open(input_fp, 'w', encoding='utf-8') as file:
        file.writelines(lines)

    with open(output_fp, 'a', encoding='utf-8') as md:
        md.write(f"## {pr_name} ({date})\n")
        md.writelines(lines)

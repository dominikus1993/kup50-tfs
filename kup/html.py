
import os
from typing import Iterable


def write(html: Iterable[str], path: str, file_name: str):
    file_path = os.path.join(path, file_name)
    with open(file_path, 'w', encoding='utf-8') as file:
        for chunk in html:
            file.write(chunk)
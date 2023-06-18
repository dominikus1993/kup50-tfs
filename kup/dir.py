import os
from azure.devops.v7_1.git.models import GitRepository
import os
import zipfile
import shutil

PATH = "kup50"
def create_dir_if_not_exists(repo: GitRepository) -> str:
    repo_name: str = repo.name if repo.name is not None else ""
    path = os.path.join(PATH, repo_name)
    if not os.path.exists(path):
        os.makedirs(path)
    return path

def path_to_html_file_name(path: str) -> str:
    res = str.replace(path, '/', '_').replace('.', '_')
    return f'{res}.html'


def write_zip(path: str):
    if os.path.exists(path=path):
        os.remove(path=path)
    with zipfile.ZipFile(path, 'w') as zf:
        zf.write(os.path.basename(PATH))
        
def remove_old():
    shutil.rmtree(PATH)
import os
from azure.devops.v7_1.git.models import GitRepository
import os
import zipfile
import shutil

def create_dir_if_not_exists(repo: GitRepository, path: str) -> str:
    repo_name: str = repo.name if repo.name is not None else ""
    path = os.path.join(path, repo_name)
    if not os.path.exists(path):
        os.makedirs(path)
    return path

def path_to_html_file_name(path: str) -> str:
    res = str.replace(path, '/', '_').replace('.', '_')
    return f'{res}.html'


def write_zip(path: str, zip_path: str):
    if os.path.exists(path=path):
        os.remove(path=path)
    with zipfile.ZipFile(path, 'w') as zf:
        zf.write(os.path.basename(zip_path))
        
def remove_old(path: str):
    shutil.rmtree(path)
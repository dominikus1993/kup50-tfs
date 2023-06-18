import os
from azure.devops.v7_1.git.models import GitRepository

def create_dir_if_not_exists(repo: GitRepository) -> str:
    repo_name: str = repo.name if repo.name is not None else ""
    path = os.path.join('kup_data', repo_name)
    if not os.path.exists(path):
        os.makedirs(path)
    return path

def path_to_html_file_name(path: str) -> str:
    res = str.replace(path, '/', '_').replace('.', '_')
    return f'{res}.html'
from typing import Iterable
from azure.devops.v7_1.core.models import TeamProjectReference
from azure.devops.v7_1.git.git_client import GitClient
from azure.devops.v7_1.git.models import GitRepository, GitCommit, GitCommitChanges
from azure.devops.v7_1.git.models import GitQueryCommitsCriteria, GitObject
from azure.devops.exceptions import AzureDevOpsServiceError
from kup.stream import stream_to_unicode
from kup.dir import create_dir_if_not_exists,path_to_html_file_name
from kup.html import write
from difflib import HtmlDiff

def list_repositories(client: GitClient, projects: Iterable[TeamProjectReference]) -> Iterable[GitRepository]:
    for project in projects:
        repositories: list[GitRepository] | None  = client.get_repositories(project=project.id)
        if repositories is not None and len(repositories) > 0:
            for repo in repositories :
                if repo is not None:
                    yield repo
                    
def list_changes(client: GitClient, repo: GitRepository, author: str, from_date: str, to_date: str) -> Iterable[GitCommitChanges]:
    try:
        criteria = GitQueryCommitsCriteria(author=author,from_date=from_date, to_date=to_date)
        commits: list[GitCommit] = client.get_commits(repository_id=repo.id, search_criteria=criteria)
        for commmit in commits:
            changes: GitCommitChanges = client.get_changes(commit_id=commmit.commit_id, repository_id=repo.id)
            yield changes
    except AzureDevOpsServiceError:
      print(f'No access to repo: {repo.name}')
      
      
def process_and_write_changes(client: GitClient, repo: GitRepository, html_diff: HtmlDiff, tmp_path: str, changes: Iterable[GitCommitChanges]):
    for change in changes:
        path = create_dir_if_not_exists(repo, tmp_path)
        chans = change.changes
        if chans is not None:
            for chan in chans:
                if chan["item"]["gitObjectType"] == "blob":
                    if chan["changeType"] == "add":
                        file_path = chan["item"]["path"]
                        file = stream_to_unicode(client.get_blob_content(repo.id, chan['item']['objectId']))
                        if file is not None:
                            diff = html_diff.make_file([], file)
                            write(diff, path, path_to_html_file_name(file_path))
                    if chan["changeType"] == "edit":
                        file_path = chan["item"]["path"]
                        new = stream_to_unicode(client.get_blob_content(repo.id, chan['item']['objectId']))
                        old = stream_to_unicode(client.get_blob_content(repo.id, chan['item']['originalObjectId']))
                        if new is not None and old is not None:
                            res = html_diff.make_file(old, new, context=True)
                            write(res, path, path_to_html_file_name(file_path))
                    
                
                
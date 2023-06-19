from typing import Iterable, Sequence
from azure.devops.connection import Connection
from msrest.authentication import BasicAuthentication
from azure.devops.v7_1.git.models import GitQueryCommitsCriteria
from azure.devops.v7_1.core.core_client import CoreClient
from azure.devops.v7_1.git.git_client import GitClient
from azure.devops.v7_1.git.models import GitRepository
from azure.devops.v7_1.core.models import TeamProjectReference
import pprint
import os

def list_projects(client: CoreClient) -> Iterable[TeamProjectReference]:
    get_projects_response = client.get_projects()
    index = 0
    while get_projects_response is not None and len(get_projects_response) > 0:
        for project in get_projects_response :
            if project is not None:
                yield project
            index += 1
        get_projects_response = client.get_projects(continuation_token=index)
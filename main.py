from azure.devops.connection import Connection
from msrest.authentication import BasicAuthentication
from azure.devops.v7_1.git.models import GitQueryCommitsCriteria
from azure.devops.v7_1.core.core_client import CoreClient
from azure.devops.v7_1.core.models import TeamProjectReference
import pprint
import os
# Fill in with your personal access token and org URL
personal_access_token = os.environ["PAT_TOKEN"]
organization_url = os.environ["ORG"]

# Create a connection to the org
credentials = BasicAuthentication('', personal_access_token)

connection = Connection(base_url=organization_url, creds=credentials)
core_client: CoreClient = connection.clients.get_core_client()

# Get the first page of projects
get_projects_response: list[TeamProjectReference] | None = core_client.get_projects()

index = 0
while get_projects_response is not None and len(get_projects_response) > 0:
    for project in get_projects_response :
        pprint.pprint("[" + str(index) + "] " + project.name)
        index += 1
    get_projects_response = core_client.get_projects(continuation_token=index)
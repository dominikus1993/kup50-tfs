from difflib import HtmlDiff
from azure.devops.connection import Connection
from msrest.authentication import BasicAuthentication
from azure.devops.v7_1.git.models import GitQueryCommitsCriteria
from azure.devops.v7_1.core.core_client import CoreClient
from azure.devops.v7_1.git.git_client import GitClient
from azure.devops.v7_1.git.models import GitRepository
from azure.devops.v7_1.core.models import TeamProjectReference
from kup.dir import write_zip,remove_old
from kup.project import list_projects 
from kup.repo import list_repositories, list_changes, read_changes
from kup.date import get_first_day_of_month_when_none, get_last_day_of_month_when_none
import pprint
import os
# Fill in with your personal access token and org URL
personal_access_token = os.environ["PAT_TOKEN"]
organization_url = os.environ["ORG"]

# Create a connection to the org
credentials = BasicAuthentication('', personal_access_token)

connection = Connection(base_url=organization_url, creds=credentials)
core_client: CoreClient = connection.clients.get_core_client()
git_client: GitClient = connection.clients.get_git_client()
diff = HtmlDiff()
# Get the first page of projects
projects = list_projects(core_client)

repos = list_repositories(git_client, projects=projects)
for repo in repos:
    changes = list_changes(git_client, repo, "Dominik.Kotecki", from_date=get_first_day_of_month_when_none(None), to_date=get_last_day_of_month_when_none(None))
    chans = read_changes(git_client, repo, diff, changes)
    
write_zip("kup.zip")
remove_old()
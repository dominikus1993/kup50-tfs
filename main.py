from azure.devops.connection import Connection
from msrest.authentication import BasicAuthentication
from azure.devops.v7_1.git.models import GitQueryCommitsCriteria
from azure.devops.v7_1.core.core_client import CoreClient
# Fill in with your personal access token and org URL
personal_access_token = 'YOURPAT'
organization_url = 'https://dev.azure.com/YOURORG'

# Create a connection to the org
credentials = BasicAuthentication('', personal_access_token)

connection = Connection(base_url=organization_url, creds=credentials)
connection.get_client()
core_client: CoreClient = connection.clients.get_core_client()

# Get the first page of projects
get_projects_response = core_client.get_projects()
import pprint

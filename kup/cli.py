import subprocess
import click
from difflib import HtmlDiff
from azure.devops.connection import Connection
from msrest.authentication import BasicAuthentication
from azure.devops.v7_1.git.models import GitQueryCommitsCriteria
from azure.devops.v7_1.core.core_client import CoreClient
from azure.devops.v7_1.git.git_client import GitClient
from azure.devops.v7_1.git.models import GitRepository
from azure.devops.v7_1.core.models import TeamProjectReference
from kup.dir import write_zip,remove_old
from kup.name import random_kup_folder_name
from kup.project import list_projects 
from kup.repo import list_repositories, list_changes, process_and_write_changes
from kup.date import get_first_day_of_month_when_none, get_last_day_of_month_when_none
import pprint
import os
CONTEXT_SETTINGS = dict(help_option_names=['-h', '--help'])

@click.group(context_settings=CONTEXT_SETTINGS)
def cli():
    pass


@cli.command("diff")
@click.option('-p', '--pat', type=str, default=os.environ["PAT_TOKEN"], help='string')
@click.option('-o', '--org', type=str, default=os.environ["ORG"])
@click.option('-a', '--author', type=str, default="Dominik.Kotecki")
@click.option('-o', "--output", type=str, default="kup")
def diff2html(pat: str, org: str, author: str, output: str):
    click.echo("Diff creation start")
    click.echo("diff from log creation start")
    click.echo("pat: {}".format(pat))
    click.echo("org: {}".format(org))
    click.echo("author: {}".format(author))
    click.echo("output: {}".format(output))
    TMP = random_kup_folder_name()
    credentials = BasicAuthentication('', pat)

    connection = Connection(base_url=org, creds=credentials)
    core_client: CoreClient = connection.clients.get_core_client()
    git_client: GitClient = connection.clients.get_git_client()
    diff = HtmlDiff()
    # Get the first page of projects
    projects = list_projects(core_client)

    repos = list_repositories(git_client, projects=projects)
    for repo in repos:
        changes = list_changes(git_client, repo, author, from_date=get_first_day_of_month_when_none(None), to_date=get_last_day_of_month_when_none(None))
        process_and_write_changes(git_client, repo, diff, TMP, changes)

    write_zip(output, TMP)
    remove_old(TMP)
    click.echo("Diff creation end")
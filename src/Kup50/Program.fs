open HtmlDiff
open Microsoft.VisualStudio.Services.Common
open Microsoft.VisualStudio.Services.WebApi
open System
open FSharp.Control
open Microsoft.TeamFoundation.SourceControl.WebApi
open Kup50

let pat = System.Environment.GetEnvironmentVariable("PAT")
let uri = System.Environment.GetEnvironmentVariable("ORG_URL")
let project = System.Environment.GetEnvironmentVariable("PROJECT_NAME")

let creds  = VssBasicCredential("", pat)
let credentials = VssCredentials(creds)

// Connect to Azure DevOps Services
let connection = new VssConnection(Uri(uri), credentials)

let struct (firstDay, lastDay) = Date.getFirstAndLastMonthDay(DateTime.Today) |> Date.formatFirstAndLastMonthDay
// Get a GitHttpClient to talk to the Git endpoints
let gitClient = connection.GetClient<GitHttpClient>();
let res = Git.getRepoChanges(gitClient) (project) ("Dominik Kotecki") (firstDay) (lastDay) |> Git.writeChanges(gitClient) |> TaskSeq.toList


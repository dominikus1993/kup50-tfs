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
let connection = new VssConnection(Uri(uri), credentials);

// Get a GitHttpClient to talk to the Git endpoints
let gitClient = connection.GetClient<GitHttpClient>();

//let res = Git.getChanges(gitClient) (project) ("Dominik Kotecki") (DateTime(2023, 3, 1).ToString()) (DateTime(2023, 3, 31).ToString()) |> TaskSeq.toList
//
// for change in res do
//     printfn "Change by %A" change.RepoName
let result = (Date.getFirstAndLastMonthDay(DateTime.Today))
printfn "%A" result
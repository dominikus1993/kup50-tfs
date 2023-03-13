namespace Kup50

open System
open System.IO
open System.Threading
open HtmlDiff
open Microsoft.TeamFoundation.SourceControl.WebApi
open FSharp.Control

type TfsGitChange = { Changes: GitChange seq; RepoName: string; RepoId: Guid; Project: string }

module Git = 
    let private commitCriteria(author) (fromDate) (toDate) =
        let q = GitQueryCommitsCriteria()
        q.Author <- author
        q.FromDate <- fromDate |> Date.format
        q.ToDate <- toDate |> Date.format
        q

    let private getRepos (client: GitHttpClient) (orgName: string)=
        taskSeq{
            let! repos = client.GetRepositoriesAsync(orgName)
            for repo in repos do
                yield repo
        }

    let private getCommits (client: GitHttpClient) (repo: GitRepository) (query: GitQueryCommitsCriteria) =
        taskSeq{
            try
                let! commits = client.GetCommitsAsync(repo.Id, query)
                for commit in commits do
                    yield commit
            with
                 | _ -> printfn $"Error! in repo %A{repo.Name}";
        }
    
    let isValidChange (change: GitChange) =
        change.Item.GitObjectType = GitObjectType.Blob && (change.ChangeType = VersionControlChangeType.Add || change.ChangeType = VersionControlChangeType.Edit)
    
    let isValidCommit (commit: GitCommitRef) (fromDate) (toDate)  =
        commit.Author.Date >= fromDate && commit.Author.Date <= toDate
        
    let private getChanges (client: GitHttpClient) (repoId: Guid) (repoName: string) struct (fromDate, toDate) (commits: taskSeq<GitCommitRef>) =
        taskSeq {
            for commit in commits do
                printfn "Repo: %A ChangeDate:  %A" repoName commit.Author.Date
                let commitDate = commit.Author.Date |> DateOnly.FromDateTime
                if commitDate >= fromDate && commitDate <= toDate then
                    let! changes = client.GetChangesAsync(commit.CommitId, repoId)
                    for change in changes.Changes do
                        if change.Item.GitObjectType = GitObjectType.Blob && (change.ChangeType = VersionControlChangeType.Add || change.ChangeType = VersionControlChangeType.Edit) then
                            yield change
        }

    let getRepoChanges(client: GitHttpClient) (orgName) (author) (fromDate) (toDate) = 
        let queryCommit = commitCriteria author fromDate toDate
        taskSeq {
            for repo in getRepos(client)(orgName) do
                printfn "Repo: %A" repo.Name
                let! changes = getCommits(client) (repo) (queryCommit) |> getChanges (client) (repo.Id) (repo.Name) struct (fromDate, toDate) |> TaskSeq.toArrayAsync
                if changes.Length > 0 then
                    yield { Changes = changes; RepoName = repo.Name; RepoId = repo.Id; Project = orgName }
        }
   
    let private getBlob (client: GitHttpClient) (project: string) (repoId: Guid) (objectId: string) =
        task {
            return! client.GetBlobZipAsync(project=project, repositoryId=repoId, sha1 = objectId, download = Nullable(true), fileName = null, userState = null, cancellationToken = CancellationToken.None)
        }
    
    let writeChanges (client: GitHttpClient) (repoChanges: taskSeq<TfsGitChange>) =
        taskSeq {
            for repoChange in repoChanges do
                let dir = $"./kup/{repoChange.RepoName}"
                do Files.createDir(dir)
                for change in repoChange.Changes do
                    let operation = change.Item.Path.Replace("/", "_")
                    let file = Path.Join(dir, $"{repoChange.RepoName}_{operation}.html")
                    printfn "Change %A %A %A" change.Item.Path repoChange.RepoName change.ChangeType
                    match change.ChangeType with
                    | VersionControlChangeType.Add ->
                        use! blob = getBlob (client) (repoChange.Project) (repoChange.RepoId) (change.Item.ObjectId)
                        do! blob |> Stream.toString |> Files.writeString(file)
                    | VersionControlChangeType.Edit ->
                        use! oldBlob = getBlob (client) (repoChange.Project) (repoChange.RepoId) (change.Item.OriginalObjectId)
                        use! newBlob = getBlob (client) (repoChange.Project) (repoChange.RepoId) (change.Item.ObjectId)
                        let newHtml = newBlob |> Stream.toString
                        let oldHtml = oldBlob |> Stream.toString
                        let diff = Html.diff (newHtml) (oldHtml)
                        do! diff |> Files.writeString(file)
                    | _ ->
                        printfn "Nope"
                    
                    yield ()
        }
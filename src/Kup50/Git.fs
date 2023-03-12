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
        q.FromDate <- fromDate
        q.ToDate <- toDate
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
    
    let private getChanges (client: GitHttpClient) (repoId: Guid) (commits: taskSeq<GitCommitRef>) =
        taskSeq {
            for commit in commits do
                let! changes = client.GetChangesAsync(commit.CommitId, repoId)
                for change in changes.Changes do
                    if change.Item.GitObjectType = GitObjectType.Blob && (change.ChangeType = VersionControlChangeType.Add || change.ChangeType = VersionControlChangeType.Edit) then
                        yield change
        }

    let getRepoChanges(client: GitHttpClient) (orgName) (author) (fromDate) (toDate) = 
        let queryCommit = commitCriteria author fromDate toDate
        taskSeq {
            for repo in getRepos(client)(orgName) do
                let! changes = getCommits(client) (repo) (queryCommit) |> getChanges (client) (repo.Id) |> TaskSeq.toArrayAsync
                if changes.Length > 0 then
                    yield { Changes = changes; RepoName = repo.Name; RepoId = repo.Id; Project = orgName }
        }
   
    let getBlob (client: GitHttpClient) (project: string) (repoId: Guid) (objectId: string) =
        task {
            return! client.GetBlobContentAsync(project=project, repositoryId=repoId, sha1 = objectId, download = Nullable(true), fileName = null, userState = null, cancellationToken = CancellationToken.None)
        }
    
    let writeChanges (client: GitHttpClient) (repoChanges: taskSeq<TfsGitChange>) =
        taskSeq {
            for repoChange in repoChanges do
                printfn "GetRepoBlobs %A" repoChange.RepoName
                let dir = $"./kup/{repoChange.RepoName}"
                do Files.createDir(dir)
                for change in repoChange.Changes do
                    let operation = change.Item.Path.Replace("/", "_")
                    let file = Path.Join(dir, $"{repoChange.RepoName}_{operation}.html")
                    match change.ChangeType with
                    | VersionControlChangeType.Add ->
                        use! blob = getBlob (client) (repoChange.Project) (repoChange.RepoId) (change.Item.ObjectId)
                        do! Files.writeAll(file) (blob)
                    | VersionControlChangeType.Edit ->
                        use! oldBlob = getBlob (client) (repoChange.Project) (repoChange.RepoId) (change.Item.OriginalObjectId)
                        use! newBlob = getBlob (client) (repoChange.Project) (repoChange.RepoId) (change.Item.ObjectId)
                        let diffHelper = HtmlDiff(oldBlob |> Stream.toString, newBlob |> Stream.toString);
                        do! Files.writeString(file)(diffHelper.Build())
                    | _ ->
                        printfn "Nope"
                    
                    yield ()
        }
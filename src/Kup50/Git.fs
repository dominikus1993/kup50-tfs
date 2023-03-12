namespace Kup50

open System
open Microsoft.TeamFoundation.SourceControl.WebApi
open FSharp.Control

type TfsGitChange = { Change: GitChange; RepoName: string; RepoId: Guid }

module Git = 
    let commitCriteria(author) (fromDate) (toDate) =
        let q = GitQueryCommitsCriteria()
        q.Author <- author
        q.FromDate <- fromDate
        q.ToDate <- toDate
        q

    let getRepos (client: GitHttpClient) (orgName: string)=
        taskSeq{
            let! repos = client.GetRepositoriesAsync(orgName)
            for repo in repos do
                yield repo
        }

    let getCommits (client: GitHttpClient) (repo: GitRepository) (query: GitQueryCommitsCriteria) =
        taskSeq{
            try
                let! commits = client.GetCommitsAsync(repo.Id, query)
                for commit in commits do
                    yield commit
            with
                 | _ -> printfn $"Error! in repo %A{repo.Name}";
        }


    let getChanges(client: GitHttpClient) (orgName) (author) (fromDate) (toDate) = 
        let queryCommit = commitCriteria author fromDate toDate
        taskSeq {
            for repo in getRepos(client)(orgName) do
                for commit in getCommits(client) (repo) (queryCommit) do
                    let! changes = client.GetChangesAsync(commit.CommitId, repo.Id)
                    for change in changes.Changes do
                        if change.ChangeType = VersionControlChangeType.Add || change.ChangeType = VersionControlChangeType.Edit then
                            yield { Change = change; RepoName = repo.Name; RepoId = repo.Id }
        }

    // let writeChanges(changes: taskSeq<TfsGitChange>) =
    //     taskSeq {
    //         
    //     }
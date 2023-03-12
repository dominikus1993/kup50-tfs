#! /usr/bin/env node
import {getPersonalAccessTokenHandler, WebApi} from "azure-devops-node-api";
import date from "date-and-time";

function getFirstAndLastDayInMonth(now: Date) : { firstDay: string, lastDay: string } {
    const year = now.getFullYear()
    const month = now.getMonth()
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);

    return { firstDay: date.format(firstDay, "MM/DD/YYYY"), lastDay: date.format(lastDay, "MM/DD/YYYY"),}
}

export function connect(token: string, orgUrl: string) : WebApi {
    const authHandler = getPersonalAccessTokenHandler(token); 
    return new WebApi(orgUrl, authHandler);    
} 

export async function getChanges(webApi: WebApi, {project, author}: {project: string, author: string}) {
    const git = await webApi.getGitApi()
    const repositories = await git.getRepositories(project)
    const { firstDay, lastDay } = getFirstAndLastDayInMonth(new Date())
    console.log(firstDay, lastDay)
    for (const repository of repositories) {
        if (!repository.id) {
            continue
        }
        const commits = await git.getCommits(repository.id, { author: author, fromDate: firstDay, toDate: lastDay })

        for (const commit of commits) {
            const changes = await git.getChanges(commit.commitId, repository.id, project)
            if (!changes?.changes) {
                continue
            }
            for (const change of changes.changes) {
                console.log(`Commit: ${repository.name} ${commit.author?.email}, ${commit.author?.date}, ${change.changeType}`)
            }
            
        }
    }
}
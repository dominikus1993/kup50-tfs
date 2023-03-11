import {getPersonalAccessTokenHandler, WebApi} from "azure-devops";
import { getFirstAndLastDayInMonth } from "../datetime/index.ts";


export function connect(token: string, orgUrl: string) : WebApi {
    const authHandler = getPersonalAccessTokenHandler(token); 
    return new WebApi(orgUrl, authHandler);    
} 

export async function getChanges(webApi: WebApi, {project, author}: {project: string, author: string}) {
    const git = await webApi.getGitApi()
    const repositories = await git.getRepositories(project)
    const { firstDay, lastDay } = getFirstAndLastDayInMonth(new Date())
    for (const repository of repositories) {
        if (!repository.id) {
            continue
        }
        const commits = await git.getCommits(repository.id, { author: author, fromDate: firstDay, toDate: lastDay })

        for (const commit of commits) {
            console.log(`Commit: ${repository.name} ${commit.author?.email}, ${commit.author?.date}`)
        }
    }
}
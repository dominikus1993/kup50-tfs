open Argu
open Argu.ArguAttributes
open HtmlDiff
open Microsoft.VisualStudio.Services.Common
open Microsoft.VisualStudio.Services.WebApi
open System
open FSharp.Control
open Microsoft.TeamFoundation.SourceControl.WebApi
open Kup50


type CmdArgs =
    | [<AltCommandLine("-g")>] Generate of proj: string * pat: string * orgUrl: string * author: string
with
    interface IArgParserTemplate with
        member this.Usage =
            match this with
            | Generate _ -> "Generate a kup diff"

let errorHandler = ProcessExiter(colorizer = function ErrorCode.HelpText -> None | _ -> Some ConsoleColor.Red)
let parser = ArgumentParser.Create<CmdArgs>(programName = "kup", errorHandler = errorHandler)

printfn "%A" (Environment.GetCommandLineArgs())
match parser.ParseCommandLine(Environment.GetCommandLineArgs()) with
| p when p.Contains(Generate) ->
    let (project, pat, orgUrl, author) = p.GetResult(Generate)
    let creds  = VssBasicCredential("", pat)
    let credentials = VssCredentials(creds)
    let connection = new VssConnection(Uri(orgUrl), credentials)
    let struct (firstDay, lastDay) = Date.getFirstAndLastMonthDay(DateTime.Today) |> Date.formatFirstAndLastMonthDay
    // Get a GitHttpClient to talk to the Git endpoints
    let gitClient = connection.GetClient<GitHttpClient>();
    let res = Git.getRepoChanges(gitClient) (project) author (firstDay) (lastDay) |> Git.writeChanges(gitClient) |> TaskSeq.toList
    Files.writeDirToZip "kup" "kup.zip"
    Files.removeDir "kup"
| _ ->
    raise (InvalidOperationException("Arguments not specified"))
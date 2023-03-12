namespace Kup50

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


module Program =
    [<EntryPoint>]
    let main argv =
        let errorHandler = ProcessExiter(colorizer = function ErrorCode.HelpText -> None | _ -> Some ConsoleColor.Red)
        let parser = ArgumentParser.Create<CmdArgs>(programName = "kup", errorHandler = errorHandler)
        
        printfn "%A" (argv)
        let parseResult = parser.ParseCommandLine(argv)
        printfn "Got parse results %A" <| parseResult.GetAllResults()
        let (project, pat, orgUrl, author) = parseResult.GetResult(Generate)
        
        printfn "Start processing kup"
        let creds  = VssBasicCredential("", pat)
        let credentials = VssCredentials(creds)
        let connection = new VssConnection(Uri(orgUrl), credentials)
        let struct (firstDay, lastDay) = Date.getFirstAndLastMonthDay(DateTime.Today) |> Date.formatFirstAndLastMonthDay
        // Get a GitHttpClient to talk to the Git endpoints
        let gitClient = connection.GetClient<GitHttpClient>();
        do Git.getRepoChanges(gitClient) (project) author (firstDay) (lastDay) |> Git.writeChanges(gitClient) |> TaskSeq.toList |> ignore
        Files.removeDir "kup"
        0

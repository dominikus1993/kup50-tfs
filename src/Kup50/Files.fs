namespace Kup50

open System.IO

module Files =
    let createDir (dirName: string) =
        Directory.CreateDirectory(dirName)

    let writeAll(fileName: string)(stream: Stream) =
        task {
            use fileStream = new FileStream(fileName, FileMode.Create ||| FileMode.Append)
            stream.Seek(0, SeekOrigin.Begin) |> ignore;
            do! stream.CopyToAsync(fileStream)
        }

    let writeString(fileName: string)(stream: string) =
        File.WriteAllTextAsync(fileName, stream)
 
 
module Stream =
    let toString(stream: Stream) =
        let reader = new StreamReader(stream);
        reader.ReadToEnd();

        
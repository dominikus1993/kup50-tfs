namespace Kup50

open System.IO
open System.IO.Compression

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
 
    let writeDirToZip(dir: string) (zipFile: string) =
        ZipFile.CreateFromDirectory(dir, zipFile);

    let removeDir(dir: string)=
        Directory.Delete(dir)
        
module Stream =
    let toString(stream: Stream) =
        let reader = new StreamReader(stream);
        reader.ReadToEnd();

        
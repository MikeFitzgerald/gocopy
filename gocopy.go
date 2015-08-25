package main

import (
    "fmt"
    "flag"
    "io"
    "log"
    "os"
)

func CopyFile(source string, dest string) (err error) {
    sourcefile, err := os.Open(source)
    if err != nil {
        return err
    }

    defer sourcefile.Close()

    destfile, err := os.Create(dest)
    if err != nil {
        return err
    }

    defer destfile.Close()
    _, err = io.Copy(destfile, sourcefile)
    if err == nil {
        sourceinfo, err := os.Stat(source)
        if err != nil {
            err = os.Chmod(dest, sourceinfo.Mode())
        }
    }

    return
}

func CopyDir(source string, dest string) (err error) {

    // Get properties of source directory
    sourceinfo, err := os.Stat(source)
    if err != nil {
        return err
    }

    // Create destination directory
    err = os.MkdirAll(dest, sourceinfo.Mode())
    if err != nil {
        return err
    }

    directory, _ := os.Open(source)

    objects, err := directory.Readdir(-1)

    for _, obj := range objects {

        sourcefilepointer := source + "/" + obj.Name()

        destfilepointer := dest + "/" + obj.Name()

        if obj.IsDir() {
            
            // Create sub-directories - recursively
            err = CopyDir(sourcefilepointer, destfilepointer)
            if err != nil {
                fmt.Println(err)
            }
        } else {
            err = CopyFile(sourcefilepointer, destfilepointer)
            if err != nil {
                fmt.Println(err)
            }
        }
    }
    return
}

func main() {
    
    // Get the source and destination of the directory
    flag.Parse()

    source_dir := flag.Arg(0)

    dest_dir := flag.Arg(1)

    fmt.Println("Source: " + source_dir)

    // Check if the source directory exists
    src, err := os.Stat(source_dir)
    if err != nil {
        log.Fatal(err)
    }

    if !src.IsDir() {
        fmt.Println("Source is not a directory")
        os.Exit(1)
    }

    // Create the destination directory
    fmt.Println("Destination:" + dest_dir)

    _, err = os.Open(dest_dir)
    if !os.IsNotExist(err) {
        fmt.Println("Desintation directory already exists. Abort!")
        os.Exit(1)
    }

    err = CopyDir(source_dir, dest_dir)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("Directory copied")
    }

}
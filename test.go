package main

import (
"log"
"encoding/json"
"github.com/fsnotify/fsnotify"
"io/ioutil"
"fmt"
"os"
"strings"
"os/exec"
"bufio"
)

var (
    explode []int
    dirs []string
    source string
    dest string
)

type Config struct {
    Source []string
    Dest []string
    Type string
}

func Watcher() {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    done := make(chan bool, 1)
    go func() {
        for {
            select {
            case event := <-watcher.Events:
                // Reset destination and destination directory
                destDir := ""
                dest = getDest()

                // Get file name
                explode := strings.Split(event.Name, "/")
                fileName := explode[len(explode) - 1]

                // Get final destination path and destination directory
                dest = strings.Replace(event.Name, source, dest, -1)
                destDir = strings.Replace(dest, fileName, "", -1)

                if event.Op&fsnotify.Remove == fsnotify.Remove {
                err := os.Remove(dest)
                check(err)

                } else {

                log.Println("event:", event)

                log.Println("modified file:", event.Name)

                
                file, err := os.Stat(event.Name)
                checkFatal(err)

                if file.Mode().IsDir() {
                    err = os.Mkdir(dest, 0755)
                    check(err)
                    err = watcher.Add(event.Name)
                    check(err)
                } else {

                    fmt.Println("Source: ", source, "Dest :", dest, "Dest Dir: ", destDir)

                    err = os.MkdirAll(destDir, 0755)
                    check(err)
                    out, err := exec.Command("/bin/cp", "-p", event.Name, dest).Output()
                    if(err != nil) {
                        fmt.Println("err: ", err)
                    } else {
                        fmt.Println("good: ", out)
                    }

                    if(err != nil) {
                        fmt.Println("N-a mers sa copiez: ", err)
                    }
                }
            }
        case err := <-watcher.Errors:
            log.Println("error:", err)
        }
    }
    }()

    dirs, err := readLines("dirs")
    for _, el := range dirs {
        err = watcher.Add(el)
    }
    check(err)
    <-done
}

func findDirs() []string {
    var c Config
    config, err := ioutil.ReadFile("/etc/gosync/config.json")

    if(err != nil) {
        fmt.Println("Couldn't open config file or the file is missing.");
    }

    if err := json.Unmarshal(config, &c); err != nil {
        panic("Config JSON might be corrupt.")
    }

    dir := c.Source[0]
    source = c.Source[0]

    cmd, err := exec.Command("find", dir, "-type", "d").Output()

    if(err != nil) {
        log.Fatal(err)
    }

    f, err := os.Create("dirs")
    checkFatal(err)

    _, err = f.Write(cmd)
    checkFatal(err)


    return nil

}

func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func writeLines(lines []string, path string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    w := bufio.NewWriter(file)
    for _, line := range lines {
        fmt.Fprintln(w, line)
    }
    return w.Flush()
}


func getDest() string {
    var c Config
    config, err := ioutil.ReadFile("/etc/gosync/config.json")
    if(err != nil) {
        fmt.Println("Couldn't open config file or the file is missing.");
    }

    if err := json.Unmarshal(config, &c); err != nil {
        panic("Config JSON might be corrupt.")
    }

    dest = c.Dest[0]

    return string(c.Dest[0])
}

func check(e error) {
    if e != nil {
        fmt.Println(e)
    }
}

func checkFatal(e error) {
    if e != nil {
        log.Fatal(e)
    }
}

func main() {
    findDirs()
    Watcher()
}

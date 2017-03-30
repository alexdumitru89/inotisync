package main
import (
"log"
//"encoding/json"
"github.com/fsnotify/fsnotify"
"io/ioutil"
"fmt"
"os"
"strings"
"os/exec"
"bufio"
"github.com/BurntSushi/toml"
"bytes"
)

var (
    explode []int
    dirs []string
    source string
    mode string
    remote bool
    dest string
    configFile string = "/home/go/inotisync/inotisync.conf"
)

type Config struct {
    Source []string
    Dest []string
    Remote bool
    Mode string
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
                // Reset destination directory
                destDir := ""

                // Get file name
                explode := strings.Split(event.Name, "/")
                fileName := explode[len(explode) - 1]

                // Get final destination path and destination directory
                destPath := strings.Replace(event.Name, source, dest, -1)
                destDir = strings.Replace(destPath, fileName, "", -1)

                // If a file/folder has been deleted
                if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
                err := os.Remove(destPath)
                check("Can't delete " + destPath + "\n", err)

                // Any other case (create, write, chmod)
                } else {

                log.Println("event:", event)

                log.Println("modified file:", event.Name)

                file, err := os.Stat(event.Name)
                checkFatal(err)

                // If modified file is a directory
                if file.Mode().IsDir() {
                    err = os.Mkdir(destPath, 0755)
                    check("Can't create directory\n", err)
                    err = watcher.Add(event.Name)
                    check("Can't add to watcher'n",err)
                } else {

                    fmt.Println("Source: ", source, "Dest :", destPath, "Dest Dir: ", destDir)

                    err = os.MkdirAll(destDir, 0755)
                    check("Can't create directories'n", err)
                    out, err := exec.Command("/bin/cp", "-p", event.Name, destPath).Output()
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
    check("Can't add to watcher\n", err)
    <-done
}

func findSourceDirs() []string {

	// Check if source path exists and is a directory
	file, err := os.Stat(source)
	if(err != nil) {
		fmt.Println("Directory doesn't exist. Attempting to create now")
		err = os.Mkdir(source, 0755)
		if(err != nil) {
			panic(err)
		}
	} else if(file.Mode().IsRegular()) {
		panic("Source path is not a directory. It's a regular file.")
	}


    cmd := exec.Command("find", source, "-type", "d")
    var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
    err = cmd.Run()

    if(err != nil) {
        log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
    }

    f, err := os.Create("dirs")
    checkFatal(err)

    _, err = f.Write(out.Bytes())
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


func readConf() {
	var c Config

	config, err := ioutil.ReadFile(configFile)
	conf := string(config)

	if(err != nil) {
        panic("Couldn't open config file or the file is missing.");
    }

	if _, err := toml.Decode(conf, &c); err != nil {
        panic("Config file might be corrupt.")
    }



	source = c.Source[0]
	dest = c.Dest[0]
	mode = c.Mode
	remote = c.Remote
}

 


func main() {
	readConf()
    findSourceDirs()
    Watcher()
}

package main
import (
"log"
//"encoding/json"
"github.com/fsnotify/fsnotify"

"fmt"
"os"
"strings"
"os/exec"
"bufio"
"bytes"
) 

var (
    explode []int
    dirs []string
    source string
    mode string
    remote bool
    dest string
    sources Sources
    configFile string = "/home/go/inotisync/src/inotisync.conf"
    ctr int
    destDir string
)

type Sources struct {
	Sync []Sync
}

type Sync struct { 
	Source string
	Destinations []string
}

func (s Sources) syncInit() {
	for i := 0; i < len(s.Sync); i++ {
		for j := 0; j < len(s.Sync[i].Destinations); j++ {	
			// Check if destination is remote or not
			// Check for the existence of the ":" character
			if strings.ContainsAny(s.Sync[i].Destinations[j], ":") {
				_, err := exec.Command("rsync", "-avzS", s.Sync[i].Source, s.Sync[i].Destinations[j]).Output()
				if(err != nil) {
	                fmt.Println("Can't copy")
	                fmt.Println(err)
	            }
	        } else {
	        	cmd := exec.Command("/bin/sh", "-c", "cp -Rp " + s.Sync[i].Source + "*" + " " + s.Sync[i].Destinations[j])
	        	stdoutErr, err := cmd.CombinedOutput()
	        	if(err != nil) {
	                fmt.Println("Can't copy")
	                fmt.Println(err)
	                fmt.Printf("%s\n", stdoutErr)
	            }
	        }
		}
	}
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

            	// Find which [sync] the file belongs to
            	for i:= 0; i < len(sources.Sync); i++ {
					if strings.HasPrefix(event.Name, sources.Sync[i].Source) {
						ctr = i
					}
				}

                // Get file name
                explode := strings.Split(event.Name, "/")
                fileName := explode[len(explode) - 1]
                // Loop through destinations
                for _, dest := range sources.Sync[ctr].Destinations {
                	destPath := strings.Replace(event.Name, sources.Sync[ctr].Source, dest, -1)
                	destDir = strings.Replace(destPath, fileName, "", -1)

                	// File or folder has been deleted
                	if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
                		if(!isRemoteDestination(destPath)) {
	                		err := os.Remove(destPath)
	                		check("Can't delete " + destPath + "\n", err)
	                	}
                	// File or folder has been created, modified
                	} else {
                		// Local destination. Use CP
                		if(!isRemoteDestination(destPath)) {
                			destPath = destPath + "/"
                			fmt.Println("Event: ", event.Name)
                			fmt.Println("Dest: ", destPath)

                			file, err := os.Stat(event.Name)
                			checkFatal(err)

                			// Path is a directory
	                		if file.Mode().IsDir() {
			                    err = os.Mkdir(destPath, 0755)
			                    check("Can't create directory\n", err)
			                    err = watcher.Add(event.Name)
			                    check("Can't add to watcher'n",err)
			                // Path is a file
			                } else {

			                    fmt.Println("Source: ", source, "Dest :", destPath, "Dest Dir: ", destDir)

			                    err = os.MkdirAll(destDir, 0755)
			                    check("Can't create directories'n", err)
			                    _, err := exec.Command("/bin/cp", "-p", event.Name, destPath).Output()
			                    if(err != nil) {
			                        fmt.Println("Can't copy")
			                        fmt.Println(err)
			                    }

			                    if(err != nil) {
			                        fmt.Println("N-a mers sa copiez: ", err)
			                    }
			                }
			            // Remote destination. Use rsync
                		} else {
                			_, err := exec.Command("rsync", "-avzS", event.Name, destPath).Output()
                			fmt.Println("rsync!")
                			if(err != nil) {
		                        fmt.Println("Rsync error")
		                        fmt.Println(err)
		                    }
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

func isRemoteDestination(destPath string) bool {
	if strings.ContainsAny(destPath, ":") {
		return true
	}
	return false
}

// Obsolete
// func findSourceDirs() []string {

// 	// Check if source path exists and is a directory
// 	file, err := os.Stat(source)
// 	if(err != nil) {
// 		fmt.Println("Directory doesn't exist. Attempting to create now")
// 		err = os.Mkdir(source, 0755)
// 		if(err != nil) {
// 			panic(err)
// 		}
// 	} else if(file.Mode().IsRegular()) {
// 		panic("Source path is not a directory. It's a regular file.")
// 	}


//     cmd := exec.Command("find", source, "-type", "d")
//     var out bytes.Buffer
// 	var stderr bytes.Buffer
// 	cmd.Stdout = &out
// 	cmd.Stderr = &stderr
//     err = cmd.Run()

//     if(err != nil) {
//         log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
//     }

//     f, err := os.Create("dirs")
//     checkFatal(err)

//     _, err = f.Write(out.Bytes())
//     checkFatal(err)


//     return nil

// }

func (s Sources) findSources() {

	f, err := os.Create("dirs")
    checkFatal(err)
    err = f.Truncate(0)

	for _, v := range s.Sync {

		file, err := os.Stat(v.Source)

		if(err != nil) {
			fmt.Println("Directory doesn't exist. Attempting to create now")
			err = os.Mkdir(source, 0755)
		}

		if(err != nil) {
			panic(err)
		} else if(file.Mode().IsRegular()) {
			panic("Source path is not a directory. It's a regular file.")
		}

		fmt.Println("Looking for directories in ", v.Source)
		cmd := exec.Command("find", v.Source, "-type", "d")
	    var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
	    err = cmd.Run()

	    if(err != nil) {
	        log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
	    }

	    _, err = f.Write(out.Bytes())
    	checkFatal(err)
	} 
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


 


func main() {
	sources.readConf()
    sources.findSources()
    Watcher()
}

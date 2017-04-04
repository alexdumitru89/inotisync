package main
import (
	"fmt"
	"log"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"strings"
)

func check(m string, e error) {
    if e != nil {
        fmt.Println(m, e)
    }
}
 
func checkFatal(e error) {
    if e != nil {
        log.Fatal("Fatal error:", e)
    }
}

func (s *Sources) readConf() *Sources {
	//var sources Sources

	config, err := ioutil.ReadFile(configFile)
	conf := string(config)

	if(err != nil) {
        panic("Couldn't open config file or the file is missing.");
    }

	if _, err := toml.Decode(conf, &s); err != nil {
		log.Fatal(err)
        panic("Config file might be corrupt.")
    }

    for i := 0; i < len(s.Sync); i++ {
    	// Add trailing slash to sources
    	if !strings.HasSuffix(s.Sync[i].Source, "/") {
    		s.Sync[i].Source = s.Sync[i].Source + "/"
    	}

    	// Add starting slash to sources
    	if !strings.HasPrefix(s.Sync[i].Source, "/") {
    		s.Sync[i].Source = "/" + s.Sync[i].Source
    	}

    	for j := 0; j < len(s.Sync[i].Destinations); j++ {
    		// Add trailing slash to destinations
    		if !strings.HasSuffix(s.Sync[i].Destinations[j], "/") {
    			s.Sync[i].Destinations[j] = s.Sync[i].Destinations[j] + "/"
    		}

    		// Add beginning slash to destinations
    		// Check if they are remote
    		if strings.ContainsAny(s.Sync[i].Destinations[j], ":") && !strings.Contains(s.Sync[i].Destinations[j], ":/") {
    			fmt.Println("A GASIT")
    			s.Sync[i].Destinations[j] = strings.Replace(s.Sync[i].Destinations[j], ":", ":/", -1)
    		} else if !strings.ContainsAny(s.Sync[i].Destinations[j], ":") && !strings.HasPrefix(s.Sync[i].Destinations[j], "/") {
    			s.Sync[i].Destinations[j] = "/" + s.Sync[i].Destinations[j]
    		}
    	}

    }

    fmt.Printf("Sources: %v\n", s.Sync)

	sources.syncInit()

	return s
}
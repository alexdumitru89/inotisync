package main
import (
	"fmt"
	"log"
	"github.com/BurntSushi/toml"
	"io/ioutil"
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

    fmt.Printf("tot: %v atat\n", s)
    fmt.Println("ln", s.Sync[0].Source)
 
	for _, v := range s.Sync {
		fmt.Printf("Sources: %v %v atat\n", v.Source, v.Destinations)
	}

	return s
}
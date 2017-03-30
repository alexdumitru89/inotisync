package main
import (
	"fmt"
	"log"
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

package main

import (
    "log"
    "os"
	"fmt"
    "path/filepath"
    "strings"
	"flag"

    "github.com/peterh/liner"
	"github.com/xoxzo/myssh"
)

var (
    history_fn = filepath.Join(os.TempDir(), ".liner_example_history")
    names      = []string{"john", "james", "mary", "nancy"}
)

type Hosts []string
func (i *Hosts) String() string {
	return fmt.Sprintf("%s", *i)
}
func (i *Hosts) Set(value string) error {
	for _, v := range strings.Split(value, ",") {
		*i = append(*i, v)
	}
	return nil
}

var hosts Hosts

func main() {
	flag.Var(&hosts, "h", "List of hosts")
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		return
	}
    line := liner.NewLiner()
    defer line.Close()

    line.SetCtrlCAborts(true)

    line.SetCompleter(func(line string) (c []string) {
        for _, n := range names {
            if strings.HasPrefix(n, strings.ToLower(line)) {
                c = append(c, n)
            }
        }
        return
    })

    if f, err := os.Open(history_fn); err == nil {
        line.ReadHistory(f)
        f.Close()
    }

	for {
		if command, err := line.Prompt("$ "); err == nil {
			//log.Print("Got: ", command)
			line.AppendHistory(command)
			if command == "exit" {
				return
			}

			myssh.Run(command, hosts)
		} else if err == liner.ErrPromptAborted {
			log.Print("Aborted")
		} else {
			log.Print("Error reading line: ", err)
		}

		if f, err := os.Create(history_fn); err != nil {
			log.Print("Error writing history file: ", err)
		} else {
			line.WriteHistory(f)
			f.Close()
		}
	}
}

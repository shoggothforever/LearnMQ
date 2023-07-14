package mes

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type Request struct {
	ID          int64  `json:"ID,omitempty"`
	Cost        int64  `json:"cost,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

func BodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "just say some thing !"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
func RouterKey(args []string) string {
	if len(args) != 3 {
		fmt.Println("please type as follow form :go run xxx.go <router_key>")
		return "anonymous.infos"
	} else {
		fmt.Println("if success, the exchanger will send message to queues whose binding_key are matched with router_key ", args[2])
		return args[2]
	}
}

package main

import (
	"fmt"
	"regexp"
)

func main() {
	comment := "  sddafdasfa   @@protogen:pkt_server_registration@@  // asfasdfasd"
	re := regexp.MustCompile("@@protogen:(.*?)@@")
	matches := re.FindStringSubmatch(comment)

	for i, match := range matches {
		fmt.Printf("[%v] %v\n", i, match)
	}

}

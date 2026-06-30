package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/egoroof/wow-server-ping/pkg/wow"
	"golang.org/x/term"
)

var PORT = flag.Int("port", 3724, "realmlist server port")
var TIMEOUT = flag.Duration("timeout", time.Second*10, "timeout for network operations")

func main() {
	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Println("Usage: realmlist user@host")
		os.Exit(1)
	}

	userHost := strings.Split(os.Args[1], "@")
	if len(userHost) < 2 {
		fmt.Println("Usage: realmlist user@host")
		os.Exit(1)
	}
	user := userHost[0]
	host := userHost[1]

	fmt.Print("Enter password: ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("")

	address := fmt.Sprintf("%v:%v", host, *PORT)
	client := wow.NewWowClient(address, user, string(password), *TIMEOUT)

	err = client.Login()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	realms := client.GetRealmList()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, realm := range realms {
		fmt.Fprintf(w, "%v\t%v\n", realm.Name, realm.Address)
	}
	w.Flush()

	filename := fmt.Sprintf("./servers/%v.json", host)
	json, err := json.MarshalIndent(realms, "", "	")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = os.WriteFile(filename, json, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Saved to %v\n", filename)
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/egoroof/wow-server-ping/pkg/wow"
)

var HOST = flag.String("host", "logon.wowcircle.me", "realmlist server host")
var PORT = flag.Int("port", 3724, "realmlist server port")
var USERNAME = flag.String("username", "", "username")
var PASSWORD = flag.String("password", "", "password")
var TIMEOUT = flag.Duration("timeout", time.Second*10, "timeout for network operations")

func main() {
	flag.Parse()

	if *HOST == "" {
		fmt.Println("Host cannot be empty")
		os.Exit(1)
	}
	if *USERNAME == "" {
		fmt.Println("Username cannot be empty")
		os.Exit(1)
	}
	if *PASSWORD == "" {
		fmt.Println("Password cannot be empty")
		os.Exit(1)
	}

	address := fmt.Sprintf("%v:%v", *HOST, *PORT)
	fmt.Printf("-> %v@%v\n", address, *USERNAME)

	client := wow.NewWowClient(address, *USERNAME, *PASSWORD, *TIMEOUT)

	err := client.Login()
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

	filename := fmt.Sprintf("./servers/%v.json", *HOST)
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

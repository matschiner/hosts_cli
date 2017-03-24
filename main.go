package main

import (
	"fmt"
	"log"
	"os"
)

const filename = "/etc/hosts"

func main() {
	log.SetFlags(0)

	if !amIRoot() {
		log.Fatal("Please run this program as Root!")
	}

	if len(os.Args) <= 1 {
		log.Fatal("Nothing to do, please specify command")
	}

	hosts := hostlist{}
	hosts.Read(filename)

	command := string(os.Args[1])

	// fmt.Println(command)
	switch command {
	default:
		log.Fatalf("Unknown command: %s", command)

	case "list", "ls":
		fmt.Println(string(hosts.Bytes()))

	case "del", "rm", "-":
		if len(os.Args) != 3 {
			log.Fatal("Please give an IP or hostname to delete")
		}
		hosts.Remove(os.Args[2])

	case "ucom":
		if len(os.Args) != 3 {
			log.Fatal("Please give an IP or hostname to uncomment")
		}
		hosts.Uncomment(os.Args[2])

	case "com":
		if len(os.Args) != 3 {
			log.Fatal("Please give an IP or hostname to comment out")
		}
		hosts.Comment(os.Args[2])

	case "add", "+":
		if len(os.Args) != 4 {
			log.Fatal("Please give arguments in the form ip, hostname")
		}
		err := hosts.Add(os.Args[2], os.Args[3])
		if err != nil {
			log.Fatal(err)
		}

	case "has", "?", "contains":
		if len(os.Args) != 4 {
			log.Fatal("Please give arguments in the form ip, hostname")
		}

		yes, err := hosts.Contains(os.Args[2], os.Args[3])
		if err != nil {
			log.Fatal(err)
		}

		if yes {
			os.Exit(0) // exit code 0 means it was contained within
		}

		os.Exit(1) // exit code 1 means not contained within

	}

	if hosts.changed {
		log.Printf("writing changes to %s", filename)
		hosts.Write(filename)
	}
}

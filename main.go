package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type line struct {
	ip, hostname string
}

const filename = "/etc/hosts"

func init() {
	if len(os.Getenv("SUDO_USER")) == 0 {
		fmt.Println("Please run this program as Root!")
		os.Exit(0)
	}
}
func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Please input command!")
		os.Exit(0)
	}
	command := os.Args[1]

	// fmt.Println(command)
	switch {
	default:
		fmt.Println("Your command is undefined!")
		os.Exit(0)
	case string(command) == "list":
		entries, comments := read()
		list(entries, comments)

	case string(command) == "rm":
		if len(os.Args) != 3 {
			fmt.Println("Please input a domain or ip address to delete!")
			os.Exit(0)
		}
		modify(string(os.Args[2]), "del")

	case string(command) == "ucom":
		if len(os.Args) != 3 {
			fmt.Println("Please input a domain or ip address to delete!")
			os.Exit(0)
		}
		modify(string(os.Args[2]), "uncomment")

	case string(command) == "com":
		if len(os.Args) != 3 {
			fmt.Println("Please input a domain or ip address to delete!")
			os.Exit(0)
		}
		modify(string(os.Args[2]), "comment")

	case string(command) == "add":
		if len(os.Args) != 4 {
			fmt.Println("Please input more arguments!")
			os.Exit(0)
		}
		add(string(os.Args[2]), string(os.Args[3]))

	}

}
func modify(modifyString string, action string) {
	inb, _ := ioutil.ReadFile(filename)
	in := string(inb)
	ins := strings.Split(in, "\n")
	var startcalled bool
	var out string

	for i := range ins {
		if len(ins[i]) != 0 {
			if ins[i] == "#hostedit" {
				startcalled = true
				out += ins[i] + "\n"

			} else if startcalled == true {
				inss := strings.Split(string(ins[i]), " ")
				if len(inss) > 1 {
					l := line{inss[0], inss[1]}
					l.ip = strings.Replace(l.ip, "#", "", -1)
					switch action {
					case "del":
						if l.ip != modifyString && l.hostname != modifyString {
							out += ins[i] + "\n"

						} else {
							fmt.Printf("Deleted %s -- %s!\n", l.ip, l.hostname)
						}
					case "comment":
						if l.ip == modifyString || l.hostname == modifyString {
							out += "#" + ins[i] + "\n"
							fmt.Printf("Commented %s -- %s!\n", l.ip, l.hostname)
						} else {
							out += ins[i] + "\n"
						}
					case "uncomment":
						if l.ip == modifyString || l.hostname == modifyString {
							ins[i] = strings.Replace(ins[i], "#", "", -1)
							out += ins[i] + "\n"
							fmt.Printf("Uncommented %s -- %s!\n", l.ip, l.hostname)
						} else {
							out += ins[i] + "\n"
						}
					}

					//fmt.Println("added -- "+ins[i])
				}

			} else {
				out += ins[i] + "\n"
			}

		}

	}
	if in != out {
		ioutil.WriteFile(filename, []byte(out), 0644)
	}
}
func add(ip string, hostname string) {
	inb, _ := ioutil.ReadFile(filename)
	in := string(inb)
	var out string
	if in[len(in)-1:len(in)] == "\n" {
		out = fmt.Sprintf("%s%s %s", in, ip, hostname)
	} else {
		out = fmt.Sprintf("%s\n%s %s", in, ip, hostname)
	}
	err := ioutil.WriteFile(filename, []byte(out), 0644)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Added: IP: %s \t HOSTNAME: %s\n", ip, hostname)
	}
}
func read() ([]line, []line) {
	var entries, comments []line
	inb, _ := ioutil.ReadFile(filename)
	in := string(inb)
	ins := strings.Split(in, "\n")
	var startcalled bool
	for i := range ins {
		if len(ins[i]) != 0 {

			if string(ins[i][0]) == "#" && ins[i] == "#hostedit" {
				startcalled = true
			} else if startcalled == true {
				inss := strings.Split(string(ins[i]), " ")
				if len(inss) > 1 {
					l := line{inss[0], inss[1]}
					if inss[0][0:1] == "#" {
						l.ip = strings.Replace(l.ip, "#", "", -1)
						comments = append(comments, l)
					} else {
						entries = append(entries, l)
					}
					//fmt.Println("added -- "+ins[i])
				}

			}

		}

	}

	return entries, comments
}
func list(entries []line, comments []line) {
	if len(entries)+len(comments) == 0 {
		fmt.Println("There was nothing added from hostsedit.")
	} else {
		if len(entries) > 0 {
			fmt.Println("Hostnames:")
			for _, h := range entries {
				fmt.Printf("IP: %s \t", h.ip)
				if len(h.ip) < 10 {
					fmt.Printf("\t")
				}
				fmt.Printf("HOSTNAME: %s \n", h.hostname)
			}
		}
		if len(comments) > 0 {
			fmt.Println("\nComments:")
			for _, h := range comments {
				fmt.Printf("IP: %s \t", h.ip)
				if len(h.ip) < 10 {
					fmt.Printf("\t")
				}
				fmt.Printf("HOSTNAME: %s \n", h.hostname)
			}
		}
	}
}

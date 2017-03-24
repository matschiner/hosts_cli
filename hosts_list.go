package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"
)

type hostlist struct {
	lines   []string
	changed bool
}

func (hl *hostlist) Read(fn string) error {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}

	hl.Parse(b)
	return nil
}

func (hl *hostlist) Parse(b []byte) {
	s := condenseNewlines(string(b))
	hl.lines = strings.Split(s, "\n")
}

func (hl *hostlist) Write(fn string) error {
	err := ioutil.WriteFile(fn, hl.Bytes(), 0644)

	if err == nil {
		hl.changed = false
	}

	return err
}

func (hl *hostlist) Bytes() []byte {
	s := condenseNewlines(strings.Join(hl.lines, "\n"))
	return []byte(s)
}

func (hl *hostlist) Contains(a, b string) (bool, error) {
	var ip, hostname string

	if net.ParseIP(a) == nil && net.ParseIP(b) == nil {
		return false, fmt.Errorf("neither %s or %s is a valid IP address", a, b)
	}

	hostname = b
	ip = a

	if net.ParseIP(a) == nil {
		hostname = a
		ip = b
	}

	for _, line := range hl.lines {
		if line == fmt.Sprintf("%s\t%s", ip, hostname) {
			return true, nil
		}
	}

	return false, nil
}

func (hl *hostlist) Add(a, b string) error {
	var ip, hostname string

	if net.ParseIP(a) == nil && net.ParseIP(b) == nil {
		return fmt.Errorf("neither %s or %s is a valid IP address", a, b)
	}

	hostname = b
	ip = a

	if net.ParseIP(a) == nil {
		hostname = a
		ip = b
	}

	hl.lines = append(hl.lines, fmt.Sprintf("%s\t%s", ip, hostname))
	hl.changed = true
	return nil
}

func (hl *hostlist) Remove(thing string) error {
	deletes := []int{}
	for i, line := range hl.lines {
		if containsPart(line, thing) {
			deletes = append(deletes, i)
		}
	}

	for _, i := range reverse(deletes) {
		hl.lines = append(hl.lines[:i], hl.lines[i+1:]...)
	}

	hl.changed = true
	return nil
}

func (hl *hostlist) Comment(thing string) error {
	for i, line := range hl.lines {
		if containsPart(line, thing) {
			hl.lines[i] = "#" + line
		}
	}

	hl.changed = true
	return nil
}

func (hl *hostlist) Uncomment(thing string) error {
	for i, line := range hl.lines {
		if containsPart(line, thing) {
			hl.lines[i] = strings.TrimLeft(line, "#")
		}
	}

	hl.changed = true
	return nil
}

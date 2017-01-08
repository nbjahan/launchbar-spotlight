package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func search(q string, byname bool, searchPath string) ([]string, error) {
	var c = "eval /usr/bin/mdfind"

	onlyin := !byname && strings.Contains(q, "-onlyin ")

	if searchPath != "/" {
		if !onlyin {
			searchPath = normalizePath(searchPath)
			c += fmt.Sprintf(` -onlyin '"%s"'`, searchPath)
		}
	}

	if byname {

		for _, k := range strings.Fields(q) {
			c += ` -name '"` + k + `"'`
		}
	} else {
		c += " " + q
	}

	cmd := exec.Command("/bin/sh", "-c", c)
	data, err := cmd.CombinedOutput()
	pb.Logger.Printf("executing %q", cmd.Args)
	if err != nil {
		pb.Logger.Printf("Error executing %q %s", cmd.Args, string(data))
	}
	if len(data) == 0 {
		return nil, nil
	}
	return strings.Split(string(data), "\n"), err
}

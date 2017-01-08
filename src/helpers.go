package main

import (
	"os"
	"path"
	"path/filepath"

	sjson "github.com/bitly/go-simplejson"
)

var InDev string

func debug(s string, ss ...string) {
	if InDev != "" && pb.Config.GetBool("indev") {
		ss = append([]string{s}, ss...)
		for _, out := range ss {
			nice := out
			js, err := sjson.NewJson([]byte(out))
			if err == nil {
				b, err := js.EncodePretty()
				if err == nil {
					nice = string(b)
				}
			}
			pb.Logger.Println(string(nice))
		}
	}
}

func normalizePath(pth string) string {
	if len(pth) > 1 && pth[:2] == "~/" {
		pth = path.Join("$HOME", pth[2:])
	} else if pth == "~" {
		pth = "$HOME"
	}
	pth, _ = filepath.Abs(os.ExpandEnv(pth))
	return pth
}

func isDir(p string) bool {
	if p == "" {
		return false
	}
	p = normalizePath(p)
	f, err := os.Open(p)
	if err != nil {
		return false
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return false
	}
	if info.IsDir() {
		return true
	}
	return false
}

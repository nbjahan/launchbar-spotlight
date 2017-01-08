package main

import (
	"fmt"
	"path"
	"strings"
	"time"

	. "github.com/nbjahan/go-launchbar"
)

var pb *Action
var start = time.Now()

var funcs = map[string]Func{
	"showView": func(c *Context) Items {
		view := c.Input.FuncArg()
		if view == "config" {
			return c.Action.GetView("config").Render()
		}
		c.Action.ShowView(view)
		return Items(nil)
	},
	"setPath": func(c *Context) {
		if !c.Config.GetBool("in-ispath") {
			return
		}

		path := c.Input.FuncArg()
		c.Config.Set("path", path)
		c.Config.Set("in-ispath", false)
		c.Action.ShowView("main")
	},
	"resetPath": func(c *Context) {
		path := c.Input.FuncArg()
		c.Config.Set("path", path)
		c.Config.Set("in-path", path)
		c.Config.Set("in-ispath", false)
		c.Action.ShowView("main")
	},
	"search": func(c *Context) Items {
		items := &Items{}
		args := c.Input.FuncArgs()
		q, searchPath, bynameStr := args[0], args[1], args[2]
		if searchPath == "" {
			searchPath = pb.Config.GetString("in-path")
		}
		if q == "" {
			q = c.Input.String()
		}

		byname := false
		bynameTitle := ""
		if bynameStr == "" {
			byname = c.Action.IsControlKey()
		} else if bynameStr == "true" {
			byname = true
		}
		if byname {
			bynameTitle = "by Name "
		}

		backItem := NewItem("Back")
		backItem.SetIcon("BackTemplate")
		backItem.Run("showView", "main")
		backItem.SetActionRunsInBackground(true)
		backItem.SetAction("spotlight")

		searchPath = normalizePath(searchPath)
		files, err := search(q, byname, searchPath)
		if err != nil {
			for _, f := range files {
				//FIXME: move this in search func
				f = strings.TrimSpace(f)
				if f != "" {
					items.Add(NewItem(f).SetIcon("/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources/UnknownFSObjectIcon.icns"))
				}
			}
			items.Add(backItem)
		} else if files == nil {
			pb.Logger.Println("Not found")
			items.Add(NewItem("not found").SetIcon("/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources/UnknownFSObjectIcon.icns"))
			if byname || !strings.Contains(q, "-onlyin ") {
				oldpath := searchPath
				for {
					newpath, _ := path.Split(oldpath)
					if oldpath == "/" {
						break
					}
					if newpath == oldpath {
						newpath = "/"
					}
					oldpath = newpath
					items.Add(NewItem(fmt.Sprintf("Search %sin %q", bynameTitle, newpath)).SetSubtitle("query: "+q).SetIcon("SearchTemplate").Run("search", q, newpath, byname).SetActionReturnsItems(true).SetActionRunsInBackground(false).SetAction("spotlight"))
				}
			}
			items.Add(backItem)
		} else {
			for _, f := range files {
				//FIXME: move this in search func
				f = strings.TrimSpace(f)
				if f != "" {
					// TODO: 6.1 nightly 6101 Has a bug with multiple selection ⇧<- , ⌃Y
					items.Add(NewItem("").SetPath(f))
				}
			}
		}
		return *items
	},
}

func init() {
	pb = NewAction("Spotlight Search", ConfigValues{
		"actionDefaultScript": "spotlight",
		"debug":               true,
		"path":                "/",
	})
}

func main() {
	pb.Init(funcs)

	debug("in", pb.Input.Raw())
	in := pb.Input

	if !in.IsLiveFeedback() {
		if in.IsPaths() && len(in.Paths()) == 1 && isDir(in.Paths()[0]) {
			pb.Config.Set("in-path", in.Paths()[0])
			pb.Config.Set("in-ispath", true)
			pb.ShowView("main")
			return
		}
	} else if !in.IsObject() {
		firstInput := strings.TrimSpace(strings.Split(in.String(), " ")[0])
		if isDir(firstInput) {
			pb.Config.Set("in-path", firstInput)
			pb.Config.Set("in-ispath", true)
			out := pb.GetView("main").Compile()
			debug("out: (livefeedback is dir)", out)
			fmt.Println(out)
			return
		} else {
			// FIXME: pb.Config.Set("in-ispath", false)
		}
	}

	if !pb.Config.GetBool("in-ispath") {
		path := pb.Config.GetString("path")
		pb.Config.Set("in-path", path)
	}

	out := pb.Run()
	debug("out:", out)
	fmt.Println(out)
}

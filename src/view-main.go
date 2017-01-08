package main

import (
	"fmt"
	"strings"

	. "github.com/nbjahan/go-launchbar"
)

func init() {
	var i *Item
	v := pb.NewView("main")

	i = v.NewItem("Search in Spotlight")
	i.SetIcon("SearchTemplate")
	i.SetActionRunsInBackground(false)
	i.SetActionReturnsItems(true)
	i.SetRender(func(c *Context) {
		q := c.Input.String()
		onlyin := !c.Action.IsControlKey() && strings.Contains(q, "-onlyin ")
		inpath := ""
		if !onlyin {
			inpath = c.Config.GetString("in-path")
		}
		if c.Input.IsEmpty() || q == inpath {
			if inpath != "" {
				c.Self.SetSubtitle(fmt.Sprintf("in: %q", inpath))
			}
		} else {
			if inpath != "" && strings.HasPrefix(q, inpath) {
				q = strings.TrimSpace(q[len(inpath):])
			}
			c.Self.SetSubtitle(fmt.Sprintf("query: %q in: %q", q, inpath))
		}
		c.Self.Run("search", q, inpath, "")
	})

	i = v.NewItem("Set Search Path")
	i.SetIcon("/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources/SmartFolderIcon.icns")
	i.SetSubtitle(pb.Config.GetString("in-path"))
	i.SetMatch(func(c *Context) bool {
		inpath := c.Config.GetString("in-path")
		c.Self.SetSubtitle(inpath)
		c.Self.Run("setPath", inpath)
		return inpath != c.Config.GetString("path")
	})

	i = v.NewItem("Reset Search Path")
	i.SetIcon("/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources/SmartFolderIcon.icns")
	i.SetSubtitle(pb.Config.GetString("path"))
	i.SetMatch(func(c *Context) bool {
		path := c.Config.GetString("path")
		c.Self.SetSubtitle(path)
		c.Self.Run("resetPath", path)
		return c.Config.GetString("in-path") != path
	})

	i = v.NewItem("CTRL+Enter to Search by Name").SetActionRunsInBackground(false).SetIcon("/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources/ToolbarInfo.icns").Run("showView", "main")

	i = v.NewItem("Spotlight Search: Preferences")
	i.SetActionRunsInBackground(false)
	i.SetActionReturnsItems(true)
	i.SetIcon("/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources/ToolbarAdvanced.icns")
	i.SetOrder(99)
	i.Run("showView", "config")
}

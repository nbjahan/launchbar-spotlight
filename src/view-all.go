package main

import (
	"fmt"
	"time"

	. "github.com/nbjahan/go-launchbar"
)

func init() {
	var i *Item
	v := pb.NewView("*")

	i = v.NewItem(fmt.Sprintf("Executed in: %v", time.Since(start)))
	i.SetSubtitle(fmt.Sprintf("%0.3f seconds", float64(time.Since(start))/float64(time.Second)))
	i.SetIcon("com.apple.Safari:ToolbarWebInspectorTemplate")
	i.SetMatch(MatchIfTrueFunc(pb.Config.GetBool("debug")))
	i.SetOrder(9998)
}

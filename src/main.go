package main

import g "github.com/andrei-galkin/imdoto/gsearch"
import b "github.com/andrei-galkin/imdoto/bsearch"
import im "github.com/andrei-galkin/imdoto/imdoto"

func main() {
	setting := im.GetSetting()
	if setting.Engine == "google" {
		g.Download(setting)
	} else {
		b.Download(setting)
	}
}

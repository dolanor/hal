package main

import (
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/aquilax/go-wakatime"
)

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	currDesk, err := ewmh.CurrentDesktopGet(X)
	if err != nil {
		panic(err)
	}
	fmt.Println("curr desk:", currDesk)

	names, err := ewmh.DesktopNamesGet(X)
	if err != nil {
		panic(err)
	}
	fmt.Println("desk names:", names)

	err = ewmh.DesktopNamesSet(X, []string{"desk0", "logcollector", "desk2", "desk3", "desk4", "desk5", "desk6"})
	if err != nil {
		panic(err)
	}

	names, err = ewmh.DesktopNamesGet(X)
	if err != nil {
		panic(err)
	}
	fmt.Println("desk names:", names)

	if int(currDesk) > len(names) {
		log.Println("name not set for desk:", currDesk)
		return
	}
	projectName := names[currDesk]
	fmt.Println("curr desk name:", projectName)

	wktr := wakatime.NewBasicTransport(wkKey)
	wk := wakatime.New(wktr)

	hbResp, err := wk.PostHeartbeat("current", wakatime.HeartbeatItem{
		Entity:   "gnome workspace",
		Project:  projectName,
		Time:     float32(time.Now().UTC().Unix()),
		Language: "Go",
	})
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("hbResp:", hbResp)

}

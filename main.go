package main

import (
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/xgb/screensaver"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/aquilax/go-wakatime"
	"github.com/coreos/go-systemd/v22/login1"
)

func main() {
	apiKey, err := loadWakatimeAPIKey()
	if err != nil {
		log.Println("api_key not found in wakatime config file. Have you set it?")
		panic(err)
	}
	fmt.Println("desk names:", names)

	err = ewmh.DesktopNamesSet(X, deskNames)
	if err != nil {
		panic(err)
	}

	names, err = ewmh.DesktopNamesGet(X)
	if err != nil {
		panic(err)
	}
	fmt.Println("desk names:", names)

	wktr := wakatime.NewBasicTransport(wkKey)
	wk := wakatime.New(wktr)

	l, err := login1.New()
	if err != nil {
		panic(err)
	}
	defer l.Close()
	users, err := l.ListUsers()
	if err != nil {
		panic(err)
	}
	fmt.Println("users:", users)

	sessions, err := l.ListSessions()
	if err != nil {
		panic(err)
	}
	fmt.Println("sessions:", sessions)

	sess, err := l.GetActiveSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("sess:", sess)

	user, err := l.GetSessionUser(sess)
	if err != nil {
		panic(err)
	}
	fmt.Println("user:", user)

	for range time.Tick(time.Second * 5) {
		info, err := screensaver.QueryInfo(X.Conn(), xproto.Drawable(rootWin)).Reply()
		if err != nil {
			log.Fatal(err)
		}
		inactiveFor := time.Duration(info.MsSinceUserInput) * time.Millisecond
		fmt.Println("Inactive for", inactiveFor)

		if inactiveFor > time.Second*30 {
			println("ignoring this activity")
			continue
		}

		sess, err := l.GetActiveSession()
		if err != nil {
			panic(err)
		}
		fmt.Println("sess:", sess)
		continue
		currDesk, err := ewmh.CurrentDesktopGet(X)
		if err != nil {
			panic(err)
		}
		if int(currDesk) > len(names) {
			log.Println("name not set for desk:", currDesk)
			return
		}

		projectName := names[currDesk]
		fmt.Println("curr desk name:", projectName)
		fmt.Println("curr desk:", currDesk)
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
}

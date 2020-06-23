package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/xgb/screensaver"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/aquilax/go-wakatime"
)

func printUsage() {
	fmt.Println(`hal: your friendly robot friend that records everything you do on task
	usage: hal PROJECT_NAME

PROJECT_NAME: the project name you want to be recorded in WakaTime.

needs a WakaTime set up on the machine.`)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	projectName := os.Args[1]

	apiKey, err := loadWakatimeAPIKey()
	if err != nil {
		log.Println("api_key not found in wakatime config file. Have you set it?")
		panic(err)
	}

	wktr := wakatime.NewBasicTransport(apiKey)
	wk := wakatime.New(wktr)

	hal, err := NewHal(projectName, time.Second*30, wk)
	if err != nil {
		panic(err)
	}

	for range time.Tick(time.Second * 5) {

		if !hal.isUserActive() {
			// let's check another tick
			continue
		}

		working, err := hal.isHumanDoingTheAssignedProject()
		if err != nil {
			log.Println(err)
		}
		if !working {
			continue
		}

		hbResp, err := hal.wakatime.PostHeartbeat("current", wakatime.HeartbeatItem{
			Entity:  "Gnome workspace",
			Project: hal.projectName,
			Time:    float32(time.Now().UTC().Unix()),
		})
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("hbResp:", hbResp)
	}
}

func (h *Hal) isHumanDoingTheAssignedProject() (bool, error) {
	currDesk, err := ewmh.CurrentDesktopGet(h.xutil)
	if err != nil {
		return false, err
	}

	currViewport, err := ewmh.DesktopViewportGet(h.xutil)
	if err != nil {
		return false, err
	}
	vp := currViewport[0]
	println("desk id", currDesk, "viewport: (", vp.X, ", ", vp.Y, ")")

	return currDesk == h.projectDesktopID &&
		vp == h.projectDesktopViewport, nil
}

type Hal struct {
	projectName            string
	projectDesktopID       uint
	projectDesktopViewport ewmh.DesktopViewport
	humanActivityThreshold time.Duration

	xutil   *xgbutil.XUtil
	rootWin xproto.Window

	wakatime *wakatime.WakaTime
}

func NewHal(projectName string, humanActivityThreshold time.Duration, wk *wakatime.WakaTime) (*Hal, error) {
	X, err := xgbutil.NewConn()
	if err != nil {
		return nil, err
	}

	err = screensaver.Init(X.Conn())
	if err != nil {
		return nil, err
	}

	setup := xproto.Setup(X.Conn())
	rootWin := setup.DefaultScreen(X.Conn()).Root
	drw := xproto.Drawable(rootWin)
	screensaver.SelectInput(X.Conn(), drw, screensaver.EventNotifyMask)

	currDesk, err := ewmh.CurrentDesktopGet(X)
	if err != nil {
		return nil, err
	}

	currViewport, err := ewmh.DesktopViewportGet(X)
	if err != nil {
		return nil, err
	}

	return &Hal{
		projectName:            projectName,
		projectDesktopID:       currDesk,
		projectDesktopViewport: currViewport[0],
		humanActivityThreshold: humanActivityThreshold,

		xutil:   X,
		rootWin: rootWin,

		wakatime: wk,
	}, nil
}

func (h *Hal) isUserActive() bool {
	info, err := screensaver.QueryInfo(h.xutil.Conn(), xproto.Drawable(h.rootWin)).Reply()
	if err != nil {
		log.Fatal(err)
		return true
	}
	inactiveFor := time.Duration(info.MsSinceUserInput) * time.Millisecond
	//fmt.Println("Inactive for", inactiveFor)

	return inactiveFor <= h.humanActivityThreshold
}

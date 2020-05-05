package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
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

	err = ewmh.DesktopNamesSet(X, []string{"desk0", "desk1", "desk2", "desk3", "desk4", "desk5", "desk6"})
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
	fmt.Println("curr desk name:", names[currDesk])

}

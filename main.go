package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"log"
	"os/exec"
)

func play() {
	cmd := exec.Command("mpv", "https://www.youtube.com/watch?v=5qap5aO4i9A", "--no-video")
	err := cmd.Start()

	defer cmd.Wait()

	if err != nil {
		log.Fatal(err)
	}
}

func onReady() {
	systray.SetTitle("lofibar")
	systray.SetTooltip("pipe youtube audio from your menubar")
	mPlay := systray.AddMenuItem("Play/Pause", "play/pause")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mPlay.ClickedCh
		fmt.Println("playing!")
		play()
		fmt.Println("Finished quitting")
	}()
	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()
}

func onExit() {
	// clean up here
}

func main() {
	systray.Run(onReady, onExit)
}

package main

import (
	_ "embed"
	"log"
	"os/exec"

	"github.com/getlantern/systray"
)

//go:embed lofi.ico
var icon []byte
var cmd *exec.Cmd

func play() {
	cmd = exec.Command("mpv", "https://www.youtube.com/watch?v=5qap5aO4i9A", "--no-video")
	err := cmd.Start()

	defer cmd.Wait()

	if err != nil {
		log.Fatal(err)
	}
}

func onReady() {
	systray.SetIcon(icon)
	systray.SetTitle("lofibar")
	systray.SetTooltip("pipe youtube audio from your menubar")
	mPlay := systray.AddMenuItem("Play/Pause", "play/pause")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mPlay.ClickedCh
		play()
	}()
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	cmd.Process.Kill()
}

func main() {
	systray.Run(onReady, onExit)
}

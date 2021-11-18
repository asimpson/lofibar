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
var isPlaying bool

func playPause() {
	if isPlaying {
		cmd.Process.Kill()
		isPlaying = false
	} else {
		cmd = exec.Command("mpv", "https://www.youtube.com/watch?v=5qap5aO4i9A", "--no-video")
		err := cmd.Start()

		if err != nil {
			log.Fatal(err)
		}
		isPlaying = true
	}
}

func onReady() {
	systray.SetIcon(icon)
	systray.SetTitle("lofibar")
	systray.SetTooltip("pipe youtube audio from your menubar")

	go func() {
		mPlay := systray.AddMenuItem("Play/Pause", "play/pause")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
		for {
			select {
			case <-mPlay.ClickedCh:
				playPause()
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {
	cmd.Process.Kill()
}

func main() {
	isPlaying = false
	systray.Run(onReady, onExit)
}

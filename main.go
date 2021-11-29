package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/getlantern/systray"
)

//go:embed lofi.ico
var icon []byte
var cmd *exec.Cmd
var isPlaying bool

func playPause(url string) {
	if isPlaying {
		cmd.Process.Kill()
		isPlaying = false
	} else {
		cmd = exec.Command("ffmpeg", "-i", url, "-f", "pulse", "lofi")
		err := cmd.Start()

		if err != nil {
			log.Fatal(err)
		}
		isPlaying = true
	}
}

func parseYT() (streamURL string) {
	type StreamPayload struct {
		StreamingData struct {
			HlsManifestUrl   string
			ExpiresInSeconds string
		}
	}

	resp, err := http.Get("https://www.youtube.com/watch?v=5qap5aO4i9A")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "var ytInitialPlayerResponse") {
			jSSlice := strings.Split(s.Text(), "var ytInitialPlayerResponse = ")
			jsonBytes := []byte(strings.TrimSuffix(jSSlice[1], ";"))
			var stream = &StreamPayload{}

			err := json.Unmarshal(jsonBytes, stream)

			if err != nil {
				log.Fatal(err)
			}

			streamURL = stream.StreamingData.HlsManifestUrl
		}
	})

	return streamURL
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
				url := parseYT()
				playPause(url)
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

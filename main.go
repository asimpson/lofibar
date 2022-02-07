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
	"github.com/shirou/gopsutil/v3/process"
)

//go:embed lofi.ico
var icon []byte
var cmd *exec.Cmd
var isPlaying bool

func quit() {
	pid, err := process.NewProcess(int32(cmd.Process.Pid))

	if err != nil {
		log.Fatal(err)
	}

	children, _ := pid.Children()

	if len(children) != 0 {
		for _, c := range children {
			err = c.Kill()

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	cmd.Process.Kill()
}

func playPause(url string) {
	if isPlaying {
		quit()
		isPlaying = false
	} else {
		cmd = exec.Command("ffplay", url, "-nodisp")
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
	systray.SetTooltip("pipe lofi beats audio from your menubar")

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
				quit()
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

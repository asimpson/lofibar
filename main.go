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

type beats struct {
	cmd       *exec.Cmd
	isPlaying bool
}

func (b *beats) quit() {
	pid, err := process.NewProcess(int32(b.cmd.Process.Pid))

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

	b.cmd.Process.Kill()
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

func (b *beats) onReady() {
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
				b.playPause(url)
			case <-mQuit.ClickedCh:
				b.quit()
				systray.Quit()
			}
		}
	}()
}

func (b *beats) onExit() {}

func main() {
	b := beats{isPlaying: false}

	systray.Run(b.onReady, b.onExit)
}

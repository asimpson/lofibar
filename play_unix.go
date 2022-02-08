//go:build !windows
// +build !windows

package main

import (
	"log"
	"os/exec"
)

func (b *beats) playPause(url string) {
	if b.isPlaying {
		b.quit()
		b.isPlaying = false
	} else {
		b.cmd = exec.Command("ffplay", url, "-nodisp")
		err := b.cmd.Start()

		if err != nil {
			log.Fatal(err)
		}
		b.isPlaying = true
	}
}

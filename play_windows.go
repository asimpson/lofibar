//go:build windows
// +build windows

package main

import (
	"log"
	"os/exec"
	"syscall"
)

func (b *beats) PlayPause(url string) {
	if b.isPlaying {
		b.quit()
		b.isPlaying = false
	} else {
		b.cmd = exec.Command("ffplay", url, "-nodisp")
		b.cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
		err := b.cmd.Start()

		if err != nil {
			log.Fatal(err)
		}
		b.isPlaying = true
	}
}

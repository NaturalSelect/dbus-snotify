package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"

	log "github.com/sirupsen/logrus"
)

const (
	paplayPath        = "/usr/bin/paplay"
	defaultSoundPath  = "./message.ogg"
	signalBufferSize  = 64
	maxInFlightEvents = 16
)

type Notifier struct {
	soundPath string
	soundChan chan interface{}
	done      sync.WaitGroup
	stop      int32
}

func NewNotifyer(soundPath string) (n *Notifier) {
	n = &Notifier{
		soundPath: soundPath,
		soundChan: make(chan interface{}, maxInFlightEvents),
		done:      sync.WaitGroup{},
		stop:      0,
	}
	n.startMonitor()
	n.startSoundPlayer()
	return
}

func (n *Notifier) startMonitor() {
	n.done.Add(1)
	go func() {
		defer n.done.Done()
		cmd := exec.Command("dbus-monitor", "path='/org/freedesktop/Notifications',interface='org.freedesktop.Notifications',member='Notify'")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return
		}
		defer stdout.Close()
		if err := cmd.Start(); err != nil {
			fmt.Println("Error starting dbus-monitor:", err)
			return
		}
		reader := bufio.NewReader(stdout)
		for atomic.LoadInt32(&n.stop) != 1 {
			lineBuf, _, err := reader.ReadLine()
			if err != nil {
				log.Errorf("Failed to read stdout, %v", err)
				continue
			}
			line := string(lineBuf)
			if strings.Contains(line, "member=Notify") {
				n.soundChan <- 1
			}
		}
		defer close(n.soundChan)
		cmd.Process.Kill()
		if err = cmd.Wait(); err != nil {
			log.Errorf("Failed to wait process stop, %v", err)
			return
		}
	}()
}

func (n *Notifier) playSound() {
	soundCmd := exec.Command(paplayPath, n.soundPath)
	if err := soundCmd.Start(); err != nil {
		log.Errorf("Failed to start paplay process, %v", err)
		return
	}
	if err := soundCmd.Wait(); err != nil {
		log.Errorf("Failed to wait paplay process %v", err)
	}
}

func (n *Notifier) startSoundPlayer() {
	n.done.Add(1)
	go func() {
		defer n.done.Done()
		for atomic.LoadInt32(&n.stop) != 1 {
			_, ok := <-n.soundChan
			if !ok {
				return
			}
			n.playSound()
		}
	}()
}

func (n *Notifier) Stop() {
	atomic.StoreInt32(&n.stop, 1)
	n.done.Wait()
}

func main() {
	soundPath := defaultSoundPath
	if len(os.Args) > 1 {
		soundPath = os.Args[1]
	}
	log.Infof("Using sound file %v", soundPath)
	if err := syscall.Access(soundPath, syscall.F_OK); err != nil {
		log.Fatalf("Failed to access sound file %v, %v", soundPath, err)
		return
	}
	sigChan := make(chan os.Signal, signalBufferSize)
	defer close(sigChan)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Reset(os.Interrupt, syscall.SIGTERM)
	n := NewNotifyer(soundPath)
	// NOTE: wait for stop signal
	<-sigChan
	log.Infof("Process stop!")
	n.Stop()
}

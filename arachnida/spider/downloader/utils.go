package downloader

import (
	"fmt"
	"time"
)

type Spinner struct {
	stopChan chan struct{}
}

// StartSpinner runs a spinning animation in a background goroutine
func StartSpinner(message string) *Spinner {
	s := &Spinner{
		stopChan: make(chan struct{}),
	}

	// Dynamic frames for the animation
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		
		i := 0
		for {
			select {
			case <-s.stopChan:
				// Clear the line when stopping
				fmt.Print("\r\033[K") 
				return
			case <-ticker.C:
				// \r moves cursor to start, \033[K clears the line text ahead
				fmt.Printf("\r\033[K%s %s ", frames[i%len(frames)], message)
				i++
			}
		}
	}()

	return s
}

func (s *Spinner) Stop() {
	close(s.stopChan)
	// Brief pause to allow the terminal cursor to reset cleanly
	time.Sleep(50 * time.Millisecond) 
}
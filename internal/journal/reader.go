package journal

import (
	"bufio"
	"bytes"
	"context"
	"log"
	"os/exec"
	"strconv"
)

type Reader struct {
	unit  string
	lines chan string
	ctx   context.Context
	cancel context.CancelFunc
}

func NewReader(unit string) *Reader {
	ctx, cancel := context.WithCancel(context.Background())
	return &Reader{
		unit:   unit,
		lines:  make(chan string, 100),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (r *Reader) Lines() <-chan string {
	return r.lines
}

func (r *Reader) Stop() {
	r.cancel()
}

func (r *Reader) Start() error {
	cmd := exec.CommandContext(r.ctx, "journalctl", "-u", r.unit, "-f", "-n", "0", "-o", "short-iso")
	cmd.Stderr = nil

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			select {
			case r.lines <- line:
			case <-r.ctx.Done():
				return
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Journal scanner error: %v", err)
		}
		cmd.Wait()
	}()

	return nil
}

func (r *Reader) ReadHistory(days int) ([]string, error) {
	cmd := exec.Command("journalctl", "-u", r.unit, "--since", strconv.Itoa(days)+" days ago", "-o", "short-iso")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var result []string
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result, scanner.Err()
}

func (r *Reader) ReadHistorySince(days int) ([]string, error) {
	cmd := exec.Command("journalctl", "-u", r.unit, "--since", strconv.Itoa(days)+" days ago", "-o", "cat")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var result []string
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result, scanner.Err()
}
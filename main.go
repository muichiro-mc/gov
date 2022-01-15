// for the most part this is just a small script to read the blockchain,
// and then pipes a command to a running minecraft server.
package main

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"
)

const (
	BanCmd       = "ban %s\n"
	PardonCmd    = "pardon %s\n"
	WhitelistCmd = "whitelist %s\n"
	stdinFD      = "/proc/%d/fd/0"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Printf("missing minecraft server pid")
		return
	}

	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("unable to parse pid: %v\n", err)
		return
	}

	c := ListenForServerKilled(pid)

	f, err := os.OpenFile(fmt.Sprintf(stdinFD, pid), os.O_APPEND, fs.ModeNamedPipe)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return
	}

	f.Write()
	f.ReadAt

	_ = f

	select {
	case <-c:
		fmt.Printf("watched pid killed")
		return
	}
}

func ListenForServerKilled(pid int) <-chan struct{} {

	c := make(chan struct{})
	go func() {
		p, err := os.FindProcess(pid)
		if err != nil {
			fmt.Printf("could not find %d: %v\n", pid, err)
			return
		}
		p.Wait()
		c <- struct{}{}
	}()

	return c
}

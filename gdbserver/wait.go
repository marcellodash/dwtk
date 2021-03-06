package gdbserver

import (
	"context"
	"fmt"
	"os"

	"github.com/dwtk/dwtk/debugwire"
	"github.com/dwtk/dwtk/internal/wait"
)

func waitForDwOrGdb(ctx context.Context, dw *debugwire.DebugWIRE, conn *tcpConn) ([]byte, error) {
	nctx, cancel := context.WithCancel(ctx)

	sigGdb := make(chan bool)
	sigDw := make(chan bool)

	go func() {
		if err := wait.ForFd(nctx, conn.Fd, sigGdb); err != nil {
			fmt.Fprintf(os.Stderr, "error: gdbserver: gdb: %s\n", err)
		}
	}()

	go func() {
		if err := dw.Wait(nctx, sigDw); err != nil {
			fmt.Fprintf(os.Stderr, "error: gdbserver: debugwire: %s\n", err)
		}
	}()

	var (
		err    error
		packet []byte
	)

	select {
	case <-ctx.Done():
		packet = []byte("S00")
	case <-sigGdb:
		cancel()
		packet = []byte("S02")
	case <-sigDw:
		cancel()
		packet = []byte("S05")
		err = dw.RecvBreak()
	}

	return packet, err
}

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/slavash/throt"
)

// **********************************************************************
// THIS IS STRAIGHTFORWARD IMPLEMENTATION TO DEMONSTRATE THE LIBRARY ONLY
// PLEASE DON'T CONSIDER IT AS A CODE EXAMPLE
// **********************************************************************

const (
	exitCmd          = "exit"
	defaultConnLimit = (980933 - burst) / 3 // I want to download my file (size: 980933) in 3 sec
	burst            = 32 * 1024            // using default io.Copy buffer size as the allowed burst
)

var (
	connLimit     int64
	globalLimiter *throt.Limiter
)

func main() {

	connLimit = defaultConnLimit // byte/sec

	// set bandwidth limit per server
	globalLimiter = throt.NewLimiter(980933, burst)

	ctx, cancel := context.WithCancel(context.Background())

	//ctx = context.WithValue(ctx, serverRateLimiterKey, globalLimiter)
	//ctx = context.WithValue(ctx, connectionRateLimitKey, connLimit)

	defer cancel()

	endPoint := "localhost:7777"
	l, err := net.Listen("tcp4", endPoint)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if e := l.Close(); e != nil {
			fmt.Printf("faield to shotdown the server: %s\n", e)
		}
	}()

	fmt.Printf("Listening on %s\n", l.Addr())

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handleConnection(ctx, c)
	}
}

func handleConnection(ctx context.Context, c net.Conn) {
	defer func() {
		if e := c.Close(); e != nil {
			fmt.Printf("faield to close connection: %s\n", e)
		}
	}()

	fmt.Printf("client connected from %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		cmd := strings.TrimSpace(netData)
		fmt.Printf("Received command: %s\n", cmd)

		if cmd == exitCmd {
			fmt.Printf("Closing connection for the client %s\n", c.RemoteAddr().String())
			break
		}

		// example of changing limits in runtime (applies to all existing connections)
		if len(cmd) > 5 && cmd[:4] == "setl" {
			limit, err := strconv.Atoi(cmd[5:])
			if err != nil {
				fmt.Printf("invalid limit value: %s\n", cmd[5:])
				_, _ = fmt.Fprintf(c, "invalid limit value: %s\n", cmd[5:])
				continue
			}
			if limit == 0 {
				limit = defaultConnLimit
			}

			setConnectionLimit(limit)

			fmt.Printf("Rate limit changed to %d\n", limit)
		}

		if len(cmd) > 4 && cmd[:3] == "get" {
			fileName := cmd[4:]
			err := serveFile(ctx, c, fileName)
			if err != nil {
				fmt.Printf("failed to serve data: %s\n", err)
				_, _ = fmt.Fprintf(c, "failed to serve data: %s\n", err)
			}
		}
	}
}

func setConnectionLimit(limit int) {
	atomic.StoreInt64(&connLimit, int64(limit))
}

func serveFile(ctx context.Context, c io.Writer, name string) error {

	var sent int64
	defer func(start time.Time) {
		fmt.Printf("Sent %d bytes in %.3fs\n", sent, time.Since(start).Seconds())
	}(time.Now())

	fd, err := os.Open(name)
	if err != nil {
		return err
	}

	// set bandwidth limit per connection
	connLimiter := throt.NewLimiter(int(connLimit), burst)

	r1 := throt.NewReader(ctx, fd)
	r1.ApplyLimit(connLimiter)

	r2 := throt.NewReader(ctx, r1)
	r2.ApplyLimit(globalLimiter)

	sent, err = io.Copy(c, r2)

	// The same may be done with writer:
	//writer := throt.NewWriter(ctx, c)
	//writer.ApplyLimit(connLimiter, globalLimiter)
	//sent, err = io.Copy(writer, fd)

	if err != nil {
		return err
	}

	return nil
}

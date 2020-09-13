package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/slavash/throt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	exitCmd    = "exit"
	ioCopyBuff = 32786 // use your buffer size instead
)

var (
	connLimit     int64
	globalLimiter *throt.Limiter
)

func main() {

	connLimit = 196500 // byte/sec

	// set bandwidth limit per server
	globalLimiter = throt.NewLimiter(196500*10, ioCopyBuff)

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

		// changing limits in runtime (applies to all existing connections)
		if len(cmd) > 5 && cmd[:4] == "setl" {
			limit, err := strconv.ParseInt(cmd[5:], 10, 32)
			if err != nil {
				fmt.Printf("invalid limit value: %s [%s]\n", err, cmd)
				_, _ = fmt.Fprintf(c, "invalid limit value: %s [%s]\n", err, cmd)
			}
			if connLimit == 0 {
				limit = connLimit
			}

			atomic.StoreInt64(&connLimit, limit)

			fmt.Printf("Rate limit changed to %d\n", limit)
			_, _ = fmt.Fprintf(c, "Rate limit changed to %d\n", limit)
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

func serveFile(ctx context.Context, c net.Conn, name string) error {

	var sent int64
	defer func(start time.Time) {
		fmt.Printf("Sent %d bytes in %03fs\n", sent, time.Since(start).Seconds())
	}(time.Now())

	fd, err := os.Open(name)
	if err != nil {
		return err
	}

	// set bandwidth limit per connection
	connLimiter := throt.NewLimiter(connLimit, ioCopyBuff)

	reader := throt.NewReader(ctx, fd)
	reader.ApplyLimits(connLimiter, globalLimiter)
	sent, err = io.Copy(c, reader)

	// The same may be done with writer:
	//writer := throt.NewWriter(ctx, c)
	//writer.ApplyLimits(connLimiter, globalLimiter)
	//sent, err = io.Copy(writer, fd)

	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func scanPort(ctx context.Context, host string, port int, timeout time.Duration) bool {
	address := net.JoinHostPort(host, strconv.Itoa(port))
	d := net.Dialer{
		Timeout: timeout,
	}

	conn, err := d.DialContext(ctx, "tcp", address)

	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func main() {
	startTime := time.Now()

	hostPtr := flag.String("host", "127.0.0.1", "Target host")
	startPtr := flag.Int("start", 1, "Starting port")
	endPtr := flag.Int("end", 1024, "Ending port")
	workersPtr := flag.Int("workers", 10, "Number of worker gorutine")
	timeout := flag.Duration("timeout", 500*time.Millisecond, "Time for one port scan")
	flag.Parse()

	ctx, interrupt := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer interrupt()

	var wg sync.WaitGroup

	ports := make(chan int, 100)
	results := make(chan int)
	openPorts := make([]int, 0)

	for i := 1; i <= *workersPtr; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range ports {
				if scanPort(ctx, *hostPtr, port, *timeout) {
					results <- port
				}
			}

		}()
	}

	go func() {
		defer close(ports)
		for i := *startPtr; i <= *endPtr; i++ {

			select {
			case <-ctx.Done():
				return
			case ports <- i:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		openPorts = append(openPorts, res)
	}
	slices.Sort(openPorts)
	fmt.Println(time.Since(startTime))
	fmt.Println("Open ports: ", openPorts)

}

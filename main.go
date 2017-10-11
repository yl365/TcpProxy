// TcpProxy project main.go
package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	go CheckHosts()
	time.Sleep(2 * time.Second)

	lis, err := net.Listen("tcp", ALLHOST.Listen)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer lis.Close()
	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Printf("accept err:%v\n", err)
			continue
		}
		go handle(conn)
	}
}

func handle(sconn net.Conn) {
	defer sconn.Close()

	ip, ok := getBackendHost(sconn.RemoteAddr())
	if !ok {
		return
	}
	dconn, err := net.Dial("tcp", ip)
	if err != nil {
		fmt.Printf("conn %v ERR:%v\n", ip, err)
		return
	}
	fmt.Printf("NEW CONN: %s-->%s OK!\n", sconn.RemoteAddr(), dconn.RemoteAddr())
	defer dconn.Close()

	makePair(sconn, dconn)
}
func IsDisconnect(err error) bool {
	if err != nil {
		if nerr, ok := err.(net.Error); ok {
			if nerr.Timeout() || nerr.Temporary() {
				return false
			} else {
				return true
			}
		} else {
			return true
		}
	}

	return false
}

func makePair(sconn net.Conn, dconn net.Conn) {

	var wg sync.WaitGroup
	ch1 := make(chan bool, 1)
	ch2 := make(chan bool, 1)
	wg.Add(2)
	go copyContent(dconn, sconn, &wg, ch1, ch2)
	go copyContent(sconn, dconn, &wg, ch2, ch1)
	wg.Wait()
}

func copyContent(from net.Conn, to net.Conn, wg *sync.WaitGroup, done chan bool, otherDone chan bool) {
	var err error = nil
	var data []byte = make([]byte, 5120)
	var nr, nw, nWCnt int = 0, 0, 0

	for {
		select {
		case <-otherDone:
			wg.Done()
			fmt.Printf("otherDone!\n")
			return
		default:

			from.SetReadDeadline(time.Now().Add(time.Second * 1))
			nr, err = from.Read(data)

			if IsDisconnect(err) {
				fmt.Printf("from.Read Disconnect!\n")
				done <- true
				wg.Done()
				return
			}

			if nr > 0 {
				to.SetWriteDeadline(time.Now().Add(time.Second * 1))
				nWCnt = 0
				nw = 0

				for nWCnt < nr {
					nw, err = to.Write(data[nWCnt:nr])

					if IsDisconnect(err) {
						fmt.Printf("to.Write Disconnect!\n")
						done <- true
						wg.Done()
						return
					}
					if nw > 0 {
						nWCnt += nw
					}
				}
			}
		}
	}
}

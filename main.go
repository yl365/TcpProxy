// TcpProxy project main.go
package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

func main() {
	go CheckHosts()

	lis, err := net.Listen("tcp", ALLHOST.Listen)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer lis.Close()
	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println("accept err:%v\n", err)
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

	var wg sync.WaitGroup
	ch1 := make(chan bool, 1)
	ch2 := make(chan bool, 1)
	wg.Add(2)
	go copyContent(dconn, sconn, &wg, ch1, ch2)
	go copyContent(sconn, dconn, &wg, ch2, ch1)
	wg.Wait()

	// 也可以直接使用io.Copy, 缺点是一方网络断开另一方不知道
	//	ExitChan := make(chan bool, 1)
	//	go copyData(dconn, sconn, ExitChan)
	//	go copyData(sconn, dconn, ExitChan)
	//	<-ExitChan

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

func copyContent(from net.Conn, to net.Conn, wg *sync.WaitGroup, done chan bool, otherDone chan bool) {
	var err error = nil
	var data []byte = make([]byte, 5120)
	var read int = 0
	var write int = 0
	var writeTotal int = 0

	for {
		select {
		case <-otherDone:
			wg.Done()
			fmt.Printf("otherDone!\n")
			return
		default:

			from.SetReadDeadline(time.Now().Add(time.Second * 2))
			read, err = from.Read(data)

			if IsDisconnect(err) {
				fmt.Printf("from.Read Disconnect!\n")

				done <- true
				wg.Done()
				return
			}

			if read > 0 {
				to.SetWriteDeadline(time.Now().Add(time.Second * 2))
				writeTimes = 0
				writeTotal = 0
				write = 0

				for writeTotal < read {
					write, err = to.Write(data[writeTotal:read])

					if IsDisconnect(err) {
						fmt.Printf("to.Write Disconnect!\n")
						done <- true
						wg.Done()
						return
					}
					if write > 0 {
						writeTotal += write
					}
				}
			}
		}
	}
}

func copyData(dconn net.Conn, sconn net.Conn, Exit chan bool) {

	n, err := io.Copy(dconn, sconn)

	if err != nil {
		fmt.Printf("%s-->%s err: %v\n", sconn.RemoteAddr(), dconn.RemoteAddr(), err)
		Exit <- true
	} else {
		fmt.Printf("%s-->%s transfer byte: %d\n", sconn.RemoteAddr(), dconn.RemoteAddr(), n)
	}

}

package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

func getBackendHost(sAddr net.Addr) (string, bool) {
	lock.RLock()
	defer lock.RUnlock()

	sIP := strings.Split(sAddr.String(), ":")[0]
	IP := net.ParseIP(sIP).To4()
	nIP := int(IP[0]) + int(IP[1])*255 + int(IP[2])*255*255 + int(IP[3])*255*255*255
	hash := nIP % 1000
	fmt.Printf("sIP=%s, nIP=%d, hash=%d\n", sIP, nIP, hash)

	for _, group := range ALLHOST.AllHost {

		if hash >= group.Min && hash <= group.Max {

			if ALLHOST.Mode == "master" { //优选主,如果主不可用则选从

				for _, host := range group.Hosts {
					if host.Status == 0 {
						return host.IP, true
					}
				}

			} else if ALLHOST.Mode == "hash" { //优选hash对应的,如果hash分配的不可用, 用随机分配

				hostNum := len(group.Hosts)
				r := hash % hostNum
				if group.Hosts[r].Status == 0 {
					return group.Hosts[r].IP, true
				} else { //如果hash分配的不可用, 用随机分配
					for l := 0; l < hostNum*2; l++ {
						r := rand.Intn(hostNum)
						if group.Hosts[r].Status == 0 {
							return group.Hosts[r].IP, true
						}
					}
				}
			} else if ALLHOST.Mode == "rand" {

				hostNum := len(group.Hosts)
				for l := 0; l < hostNum*2; l++ {
					r := rand.Intn(hostNum)
					if group.Hosts[r].Status == 0 {
						return group.Hosts[r].IP, true
					}
				}
			}
		}
	}

	return "", false
}

func CheckHosts() {

	for {
		for i, group := range ALLHOST.AllHost {
			//fmt.Printf("1#: %d--->%v\n", i, group)

			for ii, host := range group.Hosts {
				//fmt.Printf("\t2#: %d--->%v\n", ii, host)

				conn, err := net.DialTimeout("tcp", host.IP, 2*time.Second)

				if err != nil {
					//fmt.Printf("conn %s: err=%v\n", host.IP, err)
					lock.Lock()
					ALLHOST.AllHost[i].Hosts[ii].Status -= 1
					lock.Unlock()
				} else if host.Status != 0 {
					lock.Lock()
					ALLHOST.AllHost[i].Hosts[ii].Status = 0
					lock.Unlock()
				}

				if err == nil {
					conn.Close()
				}

			}
		}

		fmt.Printf("\nCheckHosts AllHost: %+v\n", ALLHOST)
		time.Sleep(300 * time.Second)
	}

}

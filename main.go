package main

import (
	"fmt"
	"net"
	"sort"
)

func worker(ports chan int, result chan int) {
	for p := range ports {
		address := fmt.Sprintf("20.194.168.28:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			result <- -p
			continue
		}
		conn.Close()
		result <- p
	}
}

func main() {
	ports := make(chan int, 100)
	result := make(chan int)
	var openPorts []int
	var closePorts []int

	//要检验的端口数
	var totalPorts int = 1024

	for i := 0; i < cap(ports); i++ {
		go worker(ports, result)
	}

	go func() {
		for i := 1; i <= totalPorts && i <= 65535; i++ {
			ports <- i
		}
	}()

	for i := 1; i <= totalPorts && i <= 65535; i++ {
		port := <-result
		if port > 0 {
			openPorts = append(openPorts, port)
		} else {
			closePorts = append(closePorts, -port)
		}
	}

	close(ports)
	close(result)

	sort.Ints(openPorts)
	sort.Ints(closePorts)

	fmt.Println("open ports:", openPorts)
	fmt.Println("close ports:", closePorts)
}

// start := time.Now()
// var wg sync.WaitGroup
// for i := 10; i < 20; i++ {
// 	wg.Add(1)
// 	go func(j int) {
// 		defer wg.Done()

// 		address := fmt.Sprintf("20.194.168.28:%d", j)
// 		conn, err := net.Dial("tcp", address)
// 		if err != nil {
// 			fmt.Printf("%s err\n", address)
// 			return
// 		}
// 		conn.Close()
// 		fmt.Printf("%s ok\n", address)
// 	}(i)
// }
// wg.Wait()
// end := time.Since(start).Seconds()
// fmt.Printf("time: %f\n", end)

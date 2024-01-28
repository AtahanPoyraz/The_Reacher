package Scripts

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type WebInfo struct {
	TargetIP  string
	PortsList []string
	Timeout   time.Duration
}

var (
	wg sync.WaitGroup
)

func PORTS(format string, args ...interface{}) []string {
	var port_var string
	if len(args) > 0 {
		port_var = fmt.Sprintf(format, args...)
	} else {
		fmt.Println("No arguments provided for PORTS function.")
		os.Exit(1)
	}

	if strings.Contains(port_var, ",") {
		ports_list := strings.Split(port_var, ",")

		for _, port := range ports_list {
			_, err_c := strconv.Atoi(port)
			if err_c != nil {
				fmt.Printf("Invalid Port: %s\n", port)
				os.Exit(0)
			}
		}

		return ports_list
	} else if strings.Contains(port_var, "-") {
		port_min_and_max := strings.Split(port_var, "-")

		port_min, err := strconv.Atoi(port_min_and_max[0])
		if err != nil {
			fmt.Printf("Invalid Port Min Range: %s\n", port_min_and_max[0])
			os.Exit(1)
		}

		port_max, err := strconv.Atoi(port_min_and_max[1])
		if err != nil {
			fmt.Printf("Invalid Port Max Range: %s\n", port_min_and_max[1])
			os.Exit(1)
		}

		var ports_temp_list []string

		for p_min := port_min; p_min <= port_max; p_min++ {
			port_str := strconv.Itoa(p_min)
			ports_temp_list = append(ports_temp_list, port_str)
		}

		return ports_temp_list
	} else if strings.Contains(port_var, ":") {
		port_range := strings.Split(port_var, ":")
		start, err := strconv.Atoi(port_range[0])
		if err != nil {
			fmt.Printf("Invalid Port Start Range: %s\n", port_range[0])
			os.Exit(1)
		}

		end, err := strconv.Atoi(port_range[1])
		if err != nil {
			fmt.Printf("Invalid Port End Range: %s\n", port_range[1])
			os.Exit(1)
		}

		var ports_temp_list []string

		for p := start; p <= end; p++ {
			ports_temp_list = append(ports_temp_list, strconv.Itoa(p))
		}

		return ports_temp_list
	}

	_, err := strconv.Atoi(port_var)
	if err != nil {
		fmt.Printf("Invalid Port: %s\n", port_var)
		os.Exit(1)
	}

	return []string{port_var}
}

func StartScanning(webInfo WebInfo) {
	for _, port := range webInfo.PortsList {
		wg.Add(1)
		go ScanPort(webInfo.TargetIP, port, webInfo.Timeout)
	}

	wg.Wait()
}

func ScanPort(targetIP, port string, timeoutTCP time.Duration) {
	defer wg.Done()

	d := net.Dialer{Timeout: timeoutTCP}
	_, err := d.Dial("tcp", targetIP+":"+port)
	if err != nil {
		if add_err, ok := err.(*net.AddrError); ok && add_err.Timeout() {
			return
		} else if add_err, ok := err.(*net.OpError); ok && strings.TrimSpace(add_err.Err.Error()) == "bind: An operation on a socket could not be performed because the system lacked sufficient buffer space or because a queue was full." {
			time.Sleep(timeoutTCP + (3000 * time.Millisecond))
			_, err_ae := d.Dial("tcp", targetIP+":"+port)
			if err_ae != nil && add_err.Timeout() {
				return
			}
		} else {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		return
	}

	fmt.Printf("%s[+]%s Port %s | TCP : OPEN\n", "\x1b[1;32m", "\x1b[1;0m", port)
}

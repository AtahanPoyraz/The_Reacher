package Scripts

import (
    "net"
    "fmt"
    "time"
)

type ScanResult struct {
    Port  int
    State string
}

func ScanPorts(Protocol, Host string, startPort, endPort int) []ScanResult {
    var results []ScanResult

    for port := startPort; port <= endPort; port++ {
        result := ScanPort(Protocol, Host, port)
        if result.State == "Open" {
            results = append(results, result)
        }
    }

    return results
}

func ScanPort(Protocol, Host string, Port int) ScanResult {
    result := ScanResult{Port: Port}
    address := fmt.Sprintf("%s:%d", Host, Port)
    conn, err := net.DialTimeout(Protocol, address, 60*time.Second)

    if err != nil {
        return result
    }

    defer conn.Close()

    result.State = "Open"
    return result
}

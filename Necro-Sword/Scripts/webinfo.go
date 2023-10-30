package Scripts

import (
    "fmt"
    "net"
    "net/url"
    "net/http"
    "crypto/tls"
    "strings"
    "io/ioutil"
)

type INFO struct {
    URL         string
    Host        string
}

func GetInfo(Url string) (*INFO, error) {
    u, err := url.Parse(Url)
    if err != nil {
        return nil, err
    }

    host, _, err := net.SplitHostPort(u.Host)
    if err != nil {
        host = u.Host
    }

    return &INFO{
        URL:  Url,
        Host: host,
    }, nil
}

func (i *INFO) Control() {
    fmt.Println("Gathering information for:", i.URL)
    i.GetServerInfo()
    i.GetDNSInfo()
    i.GetHTTPInfo()
    i.GetRobotsTxt()
    i.GetSSLInfo()
}

func (i *INFO) GetServerInfo() {
    fmt.Printf("%s[+]%s Getting Server Info..\n" , "\x1b[1;32m", "\x1b[1;0m")

    ipAddress, err := net.LookupHost(i.Host)
    if err != nil {
        fmt.Println("IP address lookup error:", err)
        return
    }

    conn, err := net.Dial("tcp", i.Host + ":80")
    if err != nil {
        fmt.Println("Dial error:", err)
        return
    }
    defer conn.Close()

    _, err = conn.Write([]byte("HEAD / HTTP/1.1\r\nHost: " + i.URL + "\r\n\r\n"))
    if err != nil {
        fmt.Println("Write error:", err)
        return
    }

    buffer := make([]byte, 1024)
    _, err = conn.Read(buffer)
    if err != nil {
        fmt.Println("Read error:", err)
        return
    }

    response := string(buffer)
    serverInfo := strings.Split(strings.Split(response, "Server: ")[1], "\r\n")[0]

    fmt.Println("IP Address:", ipAddress[0])
    fmt.Println("Server Version:", serverInfo)
}

func (i *INFO) GetDNSInfo() {
    fmt.Printf("%s[+]%s Getting DNS Info..\n" , "\x1b[1;32m", "\x1b[1;0m")

    addresses, err := net.LookupHost(i.Host)
    if err != nil {
        fmt.Println("DNS lookup error:", err)
        return
    }

    for _, address := range addresses {
        fmt.Println("DNS Info:", address)
    }
}

func (i *INFO) GetHTTPInfo() {
    fmt.Printf("%s[+]%s Getting HTTP Info..\n" , "\x1b[1;32m", "\x1b[1;0m")
    resp, err := http.Get(i.URL)
    if err != nil {
        fmt.Println("HTTP request error:", err)
        return
    }
    defer resp.Body.Close()

    fmt.Println("HTTP Status Code:", resp.Status)
}

func (i *INFO) GetRobotsTxt() {
    fmt.Printf("%s[+]%s Getting Robots.txt..\n" , "\x1b[1;32m", "\x1b[1;0m")
    robotsURL := i.URL + "/robots.txt"
    resp, err := http.Get(robotsURL)
    if err != nil {
        fmt.Println("HTTP request error for robots.txt:", err)
        return
    }
    defer resp.Body.Close()

    robotsContent, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading robots.txt content:", err)
        return
    }

    fmt.Println("Robots.txt Content:")
    fmt.Println(string(robotsContent))
}

func (i *INFO) GetSSLInfo() {
    fmt.Printf("%s[+]%s Getting SSL Info..\n" , "\x1b[1;32m", "\x1b[1;0m")

    conn, err := tls.Dial("tcp", i.Host+":443", &tls.Config{
        InsecureSkipVerify: true,
    })
    if err != nil {
        fmt.Println("SSL connection error:", err)
        return
    }
    defer conn.Close()

    state := conn.ConnectionState()
    if len(state.PeerCertificates) > 0 {
        cert := state.PeerCertificates[0]
        fmt.Println("Common Name (CN):", cert.Subject.CommonName)
        fmt.Println("Organization:", cert.Subject.Organization)
        fmt.Println("Organizational Unit:", cert.Subject.OrganizationalUnit)
        fmt.Println("Country:", cert.Subject.Country)
        fmt.Println("Locality:", cert.Subject.Locality)
        fmt.Println("Province:", cert.Subject.Province)
        fmt.Println("Street Address:", cert.Subject.StreetAddress)
        fmt.Println("Postal Code:", cert.Subject.PostalCode)
        fmt.Println("Serial Number:", cert.SerialNumber)
        fmt.Println("Issuer:", cert.Issuer.CommonName)
        fmt.Println("Valid From:", cert.NotBefore)
        fmt.Println("Valid Until:", cert.NotAfter)
        fmt.Println("Signature Algorithm:", cert.SignatureAlgorithm)
        fmt.Println("Public Key Algorithm:", cert.PublicKeyAlgorithm)

        fmt.Println("Certificate Extensions:")
        for _, ext := range cert.Extensions {
            fmt.Printf("%s: %v\n", ext.Id, ext.Value)
        }
        fmt.Println("Key Usage:")
        fmt.Printf("Digital Signature: %v\n", cert.KeyUsage)

        fmt.Println("Serial Number:", cert.SerialNumber.String())
        fmt.Println("Issuer:", cert.Issuer.CommonName)
        fmt.Println("Certificate Policies:")
        for _, policy := range cert.PolicyIdentifiers {
            fmt.Printf("Policy : %s\n", policy.String())
        }
    }
}

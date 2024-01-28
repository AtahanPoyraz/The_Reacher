package Scripts

import (
    "fmt"
    "net"
    "net/url"
    "net/http"
    "crypto/tls"
    "encoding/json"
    "strings"
    "io/ioutil"
)

type WhoisResponse struct {
	Server         string   `json:"server"`
	Name           string   `json:"name"`
	IDNName        string   `json:"idnName"`
	Status         string   `json:"status"`
	Nameserver     []string `json:"nameserver"`
	IPS            string   `json:"ips"`
	Created        string   `json:"created"`
	Changed        string   `json:"changed"`
	Expires        string   `json:"expires"`
	Registered     bool     `json:"registered"`
	DNSSEC         string   `json:"dnssec"`
	WhoisServer    string   `json:"whoisserver"`
	Contacts       Contacts `json:"contacts"`
	Registrar      string   `json:"registrar"`
	RawData        []string `json:"rawdata"`
	Network        string   `json:"network"`
	Exception      string   `json:"exception"`
	ParsedContacts bool     `json:"parsedContacts"`
	Template       struct {
		WhoisNIC string `json:"whois.nic.tr"`
	} `json:"template"`
}

type Contacts struct {
	Owner []struct {
		Handle       interface{} `json:"handle"`
		Type         interface{} `json:"type"`
		Name         interface{} `json:"name"`
		Organization string      `json:"organization"`
		Email        interface{} `json:"email"`
		Address      string      `json:"address"`
		Zipcode      interface{} `json:"zipcode"`
		City         string      `json:"city"`
		State        interface{} `json:"state"`
		Country      interface{} `json:"country"`
		Phone        interface{} `json:"phone"`
		Fax          interface{} `json:"fax"`
		Created      interface{} `json:"created"`
		Changed      interface{} `json:"changed"`
	} `json:"owner"`
}

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
    i.GetWhoIsInfo()
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

func (i *INFO) GetWhoIsInfo() {
    fmt.Printf("%s[+]%s Getting Whois Info..\n", "\x1b[1;32m", "\x1b[1;0m")

    whoisURL := "https://whoisjson.com/api/v1/whois"
    queryParams := url.Values{"domain": {i.Host}}
    fullURL := whoisURL + "?" + queryParams.Encode()

    // Kullanıcı https:// ile başladıysa bunu kaldır
    if strings.HasPrefix(i.Host, "https://") {
        i.Host = strings.TrimPrefix(i.Host, "https://")
    }

    req, err := http.NewRequest("GET", fullURL, nil)
    if err != nil {
        fmt.Println("Error creating request:", err)
        return
    }

    req.Header.Set("Authorization", "Token=48ba9979767cc7d0e8a6467c2b1709a2363c677b8c1eb0953e4f81c64116a0a6")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error making request:", err)
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading response body:", err)
        return
   }

    var whoisResponse WhoisResponse
    err = json.Unmarshal(body, &whoisResponse)
    if err != nil {
        fmt.Println("Error decoding JSON:", err)
        return
    }

    // Print each field of the struct
    fmt.Println("Server:", whoisResponse.Server)
    fmt.Println("Name:", whoisResponse.Name)
    fmt.Println("IDN Name:", whoisResponse.IDNName)
    fmt.Println("Status:", whoisResponse.Status)
    fmt.Println("Nameserver:", whoisResponse.Nameserver)
    fmt.Println("IPS:", whoisResponse.IPS)
    fmt.Println("Created:", whoisResponse.Created)
    fmt.Println("Changed:", whoisResponse.Changed)
    fmt.Println("Expires:", whoisResponse.Expires)
    fmt.Println("Registered:", whoisResponse.Registered)
    fmt.Println("DNSSEC:", whoisResponse.DNSSEC)
    fmt.Println("Whois Server:", whoisResponse.WhoisServer)

    // Print contacts information
    fmt.Println("Contacts:")
    for _, owner := range whoisResponse.Contacts.Owner {
        fmt.Println("  Organization:", owner.Organization)
        fmt.Println("  Address:", owner.Address)
        fmt.Println("  City:", owner.City)
        fmt.Println("  Country:", owner.Country)
    }

    fmt.Println("Registrar:", whoisResponse.Registrar)
    fmt.Println("Raw Data:", whoisResponse.RawData)
    fmt.Println("Network:", whoisResponse.Network)
    fmt.Println("Exception:", whoisResponse.Exception)
    fmt.Println("Parsed Contacts:", whoisResponse.ParsedContacts)
    fmt.Println("Template - Whois NIC:", whoisResponse.Template.WhoisNIC)
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

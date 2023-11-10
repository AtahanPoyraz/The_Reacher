package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/AtahanPoyraz/Necro-Sword/Scripts"
)

var (
	platform = strings.ToLower(runtime.GOOS)
)
var (
	ans	 	 string
	url 	 string
	IP   	 string
	path 	 string
	sport 	 int
	eport  	 int
	tout	 int
)

func main() {
	Clear()
	fmt.Printf(`%s
     ▐ ▄ ▄▄▄ . ▄▄· ▄▄▄            .▄▄ · ▄▄▌ ▐ ▄▌      ▄▄▄  ·▄▄▄▄  
    •█▌▐█▀▄.▀·▐█ ▌▪▀▄ █·▪         ▐█ ▀. ██· █▌▐█▪     ▀▄ █·██▪ ██ 
    ▐█▐▐▌▐▀▀▪▄██ ▄▄▐▀▀▄  ▄█▀▄     ▄▀▀▀█▄██▪▐█▐▐▌ ▄█▀▄ ▐▀▀▄ ▐█· ▐█▌
    ██▐█▌▐█▄▄▌▐███▌▐█•█▌▐█▌.▐▌    ▐█▄▪▐█▐█▌██▐█▌▐█▌.▐▌▐█•█▌██. ██ 
    ▀▀ █▪ ▀▀▀ ·▀▀▀ .▀  ▀ ▀█▄▀▪     ▀▀▀▀  ▀▀▀▀ ▀▪ ▀█▄▀▪.▀  ▀▀▀▀▀▀•
    _______________________________________________________________
                                                   by Atahan Poyraz
    [%s01%s] IP TRACKER
    [%s02%s] WEBSITE INFORMATION                   
    [%s03%s] PORT SCAN     
    [%s04%s] URL EXPLORATION                  

> %s`, "\x1b[1;31m","\x1b[1;32m","\x1b[1;31m","\x1b[1;32m","\x1b[1;31m","\x1b[1;32m","\x1b[1;31m","\x1b[1;32m","\x1b[1;31m","\x1b[1;0m")
	fmt.Scan(&ans)

	if ans == "1" || ans == "01" {
		fmt.Printf("%s[+]%s Target IP: ", "\x1b[1;32m", "\x1b[1;0m")
		fmt.Scan(&IP)
		fmt.Printf("%s--------------------------------%s\n", "\x1b[1;37m", "\x1b[1;0m")
		data, err := Scripts.Tracking(IP)
		if err != nil {
			fmt.Printf("%s[!]%s Error: %e", "\x1b[1;93m", "\x1b[1;0m", err)
			return
		}
		fmt.Printf("%s[?]%s Status: %s\n", "\x1b[1;93m", "\x1b[1;0m", data.Status)
		fmt.Printf("%s[01]%s Country: %s\n", "\x1b[1;93m", "\x1b[1;0m", data.Country)
		fmt.Printf("%s[02]%s Country Code: %s\n", "\x1b[1;93m", "\x1b[1;0m", data.CountryCode)
		fmt.Printf("%s[03]%s Region: %s\n", "\x1b[1;93m", "\x1b[1;0m", data.Region)
		fmt.Printf("%s[04]%s Region Name: %s\n", "\x1b[1;93m", "\x1b[1;0m", data.RegionName)
		fmt.Printf("%s[05]%s City: %s\n", "\x1b[1;93m", "\x1b[1;0m", data.City)
		fmt.Printf("%s[06]%s Zip: %s\n", "\x1b[1;93m", "\x1b[1;0m", data.Zip)
		fmt.Printf("%s[07]%s Lat: %f\n", "\x1b[1;93m", "\x1b[1;0m", data.Lat)
		fmt.Printf("%s[08]%s Lon: %f\n", "\x1b[1;93m", "\x1b[1;0m", data.Lon)
		fmt.Printf("%s[09]%s Timezone: %s\n", "\x1b[1;93m", "\x1b[1;0m", data.Timezone)
		fmt.Printf("%s[10]%s ISP: %s\n", "\x1b[1;93m", "\x1b[1;0m", data.ISP)
		fmt.Printf("%s[11]%s Org: %s\n", "\x1b[1;93m", "\x1b[1;0m", data.Org)
		fmt.Printf("%s[12]%s AS: %s\n", "\x1b[1;93m", "\x1b[1;0m", data.AS)
		fmt.Printf("%s[13]%s Query: %s\n", "\x1b[1;93m", "\x1b[1;0m", data.Query)
		
	} else if ans == "2" || ans == "02" {
		fmt.Printf("%s[+]%s Target URL: ", "\x1b[1;32m", "\x1b[1;0m")
		fmt.Scan(&url)
		fmt.Printf("%s--------------------------------%s\n", "\x1b[1;37m", "\x1b[1;0m")

		getInfo, err := Scripts.GetInfo(url)
		if err != nil {
			fmt.Printf("%s[!]%s Error: %e", "\x1b[1;93m", "\x1b[1;0m", err)
			return
		}
    	getInfo.Control()

	} else if ans == "3" || ans == "03" {
		fmt.Printf("%s[+]%s Target IP  : ", "\x1b[1;32m", "\x1b[1;0m")
		fmt.Scan(&IP)
		fmt.Printf("%s[+]%s Start Port : ", "\x1b[1;32m", "\x1b[1;0m")
		fmt.Scan(&sport)
		fmt.Printf("%s[+]%s End Port   : ", "\x1b[1;32m", "\x1b[1;0m")
		fmt.Scan(&eport)
		fmt.Printf("%s[+]%s Time Out   : ", "\x1b[1;32m", "\x1b[1;0m")
		fmt.Scan(&tout)
		fmt.Printf("%s--------------------------------%s\n", "\x1b[1;37m", "\x1b[1;0m")

		Domain := Scripts.WebInfo{
			TargetIP: IP,
			PortsList: Scripts.PORTS("%d-%d", sport, eport),
			Timeout: time.Duration(tout) * time.Second,
		}
		Scripts.StartScanning(Domain)

	} else if ans == "4" || ans == "04" {
		fmt.Printf("%s[+]%s Target URL   : ", "\x1b[1;32m", "\x1b[1;0m")
		fmt.Scan(&url)
		fmt.Printf("%s[+]%s WordList Path: ", "\x1b[1;32m", "\x1b[1;0m")
		fmt.Scan(&path)
		fmt.Printf("%s--------------------------------%s\n", "\x1b[1;37m", "\x1b[1;0m")
		_, _, err := Scripts.StartExploration(url, path)

		if err != nil {
			fmt.Printf("%s[!]%s Error: %e", "\x1b[1;93m", "\x1b[1;0m", err)
			return
		}
	} else {
		fmt.Printf("%s[!]%s Invalid Option", "\x1b[1;93m", "\x1b[1;0m")
		time.Sleep(time.Second * 1)
		defer main()
	}
}

func Clear() {
	var cmd *exec.Cmd
	if platform == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

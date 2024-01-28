package Scripts

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

var (
	dom = []string{}
	suc = []string{}
	fal = []string{}
)

func ReadFile(PATH string) ([]string, error) {
	file, err := os.Open(PATH)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		dom = append(dom, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dom, nil
}

func SiteStatus(url, subdomain string) (string, error) {
	response, err := http.Get(url + subdomain)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	return response.Status, nil
}

func StartExploration(url, path string) ([]string, []string, error) {
	url = fmt.Sprintf("%s/", url)
	domains, err := ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	for _, d := range domains {
		status, err := SiteStatus(url, d)
		if err != nil {
			return nil, nil, err
		}
		if status == "200 OK" {
			fmt.Printf("%s[+]%s %s\n", "\x1b[1;32m", "\x1b[1;0m", (url + d))
			suc = append(suc, d)
		} else {
			fal = append(fal, d)
		}
	}

	return suc, fal, nil
}
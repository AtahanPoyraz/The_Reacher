package Scripts

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
)

var (
    url_ = "http://ip-api.com/json/"
)

type GeoData struct {
    Status      string  `json:"status"`
    Country     string  `json:"country"`
    CountryCode string  `json:"countryCode"`
    Region      string  `json:"region"`
    RegionName  string  `json:"regionName"`
    City        string  `json:"city"`
    Zip         string  `json:"zip"`
    Lat         float64 `json:"lat"`
    Lon         float64 `json:"lon"`
    Timezone    string  `json:"timezone"`
    ISP         string  `json:"isp"`
    Org         string  `json:"org"`
    AS          string  `json:"as"`
    Query       string  `json:"query"`
}

func Tracking(IP string) (GeoData, error) {
    var geoData GeoData

    response, err := http.Get(url_ + IP)
    if err != nil {
        return geoData, err
    }
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return geoData, err
    }

    err = json.Unmarshal(body, &geoData)
    if err != nil {
        return geoData, err
    }

    return geoData, nil
}
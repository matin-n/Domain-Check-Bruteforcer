package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ernestosuarez/itertools"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var TLD = flag.String("tld", "com", "What is the TLD that you want to check?")

func main() {

	// Parse command line parameters
	flag.Parse()

	// Generate Combination

	// combinations of r = 3 elements chosen from iterable
	r := 3 // search for 3 character combination
	iterable := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

	for v := range itertools.CombinationsStr(iterable, r) {
		str := v[0] + v[1] + v[2]

		// Send post request to determine domain availability
		var status = googleCheckDomain(str, *TLD)

		if status != "AVAILABILITY_UNAVAILABLE" {
			fmt.Println(status, str+"."+*TLD)
		}

		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("done")

	/*
		var status = googleCheckDomain(str, "net")

		fmt.Println("Status:", status, str + ".net")

		time.Sleep(500 * time.Millisecond)
	*/
}

func googleCheckDomain(wantedDomain, tld string) string {

	map1 := map[string][]interface{}{
		"domainName": []interface{}{
			map[string]interface{}{
				"sld": wantedDomain,
				"tld": tld,
			},
		},
	}

	b, err := json.Marshal(map1)
	if err != nil {
		panic(err)
	}

	contentReader := bytes.NewReader(b)

	req, _ := http.NewRequest("POST", "https://domains.google.com/v1/Main/FeSearchService/Availability?authuser=0", contentReader)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Origin", "https://domains.google.com")
	req.Header.Set("User-Agent", " Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Referer", "https://domains.google.com/m/registrar/search?searchTerm="+wantedDomain)
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br") // comment out to let the transport automatically handle gzip
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, _ := http.DefaultClient.Do(req)

	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := trimString(string(bodyBytes))

	//fmt.Println(bodyString)

	//domainStatus := gjson.Get(bodyString, "availabilityResponse").String()
	domainStatus := gjson.Get(bodyString, "availabilityResponse.results.result.0.supportedResultInfo.availabilityInfo.availability").String()

	if domainStatus == "AVAILABILITY_AVAILABLE" {

		domainStatus_premiumCheck := gjson.Get(bodyString, "availabilityResponse.results.result.0.supportedResultInfo.purchaseInfo.aftermarketPremium").String()
		if domainStatus_premiumCheck == "AFTERMARKET_PREMIUM_NOT_AFTERMARKET" {
			return "Available"
		} else if domainStatus_premiumCheck == "AFTERMARKET_PREMIUM_FAST_TRANSFER" {
			return "Buy from broker"
		}

	} else if domainStatus == "AVAILABILITY_UNAVAILABLE" {
		return domainStatus
	}

	return domainStatus

}

func trimString(input string) string {
	input = strings.Trim(input, ")]}'")
	input = strings.TrimSpace(input)
	return input
}

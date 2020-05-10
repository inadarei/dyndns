/**
* This is go code that basically is equivalent to this bash command:
* > DDNS_ip="$(curl ifconfig.me/ip 2>/dev/null)" && \
*    curl "https://dynamicdns.park-your-domain.com/update?host=${DDNS_host}&domain=${DDNS_domain}&password=${DDNS_pwd}&ip=${DDNS_ip}"
*/
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/xml"
	"os"
)

type fnBytes func([]byte)

type ResponseFromAPI struct {
	XMLName xml.Name "interface-response";
	Command, IP  string;
	ErrCount, ResponseCount int;
	Responses []Response  `xml:"responses>response"`
	Done string;
	debug string "chardata";
}

type Response struct {
	XMLName xml.Name `xml:"response"`
	ResponseNumber int;
	ResponseString string;
}

func queryAPIError(body []byte) {
	errMsg := string(body)
	panic(fmt.Sprintf("Error API Response: '%s'", errMsg))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
			return value
	}
	return fallback
}

/** queryAPI queries an HTTP API, and performs error handling */
func queryAPI(apiURL string, errHandler fnBytes) []byte {
	res, err := http.Get(apiURL)
	if err != nil {
			panic(err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
			panic(err.Error())
	}

	if res.StatusCode > 399 {
			errHandler(body)
	}

	return body
}
	
func main() {
	domain := getEnv("DDNS_domain", "example.com")
	host := getEnv("DDNS_host", "ssh")
	pwd := getEnv("DDNS_pwd", "")

	ip := string(queryAPI("https://ifconfig.me/ip", queryAPIError))
	
	base := "https://dynamicdns.park-your-domain.com/update?"
	updateURI := base + "host=" + host + "&domain=" + domain + "&password=" + pwd + "&ip=" + ip

	bodyXML := queryAPI(updateURI, queryAPIError)
	var response ResponseFromAPI
	if err := xml.Unmarshal(bodyXML, &response); err != nil {
    panic(err)
	}

	if (response.ErrCount > 0 || response.Done != "true") {
		fmt.Printf("Done: %s \n", response.Done)
		fmt.Printf("updateURI: %s \n", updateURI)
		fmt.Printf("ErrCount: %d \n", response.ErrCount)
		fmt.Printf("error: %s \n", response.Responses[0].ResponseString)
		os.Exit(1)
	} else {
		fmt.Printf("Successfully updated to: %s \n", ip)
	}
	

}


/*

<?xml version="1.0"?>
<interface-response>
	<Command>SETDNSHOST</Command>
	<Language>eng</Language>
	<ErrCount>1</ErrCount>
	<errors>
		<Err1>Domain name not found</Err1>
	</errors>
	<ResponseCount>1</ResponseCount>
	<responses>
		<response>
			<ResponseNumber>316153</ResponseNumber>
			<ResponseString>Validation error; not found; domain name(s)</ResponseString>
		</response>
	</responses>
	<Done>true</Done>
	<debug><![CDATA[]]></debug>
</interface-response>

*/
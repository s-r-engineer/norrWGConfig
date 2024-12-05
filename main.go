package main

import (
	"flag"
	"fmt"
	"os"

	libraryErrors "github.com/s-r-engineer/library/errors"
	libraryIO "github.com/s-r-engineer/library/io"
	libraryNordvpn "github.com/s-r-engineer/library/nordvpn"
)

const template = `[Interface]
Address = 10.5.0.2/32
ListenPort = 51820
PrivateKey = %s
DNS = 103.86.96.100, 103.86.99.100

[Peer]
PublicKey = %s
AllowedIPs = 0.0.0.0/0
Endpoint = %s.%s:51820
`

func main() {
	flag.Parse()
	country := flag.Arg(0)
	host := flag.Arg(1)
	if country == "" || host == "" {
		panic("first country code then host!!!!!")
	}

	if country == "uk" {
		country = "gb"
	}
	token, err := libraryIO.ReadFileToString("./token")
	libraryErrors.Panicer(err)
	countryID, err := libraryNordvpn.GetCountryCode(country)
	libraryErrors.Panicer(err)
	if countryID == -1 {
		panic("no such country " + country)
	}
	privateKey, _ := libraryNordvpn.FetchOwnPrivateKey(token)
	_, _, publicKey, _, _ := libraryNordvpn.FetchServerData(countryID)
	file, err := os.OpenFile(country+".conf", os.O_CREATE|os.O_WRONLY, 0600)
	libraryErrors.Panicer(err)
	defer file.Close()
	_, err = file.Write([]byte(fmt.Sprintf(template, privateKey, publicKey, country, host)))
	libraryErrors.Panicer(err)
}

package main

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"github.com/brianium/mnemonic"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/wemeetagain/go-hdwallet"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	addresses := fetch()
	c := 0
	for {
		m, _ := mnemonic.NewRandom(256, mnemonic.English)

		// Create a master private key
		privateKey := hdwallet.MasterKey([]byte(m.Sentence()))

		// Convert a private key to public key
		publicKey := privateKey.Pub()

		// Get the address
		address := publicKey.Address()

		witnessProg := btcutil.Hash160(publicKey.Serialize())
		addressWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, &chaincfg.MainNetParams)
		if err != nil {
			panic(err)
		}
		bipAddress := addressWitnessPubKeyHash.EncodeAddress()

		if c > 100 {
			c = 0
			fmt.Print(".")
		}

		for _, a := range addresses {
			if strings.HasPrefix(a, "bc1") {
				if bipAddress == a {
					foundBTC(privateKey.String(), publicKey.String(), address)
				}
			} else {
				if address == a {
					foundBTC(privateKey.String(), publicKey.String(), address)
				}
			}
		}
		c++
	}
}

func foundBTC(privateKey string, publicKey string, address string) {
	log.Println("Found some BTC!")
	log.Println(privateKey)
	log.Println(publicKey)
	log.Println(address)
	store(fmt.Sprintf("%s,%s,%s\n", privateKey, publicKey, address))
}

func store(s string) {
	f, err := os.OpenFile("found.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(s); err != nil {
		log.Println(err)
	}
}

func fetch() []string {
	var addresses []string
	log.Println("Attempting to load from file")
	const storeFile = "addresses.bin"
	rd, err := ioutil.ReadFile(storeFile)
	if err != nil {
		log.Println("Downloading new file")
		records, err := downloadFile("https://bitkeys.work/btc_balance_sorted.csv")
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Found %d records\n", len(records))
		// TODO should normalize all the different types of addresses here
		for i, r := range records {
			// Skip headers
			if i > 0 {
				addresses = append(addresses, r[0])
			}
		}
		log.Println("finished processing")
		buf := &bytes.Buffer{}
		err = gob.NewEncoder(buf).Encode(addresses)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(storeFile, buf.Bytes(), 0666)
		if err != nil {
			panic(err)
		}
		return addresses
	}
	err = gob.NewDecoder(bytes.NewReader(rd)).Decode(&addresses)
	if err != nil {
		panic(err)
	}
	log.Printf("Loaded %d records\n", len(addresses))
	return addresses
}

func downloadFile(url string) ([][]string, error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		return nil, err1
	}
	defer resp.Body.Close()

	csvReader := csv.NewReader(resp.Body)
	records, err2 := csvReader.ReadAll()
	if err2 != nil {
		return nil, err2
	}

	return records, nil
}

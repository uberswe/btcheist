package main

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/wemeetagain/go-hdwallet"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const storeFile = "addresses.bin"

func main() {
	addresses := fetch()

	log.Println("Begin processing address loop")
	c := 0
	for {
		m, err2 := hdwallet.GenSeed(256)
		if err2 != nil {
			panic(err2)
		}

		// Create a master private key
		privateKey := hdwallet.MasterKey(m)

		// Convert a private key to public key
		publicKey := privateKey.Pub()

		// Get the address
		address := publicKey.Address()

		// TODO while running in docker we seem to get stuck here
		witnessProg := btcutil.Hash160(publicKey.Serialize())
		addressWitnessPubKeyHash, err2 := btcutil.NewAddressWitnessPubKeyHash(witnessProg, &chaincfg.MainNetParams)
		if err2 != nil {
			panic(err2)
		}
		addressScriptHash, err3 := btcutil.NewAddressScriptHash(witnessProg, &chaincfg.MainNetParams)
		if err3 != nil {
			panic(err3)
		}
		bipAddress := addressWitnessPubKeyHash.EncodeAddress()
		scriptAddress := addressScriptHash.EncodeAddress()

		if c > 10 {
			c = 0
			log.Println(".")
		}

		for _, a := range addresses {
			if bipAddress == a {
				foundBTC(privateKey.String(), publicKey.String(), a)
			} else if scriptAddress == a {
				foundBTC(privateKey.String(), publicKey.String(), a)
			} else if address == a {
				foundBTC(privateKey.String(), publicKey.String(), a)
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

func fetch() (addresses []string) {
	log.Println("Attempting to load from file")
	rd, err := ioutil.ReadFile(storeFile)
	if err != nil {
		log.Println("Downloading new file")
		records, err2 := downloadFile("https://bitkeys.work/btc_balance_sorted.csv")
		if err2 != nil {
			log.Fatalln(err2)
		}
		log.Printf("Found %d records\n", len(records))
		for i, r := range records {
			// Skip headers
			if i > 0 {
				addresses = append(addresses, r[0])
			}
		}
		records = nil
		log.Println("Encoding results")
		buf := &bytes.Buffer{}
		err = gob.NewEncoder(buf).Encode(addresses)
		if err != nil {
			panic(err)
		}
		log.Println("Writing file to disk")
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

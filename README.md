## BTC Heist

Basically a copy of [btc-heist](https://github.com/theden/btc-heist) but written in Go.

This is just for fun if you like playing the lottery, the chances of actually guessing a private key with a balance is so tiny it's not really worth trying.

This is a simple Go program which will download a [CSV dump of all bitcoin addresses with a balance](https://bitkeys.work/download.php) and then randomly generates private keys and checks if they match one of the addresses.

The program will store a bin file to prevent having to re-download the csv file every time the program runs.

Currently there is an issue with this program since it doesn't handle different types of keys. The file that we download should be normalized or we need to generate every different type of address. Some ideas on how this can be done can be taken from https://github.com/modood/btckeygen/blob/master/main.go

Run using `go run main.go`. You can use `go build` first to create an executable as it seems to use less memory.
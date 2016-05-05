package main

import (
	"flag"
	"log"
	"os"
	"snailmail/fileprocessor"
)

func main() {
	log.SetOutput(os.Stdout)

	action := flag.String("a", "send", "creates a zip file and sends it to the given email address. Possible values are: 'send' and 'receive'")
	emailAddress := flag.String("e", "", "email address to send/receive file")
	fileName := flag.String("f", "", "file name to use when sending a file")

	flag.Parse()
	if *emailAddress == "" {
		log.Fatal("Email address is required!!")
	}

	if *action == "send" {
		fileprocessor.ReadFileAndZip(*fileName)
	} else if *action == "receive" {
		fileprocessor.UnzipAndJoin(*fileName)
	}
}

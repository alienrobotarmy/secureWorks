package main

import "fmt"
import "os"
import "secureWorks"
import "flag"
import "encoding/base64"
import "io/ioutil"

func main() {
	fileName := flag.String("c", "", "Config File <required>")
	Help := flag.Bool("h", false, "Help")
	AtId := flag.String("i", "", "Attachment Id <required>")
	Out := flag.String("o", "", "Filename <optional> (Output attachment to file)")
	TicketNumber := flag.String("t", "", "Ticket Number <required>")
	flag.Parse()
	if len(*fileName) == 0 || len(*TicketNumber) == 0 || len(*AtId) == 0 || *Help == true {
		fmt.Fprintf(os.Stderr, "Go Interface for SecureWorks Soap API by Jess Mahan\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Must specify Config file with -c\n")
		os.Exit(0)
	}

	l, err := secureWorks.ReadConfig(*fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	a, err := secureWorks.GetAttachment(l, *TicketNumber, *AtId)

	if len(*Out) == 0 {
		fmt.Printf("Content: %s\nFilename: %s\nmd5Sum: %s\n",
			a.Content,
			a.Filename,
			a.Md5Sum)
	} else {
		d, err := base64.StdEncoding.DecodeString(a.Content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		} else {
			ioutil.WriteFile(*Out, d, 0722)
		}
	}
}

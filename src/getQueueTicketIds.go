package main

import "fmt"
import "os"
import "secureWorks"
import "flag"

func main() {
	fileName := flag.String("c", "", "Config File")
	TicketType := flag.String("t", "", "Ticket Type")
	Limit := flag.Int("l", 25, "Ticket Limit (Max is 500) Default is 25")
	Help := flag.Bool("h", false, "Help")
	flag.Parse()
	if len(*fileName) == 0 || *Help == true {
		fmt.Fprintf(os.Stderr, "Go Interface for SecureWorks Soap API by Jess Mahan\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Must specify Config file with -c\n")
		os.Exit(0)
	}

	if len(*TicketType) == 0 {
		*TicketType = "INCIDENT"
	}

	l, err := secureWorks.ReadConfig(*fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	y, _ := secureWorks.GetQueueTicketIds(l, *TicketType, *Limit)
	for _, v := range y.TicketIds {
		fmt.Printf("%s\n", v)
	}
}

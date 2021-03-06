package main

import "fmt"
import "os"
import "secureWorks"
import "flag"

func main() {
	fileName := flag.String("c", "", "Config File <required>")
	TicketNumber := flag.String("t", "", "Ticket Number <required>")
	Csv := flag.Bool("C", false, "CSV Output")
	Long := flag.Bool("L", false, "Long Output")
	Short := flag.Bool("S", false, "Short Output (don't include work logs)")
	Work := flag.Bool("W", false, "Show Work Logs Only")
	Help := flag.Bool("h", false, "Help")
	flag.Parse()
	if len(*fileName) == 0 || *Help == true || len(*TicketNumber) == 0 {
		fmt.Fprintf(os.Stderr, "Go Interface for SecureWorks Soap API by Jess Mahan\n")
		flag.PrintDefaults()
		os.Exit(0)

	}
	if *Csv == false && *Work == false && *Short == false && *Long == false {
		fmt.Fprintf(os.Stderr, "WARN: No Output Option Selected, using default: Long\n")
		*Long = true
	}

	l, err := secureWorks.ReadConfig(*fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	d, _ := secureWorks.GetTicketDetail(l, *TicketNumber)
	if *Csv == true {
		d.Detail.PrintCsv()
	}
	if *Work == true {
		d.Detail.PrintWorkLogs()
	}
	if *Short == true {
		d.Detail.PrintDetails()
	}
	if *Long == true {
		d.Detail.PrintDetails()
		d.Detail.PrintWorkLogs()
	}
}

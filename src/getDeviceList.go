package main

import "fmt"
import "os"
import "secureWorks"
import "flag"

func main() {
	fileName := flag.String("c", "", "Config File")
	flag.Parse()
	if len(*fileName) == 0 {
		fmt.Fprintf(os.Stderr, "Go Interface for SecureWorks Soap API by Jess Mahan\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Must specify Config file with -c\n")
		os.Exit(0)
	}

	l, err := secureWorks.ReadConfig(*fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	x, err := secureWorks.GetDeviceList(l)
	if err != nil {
		fmt.Printf("Error: %q\n", err)
	}
	fmt.Printf("Id,Device\n")
	for _, v := range x.Devices {
		fmt.Printf("%d,%s\n", v.DeviceId, v.DeviceAlias)
	}

}

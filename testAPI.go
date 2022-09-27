package main

import (
	"fmt"
        bigip "github.com/f5devcentral/go-bigip"
)

// Note: Had to call main something different as the package complained main was in LTM's example
func main() {
	// Connect to the BIG-IP system.
	f5 := bigip.NewSession("10.124.5.244", "443","admin", "chakku@chakku1",nil)
	/*var records  []bigip.DeviceRecord
        records = append(records, bigip.DeviceRecord{Name: "test1.ns.ctc",Address: "3.3.3.3"})
        config := &bigip.Server{
                Name:                     "test1.ns.ctc",
                Datacenter:               "test2",
                Virtual_server_discovery: true,
		Devices:                  records,
        }

	err := f5.CreateGtmserver(config)*/
	var members []string
	//members := make([]string, 1, 4)
	members = append(members,"2.2.2.2")
        config := &bigip.Pool_a{
                Name:    "test1.ns.ctc",
                Members: members,
        }
        err := f5.CreatePool_a(config)

        
        if err != nil {
		fmt.Println(err)
	}
}

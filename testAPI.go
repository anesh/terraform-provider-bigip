package main

import (
	"fmt"
        bigip "github.com/f5devcentral/go-bigip"
)

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
	/*var members []string
	//members := make([]string, 1, 4)
	members = append(members,"Gw-cr-F5-2-lab.ctc:test_vs")
        config := &bigip.Pool_a{
                Name:    "test1.ns.ctc",
                Members: members,
        }
        //err := f5.CreatePool_a(config)
	err := f5.ModifyPool_a("test1.ns.ctc",config)*/
        var records  []bigip.MemberRecord
        records = append(records, bigip.MemberRecord{Name: "/Common/test.ns.ctc"})
        config := &bigip.Wideip_a{
                Name:               "test1.wip.ns.ctc",
                Pools:               records,
        }

        err := f5.CreateWideip_a(config)

        //gw, err := f5.GetWideip_a("/Common/test1.gtm.net")
	//gw,err := f5.GetPool_a("/Common/test1.ns.ctc")
        fmt.Println(err)


        
}

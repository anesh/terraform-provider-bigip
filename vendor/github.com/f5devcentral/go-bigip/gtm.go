package bigip

import "encoding/json"
import "log"


type Datacenters struct {
	Datacenters []Datacenter `json:"items"`
}

type Datacenter struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Contact     string `json:"contact,omitempty"`
	AppService  string `json:"appService,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
	Enabled     bool   `json:"enabled,omitempty"`
	ProberPool  string `json:"proberPool,omitempty"`
}

type Gtmmonitors struct {
	Gtmmonitors []Gtmmonitor `json:"items"`
}

type Gtmmonitor struct {
	Name          string `json:"name,omitempty"`
	Defaults_from string `json:"defaultsFrom,omitempty"`
	Interval      int    `json:"interval,omitempty"`
	Probe_timeout int    `json:"probeTimeout,omitempty"`
	Recv          string `json:"recv,omitempty"`
	Send          string `json:"send,omitempty"`
}

type Servers struct {
	Servers []Server `json:"items"`
}

type Server struct {
	Name                     string
	Datacenter               string
	Monitor                  string
	Virtual_server_discovery bool
	Product                  string
	Devices                  []DeviceRecord
	GTMVirtual_Server        []VSrecord
}

type DeviceRecord struct {
     Name        string `json:"name"`
     Address string `json:"addresses,omitempty"`
}

type serverDTO struct {
	Name                     string `json:"name"`
	Datacenter               string `json:"datacenter,omitempty"`
	Monitor                  string `json:"monitor,omitempty"`
	Virtual_server_discovery bool   `json:"virtual_server_discovery"`
	Product                  string `json:"product,omitempty"`
	Devices                  []DeviceRecord `json:"devices,omitempty"`
	GTMVirtual_Server struct {
		Items []VSrecord `json:"items,omitempty"`
	} `json:"virtualServersReference,omitempty"`
}

func (p *Server) MarshalJSON() ([]byte, error) {
	return json.Marshal(serverDTO{
		Name:                     p.Name,
		Datacenter:               p.Datacenter,
		Monitor:                  p.Monitor,
		Virtual_server_discovery: p.Virtual_server_discovery,
		Product:                  p.Product,
		Devices:                  p.Devices,
		GTMVirtual_Server: struct {
			Items []VSrecord `json:"items,omitempty"`
		}{Items: p.GTMVirtual_Server},
	})
}

func (p *Server) UnmarshalJSON(b []byte) error {
	var dto serverDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}

	p.Name = dto.Name
	p.Datacenter = dto.Datacenter
	p.Monitor = dto.Monitor
	p.Virtual_server_discovery = dto.Virtual_server_discovery
	p.Product = dto.Product
	p.Devices = dto.Devices
	p.GTMVirtual_Server = dto.GTMVirtual_Server.Items
	return nil
}

type ServerAddressess struct {
	Items []ServerAddresses `json:"items,omitempty"`
}

type ServerAddresses struct {
	Name        string `json:"name"`
	Device_name string `json:"deviceName,omitempty"`
	Translation string `json:"translation,omitempty"`
}

type VSrecords struct {
	Items []VSrecord `json:"items,omitempty"`
}

type VSrecord struct {
	Name        string `json:"name"`
	Destination string `json:"destination,omitempty"`
}

type Pool_as struct {
	Pool_as []Pool_a `json:"items"`
}

type Pool_a struct {
	Name                 string   `json:"name,omitempty"`
	Monitor              string   `json:"monitor,omitempty"`
	Load_balancing_mode  string   `json:"load_balancing_mode,omitempty"`
	Max_answers_returned int      `json:"max_answers_returned,omitempty"`
	Alternate_mode       string   `json:"alternate_mode,omitempty"`
	Fallback_ip          string   `json:"fallback_ip,omitempty"`
	Fallback_mode        string   `json:"fallback_mode,omitempty"`
	Members              []string `json:"members,omitempty"`
}

const (
	uriGtm        = "gtm"
	uriServer     = "server"
	uriDatacenter = "datacenter"
	uriGtmmonitor = "monitor"
	uriHttp       = "http"
	uriPool_a     = "pool/a"
)

func (b *BigIP) GetDatacenter(name string) (*Datacenter, error) {
	var datacenter Datacenter
	err, ok := b.getForEntity(&datacenter, uriGtm, uriDatacenter,name)

	if err != nil {
		return nil, err
	}

        if !ok  {
                return nil, nil
        }


	return &datacenter, nil
}

func (b *BigIP) CreateDatacenter(name string) error {
	config := &Datacenter{
		Name: name,
	}
	return b.post(config, uriGtm, uriDatacenter)
}

func (b *BigIP) ModifyDatacenter(name string, config *Datacenter) error {
	return b.put(config,uriGtm, uriDatacenter, name)
}

func (b *BigIP) DeleteDatacenter(name string) error {
	return b.delete(uriGtm, uriDatacenter, name)
}

func (b *BigIP) Gtmmonitors() (*Gtmmonitor, error) {
	var gtmmonitor Gtmmonitor
	err, _ := b.getForEntity(&gtmmonitor, uriGtm, uriGtmmonitor, uriHttp)

	if err != nil {
		return nil, err
	}

	return &gtmmonitor, nil
}
func (b *BigIP) CreateGtmmonitor(name, defaults_from string, interval, probeTimeout int, recv, send string) error {
	config := &Gtmmonitor{
		Name:          name,
		Defaults_from: defaults_from,
		Interval:      interval,
		Probe_timeout: probeTimeout,
		Recv:          recv,
		Send:          send,
	}
	return b.post(config, uriGtm, uriGtmmonitor, uriHttp)
}

func (b *BigIP) ModifyGtmmonitor(*Gtmmonitor) error {
	return b.patch(uriGtm, uriGtmmonitor, uriHttp)
}

func (b *BigIP) DeleteGtmmonitor(name string) error {
	return b.delete(uriGtm, uriGtmmonitor, uriHttp, name)
}

func (b *BigIP) CreateGtmserver(p *Server) error {
	log.Println(" what is the complete payload    ", p)
	return b.post(p, uriGtm, uriServer)
}

//Update an existing policy.
func (b *BigIP) UpdateGtmserver(name string, p *Server) error {
	return b.put(p, uriGtm, uriServer, name)
}

//Delete a policy by name.
func (b *BigIP) DeleteGtmserver(name string) error {
	return b.delete(uriGtm, uriServer, name)
}

func (b *BigIP) GetGtmserver(name string) (*Server, error) {
	var p Server
	err, ok := b.getForEntity(&p, uriGtm, uriServer, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &p, nil
}

func (b *BigIP) CreatePool_a(name, monitor, load_balancing_mode string, max_answers_returned int, alternate_mode, fallback_ip, fallback_mode string, members []string) error {
	config := &Pool_a{
		Name:                 name,
		Monitor:              monitor,
		Load_balancing_mode:  load_balancing_mode,
		Max_answers_returned: max_answers_returned,
		Alternate_mode:       alternate_mode,
		Fallback_ip:          fallback_ip,
		Fallback_mode:        fallback_mode,
		Members:              members,
	}
	log.Println("in poola now", config)
	return b.patch(config, uriGtm, uriPool_a)
}

func (b *BigIP) ModifyPool_a(config *Pool_a) error {
	return b.put(config, uriGtm, uriPool_a)
}

func (b *BigIP) Pool_as() (*Pool_a, error) {
	var pool_a Pool_a
	err, _ := b.getForEntity(&pool_a, uriGtm, uriPool_a)

	if err != nil {
		return nil, err
	}

	return &pool_a, nil
}

// GTMAPool contains information about each gtm/pool/a
type GTMAPool struct {
	Name                      string `json:"name,omitempty"`
	Partition                 string `json:"partition,omitempty"`
	FullPath                  string `json:"fullPath,omitempty"`
	Generation                int    `json:"generation,omitempty"`
	AppService                string `json:"appService,omitempty"`
	Description               string `json:"description,omitempty"`
	Disabled                  bool   `json:"disabled,omitempty"`
	DynamicRatio              string `json:"dynamicRatio,omitempty"`
	Enabled                   bool   `json:"enabled,omitempty"`
	FallbackIP                string `json:"fallbackIp,omitempty"`
	FallbackMode              string `json:"fallbackMode,omitempty"`
	LimitMaxBps               uint64 `json:"limitMaxBps,omitempty"`
	LimitMaxBpsStatus         string `json:"limitMaxBpsStatus,omitempty"`
	LimitMaxConnections       uint64 `json:"limitMaxConnections,omitempty"`
	LimitMaxConnectionsStatus string `json:"limitMaxConnectionsStatus,omitempty"`
	LimitMaxPps               uint64 `json:"limitMaxPps,omitempty"`
	LimitMaxPpsStatus         string `json:"limitMaxPpsStatus,omitempty"`
	LoadBalancingMode         string `json:"loadBalancingMode,omitempty"`
	ManualResume              string `json:"manualResume,omitempty"`
	MaxAnswersReturned        int    `json:"maxAnswersReturned,omitempty"`
	Monitor                   string `json:"monitor,omitempty"`
	TmPartition               string `json:"tmPartition,omitempty"`
	QosHitRatio               int    `json:"qosHitRatio,omitempty"`
	QosHops                   int    `json:"qosHops,omitempty"`
	QosKilobytesSecond        int    `json:"qosKilobytesSecond,omitempty"`
	QosLcs                    int    `json:"qosLcs,omitempty"`
	QosPacketRate             int    `json:"qosPacketRate,omitempty"`
	QosRtt                    int    `json:"qosRtt,omitempty"`
	QosTopology               int    `json:"qosTopology,omitempty"`
	QosVsCapacity             int    `json:"qosVsCapacity,omitempty"`
	QosVsScore                int    `json:"qosVsScore,omitempty"`
	TTL                       int    `json:"ttl,omitempty"`
	VerifyMemberAvailability  string `json:"verifyMemberAvailability,omitempty"`
	MembersReference          struct {
		Link            string `json:"link,omitempty"`
		IsSubcollection bool   `json:"isSubcollection,omitempty"`
	}
}

// AddGTMAPool adds a Pool/A by config to the BIG-IP system.
func (b *BigIP) AddGTMAPool(config *GTMAPool) error {
	return b.post(config, uriGtm, uriPool, string(ARecord))
}

// ModifyGTMAPool adds a Pool/A by config to the BIG-IP system.
func (b *BigIP) ModifyGTMAPool(fullPath string, config *GTMAPool) error {
	return b.put(config, uriGtm, uriPool, string(ARecord), fullPath)
}

// GetGTMAPool get's a Pool/A by name
func (b *BigIP) GetGTMAPool(name string) (*GTMAPool, error) {
	var w GTMAPool

	err, ok := b.getForEntity(&w, uriGtm, uriPool, string(ARecord), name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &w, nil
}


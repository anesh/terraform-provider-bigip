package bigip

import (
	"fmt"
	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceBigipGtmServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipGtmServerCreate,
		Read:   resourceBigipGtmServerRead,
		Update: resourceBigipGtmServerUpdate,
		Delete: resourceBigipGtmServerDelete,
		Exists: resourceBigipGtmServerExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of Server",
				ForceNew:     true,
				ValidateFunc: validateF5NameWithDirectory,
			},

			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
			},
                        
			"virtualServerDiscovery":{
                               Type:     schema.TypeBool,
			       Optional: true,

		        },
                        "Addresses":{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
		},
	}
}

func resourceBigipGtmServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	datacenter := d.Get("datacenter").(string)
	virtualServerDiscovery := d.Get("virtualServerDiscovery").(string)
	log.Println("[INFO] Creating GTM Server " + name)
        config := &bigip.Server{
			Name:                      name,
			Datacenter:                datacenter,
			Virtual_Server_discovery:  virtualServerDiscovery,
		}
	err := client.CreateGtmserver(config)
        if err != nil {
			log.Printf("[ERROR] Unable to Create GTM Server  (%s) (%v)", name, err)
			return err
		}
	d.SetId(name)
	err = resourceBigipGtmServerUpdate(d, meta)
	if err != nil {
		_ = client.DeleteGtmserver(name)
		return err
	}
	return resourceBigipGtmServerRead(d, meta)
}

func resourceBigipGtmServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	_ = d.Set("name", name)
	log.Println("[INFO] Fetching GTM Servers " + name)
	sv, err := client.GetGtmserver(name)
	if err != nil {
		return err
	}
	if sv == nil {
		log.Printf("[WARN] GTM Server (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
        _ = d.Set("name",name)
	_ = d.Set("datacenter",sv.Datacenter)
        _ = d.Set("virtualServerDiscovery",sv.Virtual_server_discovery)

        deviceNames := schema.NewSet(schema.HashString, make([]interface{}, 0, len(sv.Addresses)))
	for _, device := range sv.Addresses {
		FullDeviceName := device.Name
		deviceNames.Add(FullDeviceName)
	}
	if deviceNames.Len() > 0 {
		_ = d.Set("Addresses", deviceNames)
	}

	return nil
}

func resourceBigipGtmServerExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching Data center " + name)

	dc, err := client.GetDatacenter(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Data center  (%s) (%v)", name, err)
		return false, err
	}

	if dc == nil {
		d.SetId("")
	}

	return dc != nil, nil
}

func resourceBigipGtmDataCenterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	dc := &bigip.Datacenter{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Contact:     d.Get("contact").(string),
		AppService:  d.Get("appservice").(string),
		Enabled:     d.Get("enabled").(bool),
		ProberPool:  d.Get("prober_pool").(string),
	}
	err := client.ModifyDatacenter(name, dc)
	if err != nil {
		log.Printf("[ERROR] Unable to modify Data center (%s) (%v)", name, err)
		return err
	}

	return resourceBigipGtmDataCenterRead(d, meta)
}

func resourceBigipGtmDataCenterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Data center " + name)

	err := client.DeleteDatacenter(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Datacenter  (%s) (%v)", name, err)
		return err
	}
	d.SetId("")
	return nil
}

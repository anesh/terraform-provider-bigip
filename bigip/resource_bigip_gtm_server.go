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

			"virtualserverdiscovery": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"devices": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"address": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}
func resourceBigipGtmServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	datacenter := d.Get("datacenter").(string)
	virtualserverdiscovery := d.Get("virtualserverdiscovery").(bool)
	rs := d.Get("devices").(*schema.Set)
	var records []bigip.DeviceRecord
	if rs.Len() > 0 {
		for _, r := range rs.List() {
			record := r.(map[string]interface{})
			records = append(records, bigip.DeviceRecord{Name: record["name"].(string), Address: record["address"].(string)})
		}
	}
	log.Println("[INFO] Creating GTM Server " + name)
	config := &bigip.Server{
		Name:                     name,
		Datacenter:               datacenter,
		Virtual_server_discovery: virtualserverdiscovery,
		Devices:                  records,
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
	var records []map[string]interface{}
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
	for _, record := range sv.Devices {
		dRecord := map[string]interface{}{
			"name":    record.Name,
			"address": record.Address,
		}
		records = append(records, dRecord)
	}
	if err := d.Set("devices", records); err != nil {
		return fmt.Errorf("Error updating devicein state for GTM Server %s: %v ", name, err)
	}
	_ = d.Set("name", name)
	_ = d.Set("datacenter", sv.Datacenter)
	_ = d.Set("virtualserverdiscovery", sv.Virtual_server_discovery)

	return nil
}

func resourceBigipGtmServerExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching Data center " + name)

	sv, err := client.GetGtmserver(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve GTM server  (%s) (%v)", name, err)
		return false, err
	}

	if sv == nil {
		d.SetId("")
	}

	return sv != nil, nil
}

func resourceBigipGtmServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	rs := d.Get("devices").(*schema.Set)
	var records []bigip.DeviceRecord
	if rs.Len() > 0 {
		for _, r := range rs.List() {
			record := r.(map[string]interface{})
			records = append(records, bigip.DeviceRecord{Name: record["name"].(string), Address: record["address"].(string)})
		}
	}

	sv := &bigip.Server{
		Name:                     d.Get("name").(string),
		Datacenter:               d.Get("datacenter").(string),
		Virtual_server_discovery: d.Get("virtualserverdiscovery").(bool),
		Devices:                  records,
	}
	err := client.UpdateGtmserver(name, sv)
	if err != nil {
		log.Printf("[ERROR] Unable to modify GTM Server (%s) (%v)", name, err)
		return err
	}

	return resourceBigipGtmServerRead(d, meta)
}

func resourceBigipGtmServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting GTM server " + name)

	err := client.DeleteGtmserver(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete GTM server  (%s) (%v)", name, err)
		return err
	}
	d.SetId("")
	return nil
}

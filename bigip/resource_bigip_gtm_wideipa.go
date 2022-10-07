package bigip

import (
	"fmt"
	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceBigipGtmWideipa() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipGtmWideipaCreate,
		Read:   resourceBigipGtmWideipaRead,
		Update: resourceBigipGtmWideipaUpdate,
		Delete: resourceBigipGtmWideipaDelete,
		Exists: resourceBigipGtmWideipaExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of wideip",
				ForceNew:     true,
				ValidateFunc: validateF5NameWithDirectory,
			},
			"pools": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"order": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"ratio": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}
func resourceBigipGtmWideipaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	pools := d.Get("pools").(*schema.Set)
	var memrecords []bigip.MemberRecord
	if pools.Len() > 0 {
		for _, r := range pools.List() {
			record := r.(map[string]interface{})
			memrecords = append(memrecords, bigip.MemberRecord{Name: record["name"].(string), Order: record["order"].(int), Ratio: record["ratio"].(int)})
		}
	}
	log.Println("[INFO] Creating GTM Wideip " + name)
	config := &bigip.Wideip_a{
		Name:  name,
		Pools: memrecords,
	}
	err := client.CreateWideip_a(config)
	if err != nil {
		log.Printf("[ERROR] Unable to Create GTM widiep  (%s) (%v)", name, err)
		return err
	}
	d.SetId(name)
	err = resourceBigipGtmWideipaUpdate(d, meta)
	if err != nil {
		_ = client.DeleteGtmWideip_a(name)
		return err
	}
	return resourceBigipGtmWideipaRead(d, meta)
}

func resourceBigipGtmWideipaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	var records []map[string]interface{}
	name := d.Id()
	log.Println("[INFO] Fetching GTM Wideip for read " + name)
	gw, err := client.GetWideip_a(name)
	if err != nil {
		return err
	}
	if gw == nil {
		log.Printf("[WARN] GTM Wideip (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	for _, record := range gw.Pools {
		dRecord := map[string]interface{}{
			"name":    record.Name,
			"address": record.Order,
			"ratio":   record.Ratio,
		}
		records = append(records, dRecord)
	}

	if err := d.Set("pools", records); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Pools to state for GTM Wideip (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceBigipGtmWideipaExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching GTM Wideip " + name)

	gw, err := client.GetWideip_a(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve GTM Wideip  (%s) (%v)", name, err)
		return false, err
	}

	if gw == nil {
		d.SetId("")
	}

	return gw != nil, nil
}

func resourceBigipGtmWideipaUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	pools := d.Get("pools").(*schema.Set)
	var memrecords []bigip.MemberRecord
	if pools.Len() > 0 {
		for _, r := range pools.List() {
			record := r.(map[string]interface{})
			memrecords = append(memrecords, bigip.MemberRecord{Name: record["name"].(string), Order: record["order"].(int), Ratio: record["ratio"].(int)})
		}
	}

	r := &bigip.Wideip_a{
		Name:  d.Get("name").(string),
		Pools: memrecords,
	}
	err := client.ModifyWideip_a(name, r)
	if err != nil {
		log.Printf("[ERROR] Unable to modify GTM Wideip (%s) (%v)", name, err)
		return err
	}

	return resourceBigipGtmWideipaRead(d, meta)
}

func resourceBigipGtmWideipaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting GTM Wideip " + name)

	err := client.DeleteGtmWideip_a(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete GTM Wideip  (%s) (%v)", name, err)
		return err
	}
	d.SetId("")
	return nil
}

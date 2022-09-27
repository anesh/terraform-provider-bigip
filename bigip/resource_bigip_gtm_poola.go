package bigip

import (
	"fmt"
	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceBigipGtmPoola() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipGtmPoolaCreate,
		Read:   resourceBigipGtmPoolaRead,
		Update: resourceBigipGtmPoolaUpdate,
		Delete: resourceBigipGtmPoolaDelete,
		Exists: resourceBigipGtmPoolaExists,
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
			"members": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
		},
	}
}
func resourceBigipGtmPoolaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	members := setToStringSlice(d.Get("members").(*schema.Set))

	log.Println("[INFO] Creating GTM Pool " + name)
	config := &bigip.Pool_a{
		Name:    name,
		Members: members,
	}
	err := client.CreatePool_a(config)
	if err != nil {
		log.Printf("[ERROR] Unable to Create GTM Pool  (%s) (%v)", name, err)
		return err
	}
	d.SetId(name)
	err = resourceBigipGtmPoolaUpdate(d, meta)
	if err != nil {
		_ = client.DeleteGtmserver(name)
		return err
	}
	return resourceBigipGtmPoolaRead(d, meta)
}

func resourceBigipGtmPoolaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Fetching GTM Pool " + name)
	gp, err := client.GetPool_a(name)
	if err != nil {
		return err
	}
	if gp == nil {
		log.Printf("[WARN] GTM Pool (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	if err := d.Set("members", gp.Members); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Members to state for GTM Pool (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceBigipGtmPoolaExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching GTM Pool " + name)

	gp, err := client.GetPool_a(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve GTM Pool  (%s) (%v)", name, err)
		return false, err
	}

	if gp == nil {
		d.SetId("")
	}

	return gp != nil, nil
}

func resourceBigipGtmPoolaUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	r := &bigip.Pool_a{
		Name:    d.Get("name").(string),
		Members: setToStringSlice(d.Get("members").(*schema.Set)),
	}
	err := client.ModifyPool_a(name, r)
	if err != nil {
		log.Printf("[ERROR] Unable to modify GTM Pool (%s) (%v)", name, err)
		return err
	}

	return resourceBigipGtmPoolaRead(d, meta)
}

func resourceBigipGtmPoolaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting GTM Pool " + name)

	err := client.DeleteGtmPool_a(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete GTM Pool  (%s) (%v)", name, err)
		return err
	}
	d.SetId("")
	return nil
}

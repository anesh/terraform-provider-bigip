package bigip

import (
	"fmt"
	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceBigipGtmDataCenter() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipGtmDataCenterCreate,
		Read:   resourceBigipGtmDataCenterRead,
		Update: resourceBigipGtmDataCenterUpdate,
		Delete: resourceBigipGtmDataCenterDelete,
		Exists: resourceBigipGtmDataCenterExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the data center",
				ForceNew:     true,
				ValidateFunc: validateF5NameWithDirectory,
			},

			"appservice": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The application service that the object belongs to",
			},
			"contact": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"prober_fallback": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"prober_pool": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"prober_preference": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceBigipGtmDataCenterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	d.SetId(name)
	log.Println("[INFO] Creating Data Center " + name)
	err := client.CreateDatacenter(name)
	if err != nil {
		return fmt.Errorf("Error retrieving datacenter (%s): %s", name, err)
	}
	err = resourceBigipGtmDataCenterUpdate(d, meta)
	if err != nil {
		_ = client.DeleteDatacenter(name)
		return err
	}
	return resourceBigipGtmDataCenterRead(d, meta)
}

func resourceBigipGtmDataCenterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	_ = d.Set("name", name)
	log.Println("[INFO] Fetching Data center " + name)
	dc, err := client.GetDatacenter(name)
	if err != nil {
		return err
	}
	if dc == nil {
		log.Printf("[WARN] DataCenter (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	return nil
}

func resourceBigipGtmDataCenterExists(d *schema.ResourceData, meta interface{}) (bool, error) {
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

package dome9

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/dome9/dome9-sdk-go/services/compliance/continuous_compliance_policy"
)

func resourceContinuousCompliancePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceContinuousCompliancePolicyCreate,
		Read:   resourceContinuousCompliancePolicyRead,
		Update: resourceContinuousCompliancePolicyUpdate,
		Delete: resourceContinuousCompliancePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cloud_account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"external_account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cloud_account_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bundle_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"notification_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceContinuousCompliancePolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	resp, _, err := client.continuousCompliancePolicy.Get(d.Id())
	if err != nil {
		return nil
	}

	log.Printf("[INFO] Getting continuous compliance policy: %+v\n", resp)
	d.SetId(resp.ID)
	_ = d.Set("cloud_account_id", resp.CloudAccountID)
	_ = d.Set("external_account_id", resp.ExternalAccountID)
	_ = d.Set("cloud_account_type", resp.CloudAccountType)
	_ = d.Set("bundle_id", resp.BundleID)
	if err := d.Set("notification_ids", flattenNotificationIDs(resp)); err != nil {
		return err
	}

	return nil
}

func resourceContinuousCompliancePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	req := expandContinuousCompliancePolicyRequest(d)
	log.Printf("[INFO] Creating compliance policy request %+v\n", req)
	resp, _, err := client.continuousCompliancePolicy.Create(&req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Created compliance policy with ID: %v\n", resp.ID)
	d.SetId(resp.ID)

	return resourceContinuousCompliancePolicyRead(d, meta)
}

func resourceContinuousCompliancePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	log.Printf("[INFO] Updating continuous compliance policy ID: %v\n", d.Id())
	req := expandContinuousCompliancePolicyRequest(d)

	if _, _, err := client.continuousCompliancePolicy.Update(d.Id(), &req); err != nil {
		return err
	}

	return nil
}

func resourceContinuousCompliancePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	log.Printf("[INFO] Deleting continuous compliance policy ID: %v\n", d.Id())

	if _, err := client.continuousCompliancePolicy.Delete(d.Id()); err != nil {
		return err
	}

	return nil
}

func expandNotificationIDs(d *schema.ResourceData, key string) []string {
	notificationsIDsData := d.Get(key).([]interface{})
	notificationIDsList := make([]string, len(notificationsIDsData))
	for i, notificationID := range notificationsIDsData {
		notificationIDsList[i] = notificationID.(string)
	}

	return notificationIDsList
}

func expandContinuousCompliancePolicyRequest(d *schema.ResourceData) continuous_compliance_policy.ContinuousCompliancePolicyRequest {
	return continuous_compliance_policy.ContinuousCompliancePolicyRequest{
		CloudAccountID:    d.Get("cloud_account_id").(string),
		ExternalAccountID: d.Get("external_account_id").(string),
		BundleID:          d.Get("bundle_id").(int),
		NotificationIds:   expandNotificationIDs(d, "notification_ids"),
		CloudAccountType:  d.Get("cloud_account_type").(string),
	}
}

func flattenNotificationIDs(resp *continuous_compliance_policy.ContinuousCompliancePolicyResponse) []string {
	nIDs := make([]string, len(resp.NotificationIds))
	for i, nID := range resp.NotificationIds {
		nIDs[i] = nID
	}

	return nIDs
}
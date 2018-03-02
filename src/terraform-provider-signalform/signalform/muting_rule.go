package signalform

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

const (
	MUTING_RULE_API_URL = "https://api.signalfx.com/v2/alertmuting"
	MUTING_RULE_URL     = "https://api.signalfx.com/v2/alertmuting/<id>"
)

func mutingRuleResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"synced": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the resource in SignalForm and SignalFx are identical or not. Used internally for syncing.",
			},
			"last_updated": &schema.Schema{
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Latest timestamp the resource was updated",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Url of the muting rule",
			},
			// "resource_url": &schema.Schema{
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Default:     MUTING_RULE_URL,
			// 	Description: "Base muting rule url",
			// },
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the muting rule",
			},
			"start_time": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				// ConflictsWith: []string{"time_range"},
				Description: "Start time of muting period in milliseconds since epoch. If startTime is not specified, it will default its value to the current time in UTC milliseconds.",
			},
			"stop_time": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				// ConflictsWith: []string{"time_range"},
				Description: "Stop time of muting period in milliseconds since epoch. 0 indicates that the detectors matching this rule will be muted indefinitely. If stopTime is not specified, the value will default to 0.",
			},
			"filter": &schema.Schema{
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Set of rules used for muting",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"property": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter property",
						},
						"property_value": &schema.Schema{
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Filter property value",
						},
						"not": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Add not to the rule",
						},
					},
				},
			},
		},

		Create: mutingRuleCreate,
		Read:   mutingRuleRead,
		Update: mutingRuleUpdate,
		Delete: mutingRuleDelete,
	}
}

func getPayloadMutingRule(mr *schema.ResourceData) ([]byte, error) {

	tf_filters := mr.Get("filter").(*schema.Set).List()
	filters_list := make([]map[string]interface{}, len(tf_filters))

	for i, tf_filter := range tf_filters {
		tf_filter := tf_filter.(map[string]interface{})
		item := make(map[string]interface{})

		item["property"] = tf_filter["property"].(string)

		if filterValues, ok := tf_filter["property_value"]; ok {
			values := strings.Join(filterValues.([]string), ",")
			item["propertyValue"] = values
		}

		if val, ok := tf_filter["not"]; ok {
			item["not"] = val.(string)
		}

		filters_list[i] = item
	}

	payload := map[string]interface{}{
		"description": mr.Get("description").(string),
		"start_time":  mr.Get("start_time").(string),
		"stop_time":   mr.Get("stop_time").(string),
		"filters":     filters_list,
	}

	return json.Marshal(payload)
}

func mutingRuleCreate(mr *schema.ResourceData, meta interface{}) error {
	config := meta.(*signalformConfig)
	payload, err := getPayloadMutingRule(mr)
	if err != nil {
		return fmt.Errorf("Failed creating json payload: %s", err.Error())
	}

	return resourceCreate(MUTING_RULE_API_URL, config.AuthToken, payload, mr)
}

func mutingRuleRead(mr *schema.ResourceData, meta interface{}) error {
	config := meta.(*signalformConfig)
	url := fmt.Sprintf("%s/%s", MUTING_RULE_API_URL, mr.Id())

	return resourceRead(url, config.AuthToken, mr)
}

func mutingRuleUpdate(mr *schema.ResourceData, meta interface{}) error {
	config := meta.(*signalformConfig)
	payload, err := getPayloadDetector(mr)
	if err != nil {
		return fmt.Errorf("Failed creating json payload: %s", err.Error())
	}
	url := fmt.Sprintf("%s/%s", MUTING_RULE_API_URL, mr.Id())

	return resourceUpdate(url, config.AuthToken, payload, mr)
}

func mutingRuleDelete(mr *schema.ResourceData, meta interface{}) error {
	config := meta.(*signalformConfig)
	url := fmt.Sprintf("%s/%s", MUTING_RULE_API_URL, mr.Id())

	return resourceDelete(url, config.AuthToken, mr)
}

// func resourceFilterHash(v interface{}) int {
// 	var buf bytes.Buffer
// }

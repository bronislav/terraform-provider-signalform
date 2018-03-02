package signalform

// import (
// 	"github.com/hashicorp/terraform/helper/hashcode"
// 	"github.com/stretchr/testify/assert"
// 	"testing"
// )

// func TestFilterHash(t *testing.T) {
// 	values := map[string]interface{}{
// 		"property":       "property name",
// 		"property_value": "property value",
// 	}

// 	expected := hashcode.String("property name-property value-false")
// 	assert.Equal(t, expected, resourceFilterHash(values))

// 	values = map[string]interface{}{
// 		"property": "property name",
// 		"property_value": []interface{}{
// 			"value1",
// 			"value2",
// 		},
// 	}

// 	expected := hashcode.String("property name-property value-false")
// 	assert.Equal(t, expected, resourceFilterHash(values))
// }

package eventprocessing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sensu/sensu-go/types"
)

var (
	stdin *os.File
)

func GetPipedEvent() (*types.Event, error) {
	if stdin == nil {
		stdin = os.Stdin
	}

	eventJSON, err := ioutil.ReadAll(stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read stdin: %s", err.Error())
	}

	event := &types.Event{}
	err = json.Unmarshal(eventJSON, event)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal stdin data: %s", err.Error())
	}

	if err = event.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate event: %s", err.Error())
	}

	if !event.HasMetrics() {
		return nil, fmt.Errorf("event does not contain metrics")
	}

	return event, nil
}

type Tag struct {
	tag   string
	value string
}
type MetricValue struct {
	timestamp string
	name      string
	entity    string
	value     float64
	namespace string
	tags      []Tag
}

// {
// 	"name": "avg_cpu",
// 	"value": "56.0",
// 	"timestamp": "2019-03-30 12:30:00.45",
// 	"entity": "demo_test_agent",
// 	"namespace": "demo_jk185160",
// 	"tags": [
// 			{
// 					"tag": "company",
// 					"value": "JKTE001"
// 			},
// 			{
// 					"tag": "site",
// 					"value": "1001"
// 			}
// 	]
// }

func parsePointTimestamp(point *types.MetricPoint) (string, error) {
	stringTimestamp := strconv.FormatInt(point.Timestamp, 10)
	if len(stringTimestamp) > 10 {
		stringTimestamp = stringTimestamp[:10]
	}
	t, err := strconv.ParseInt(stringTimestamp, 10, 64)
	if err != nil {
		return "", err
	}
	return time.Unix(t, 0).Format(time.RFC3339), nil
}

func GetMetricFromPoint(point *types.MetricPoint, entityID string, namespaceID string) (MetricValue, error) {
	var metric MetricValue

	metric.entity = entityID
	metric.namespace = namespaceID
	// Find metric name
	nameField := strings.Split(point.Name, ".")
	metric.name = nameField[0]

	// Find metric timstamp
	unixTimestamp, err := parsePointTimestamp(point)
	if err != nil {
		return *new(MetricValue), fmt.Errorf("failed to validate event: %s", err.Error())
	}
	metric.timestamp = unixTimestamp
	metric.tags = make([]Tag, len(point.Tags)+1)
	i := 0
	for _, tag := range point.Tags {
		var thisTag Tag
		thisTag.tag = tag.Name
		thisTag.value = tag.Value
		metric.tags[i] = thisTag
		i++
	}
	var entityNameTag Tag
	entityNameTag.tag = "sensu_entity_name"
	entityNameTag.value = entityID
	metric.value = point.Value
	return metric, nil
}

package eventprocessing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

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

type MetricValue struct {
	Timestamp int64
	Name      string
	Entity    string
	Value     float64
	Tags      map[string]string
}

func parsePointTimestamp(point *types.MetricPoint) (int64, error) {
	stringTimestamp := strconv.FormatInt(point.Timestamp, 10)
	if len(stringTimestamp) > 10 {
		stringTimestamp = stringTimestamp[:10]
	}
	t, err := strconv.ParseInt(stringTimestamp, 10, 64)
	if err != nil {
		return 0, err
	}
	return t, nil
}

func GetMetricFromPoint(point *types.MetricPoint, entityID string) (MetricValue, error) {
	var metric MetricValue

	metric.Entity = entityID
	// Find metric name
	nameField := strings.Split(point.Name, ".")
	metric.Name = nameField[0]

	// Find metric timstamp
	unixTimestamp, err := parsePointTimestamp(point)
	if err != nil {
		return *new(MetricValue), fmt.Errorf("failed to validate event: %s", err.Error())
	}
	metric.Timestamp = unixTimestamp
	metric.Tags = make(map[string]string)
	metric.Tags["sensu_entity_name"] = entityID
	for _, tag := range point.Tags {
		metric.Tags[tag.Name] = tag.Value
	}
	metric.Value = point.Value
	return metric, nil
}

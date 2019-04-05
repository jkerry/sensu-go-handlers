// +build linux darwin
package main

import (
	"github.com/jkerry/sensu_gcp_pubsub_handler/lib/app/pubsub"
)

func main() {
	pubsub.SendMetric()
}

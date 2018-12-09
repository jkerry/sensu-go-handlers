// +build linux darwin
package main

import (
	"github.com/jkerry/sensu-go-handlers/lib/app/pubsub"
)

func main() {
	pubsub.SendMetric()
}

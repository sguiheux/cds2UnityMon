package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"

	"fmt"
	"github.com/ovh/cds/sdk"
	"github.com/ovh/cds/sdk/event"
)

var mapPb map[string]string

func consumeFromKafka(kafka, topic, group, username, password string) {
	mapPb = make(map[string]string)
	event.ConsumeKafka(kafka, topic, group, username, password, func(e sdk.Event) error {
		return process(e)
	}, log.Errorf)
}

// Process send message to all notifications backend
func process(event sdk.Event) error {
	log.Debugf("process> receive: type:%s all: %+v", event.EventType, event)

	if event.EventType == fmt.Sprintf("%T", sdk.EventPipelineBuild{}) {
		var e sdk.EventPipelineBuild
		if err := mapstructure.Decode(event.Payload, &e); err != nil {
			log.Errorf("Error during consumption EventPipelineBuild: %s", err)
		} else {
			return processEventPipelineBuild(&e)
		}
	}
	return nil
}

func processEventPipelineBuild(e *sdk.EventPipelineBuild) error {
	key := fmt.Sprintf("%s-%s-%s-%s-%s-%d", e.ProjectKey, e.ApplicationName, e.PipelineName, e.EnvironmentName, e.BranchName, e.BuildNumber)
	_, found := mapPb[key]
	if !found {
		switch e.Status.String() {
		case "Success":
		case "Fail":
		default:
			mapPb[key] = e.Status.String()
		}
		return nil
	}

	switch e.Status.String() {
	case "Success":
		delete(mapPb, key)
	case "Fail":
		delete(mapPb, key)
	default:
		mapPb[key] = e.Status.String()
	}

	return nil
}
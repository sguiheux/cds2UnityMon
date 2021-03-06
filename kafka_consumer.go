package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"
	"github.com/ovh/cds/sdk"
	"github.com/ovh/cds/sdk/event"
	"strings"
)

var mapPb map[string]string

func consumeFromKafka(kafka, topic, group, username, password string) {
	log.Info("Init kafka consumer")
	mapPb = make(map[string]string)
	if err := event.ConsumeKafka(kafka, topic, group, username, password, func(e sdk.Event) error {
		return process(e)
	}, log.Errorf); err != nil {
		log.Errorf("Unable to init kafka consumer: %s", err)
	}
}

// Process send message to all notifications backend
func process(event sdk.Event) error {
	log.Infof("process> receive: type:%s all: %+v", event.EventType, event)
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
	key := strings.Replace(fmt.Sprintf("%s-%s-%s-%s-%s-%d", e.ProjectKey, e.ApplicationName, e.PipelineName, e.EnvironmentName, e.BranchName, e.BuildNumber), "/", "-", -1)
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

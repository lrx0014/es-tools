package es

import (
	"log"
	"testing"
)

func TestGetPodListFromeES(t *testing.T) {
	_, err := GetPodList("game2048")
	if err != nil {
		log.Fatalf("Unable to get pods list: %v", err)
	}
}

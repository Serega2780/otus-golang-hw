package integration

import (
	"log"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

const delay = 3 * time.Second

func TestFeatures(t *testing.T) {
	log.Printf("wait %s for service availability...", delay)
	time.Sleep(delay)

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:    "pretty",
			Paths:     []string{"features"},
			TestingT:  t, // Testing instance that will run subtests.
			Randomize: 0,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

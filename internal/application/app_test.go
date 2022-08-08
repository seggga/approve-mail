package application

import (
	"strings"
	"testing"
)

var (
	configText = `
# config.yaml
# logger level, kafka pub and sub parameters, smtp ratelimit in seconds and email comopsers run simultaneously 

logger:
  level: debug

kafka:
  sub:
    server: "91.185.95.87:9094"
    topic: "team9-task-mail"
    group_id: "team9-consumer-group"

  pub:
    server: "91.185.95.87:9094"
    topic: "team9-mail-analytics"

mail:
  rate: 5  # rate limit in seconds on sending emails

compose:
  workers: 2 # number of workers for composing smtp-messages simultaneously
`

	cfgExpected = Config{
		Logger: Logger{
			Level: "debug",
		},
		Kafka: Kafka{
			Sub: KafkaSub{
				Server:  "91.185.95.87:9094",
				Topic:   "team9-task-mail",
				GroupID: "team9-consumer-group",
			},
			Pub: KafkaPub{
				Server: "91.185.95.87:9094",
				Topic:  "team9-mail-analytics",
			},
		},
		Mail: Mail{
			Rate: 5,
		},
		Compose: Compose{
			Workers: 2,
		},
	}
)

func TestReadConfig(t *testing.T) {

	cfg := readConfigFile(strings.NewReader(configText))

	if cfgExpected != *cfg {
		t.Errorf("error reading config: expected %v, got %v", cfgExpected, *cfg)
	}
}

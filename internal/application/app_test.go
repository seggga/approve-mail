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
    server: "127.0.0.1:9093"
    topic: "task-mail"
    group_id: "consumer-group"

  pub:
    server: "127.0.0.1:9093"
    topic: "mail-analytics"

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
				Server:  "127.0.0.1:9093",
				Topic:   "task-mail",
				GroupID: "consumer-group",
			},
			Pub: KafkaPub{
				Server: "127.0.0.1:9093",
				Topic:  "mail-analytics",
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

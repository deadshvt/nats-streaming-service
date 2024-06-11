package nats

import (
	"github.com/nats-io/stan.go"
)

func Init(clusterID, clientID string, url string) (stan.Conn, error) {
	nc, err := stan.Connect(clusterID, clientID, stan.NatsURL(url))
	if err != nil {
		return nil, err
	}

	return nc, nil
}

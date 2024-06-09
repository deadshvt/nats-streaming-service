package nats

import (
	"github.com/nats-io/stan.go"
)

func Init(clusterID, clientID string) (stan.Conn, error) {
	nc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		return nil, err
	}

	return nc, nil
}

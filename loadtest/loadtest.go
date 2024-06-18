package main

import (
	"fmt"
	"os"
	"time"

	"github.com/deadshvt/nats-streaming-service/config"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	config.Load(".env")

	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 3 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL: fmt.Sprintf("http://localhost:%s/order/%s",
			os.Getenv("SERVER_PORT"),
			os.Getenv("VEGETA_ORDER_ID")),
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Get order") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Rate: %f\n", metrics.Rate)
	fmt.Printf("Success: %f\n", metrics.Success*100)
	fmt.Printf("Latency: \n")
	fmt.Printf("  - Mean: %s\n", metrics.Latencies.Mean)
	fmt.Printf("  - P99: %s\n", metrics.Latencies.P99)

	if len(metrics.Errors) > 0 {
		fmt.Println("Errors:")
		for _, err := range metrics.Errors {
			fmt.Println(err)
		}
	}
}

package main

// otlpumper
// written by marTIN FISHer
// 04.11.2023
// the purpose of this script is pump a configured o11y narrative from a supplied
// yaml file (narrative.yaml) 


import (
	"context"
	"fmt"
	"io/ioutil"
	"sync"
	"time"
        "os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"gopkg.in/yaml.v2"
)

type MetricConfig struct {
	Narrative map[string]MetricTune
}

type MetricTune struct {
        Metric     string
	Attributes []map[string]string
	Type       string
	Sequence   []int
}

func initMTS(metricIn MetricTune, narrative string, ch chan<- MetricTune) {
        otlpEP := os.Getenv("OTLP_ENDPOINT")
        if otlpEP == "" {
        	fmt.Println("OTLP_ENDPOINT environment variable is not set, using 127.0.0.1:4317")
        	otlpEP="127.0.0.1:4317"
        }
	fmt.Printf("InitMTS: %s Metric: %s Attributes: %s Type: %s Sequence: %d metrics\n", narrative, metricIn.Metric, metricIn.Attributes, metricIn.Type, len(metricIn.Sequence))
	ctx := context.Background()
	exp, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(otlpEP),
	)
	if err != nil {
		panic(err)
	}
	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp)))
	defer func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	otel.SetMeterProvider(meterProvider)
	var meter = otel.Meter(metricIn.Metric)
	cardinality := []attribute.KeyValue{}
	for _, attr := range metricIn.Attributes {
		for key, val := range attr {
			cardinality = append(cardinality, attribute.String(key, val))
		}
	}
	index := 0
	for {
		value := metricIn.Sequence[index]
		_, obsErr := meter.Int64ObservableGauge(
			metricIn.Metric,
			metric.WithInt64Callback(func(_ context.Context, obsrv metric.Int64Observer) error {
				fmt.Printf("Emitting Metric: %s, Value: %d (iter:%d)\n", metricIn.Metric, value, index)
				obsrv.Observe(int64(value), metric.WithAttributes(cardinality...))
				return nil
			}),
		)
		if obsErr != nil {
			fmt.Println("failed to register instrument")
			panic(obsErr)
		}
		index = (index + 1) % len(metricIn.Sequence)
		time.Sleep(time.Minute)
	}
	ch <- metricIn
}

func main() {
	yamlFile, err := ioutil.ReadFile("./narrative.yaml")
	if err != nil {
		fmt.Printf("yaml read error: %v\n", err)
		return
	}
	var config MetricConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Printf("yaml unmarshal fail: %v\n", err)
		return
	}
	channels := make(map[string]chan MetricTune)
	var wg sync.WaitGroup
	for narrative, metric := range config.Narrative {
		ch := make(chan MetricTune)
		channels[narrative] = ch
		wg.Add(1)
		go func(narrative string, metric MetricTune, ch chan MetricTune) {
			defer wg.Done()
			initMTS(metric, narrative, ch)
		}(narrative, metric, ch)
	}
	wg.Wait()
}


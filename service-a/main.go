package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

func initProvider(serviceName, collectorURL string) (func(ctx context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(serviceName)))
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, collectorURL, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to collector: %w", err)
	}
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp))
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}

type CEPRequest struct {
	CEP string `json:"cep"`
}

func main() {
	http.HandleFunc("/", CEPHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Starting server on port %s...", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func CEPHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"message": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"message": "failed to read request body"}`, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var req CEPRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, `{"message": "invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}

	if !validateCEP(req.CEP) {
		http.Error(w, `{"message": "invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}

	respBody, status, err := forwardToServiceB(req.CEP)
	if err != nil {
		http.Error(w, `{"message": "failed to contact service B"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(respBody)
}

func validateCEP(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}

func forwardToServiceB(cep string) ([]byte, int, error) {
	url := "http://service-b:8082/" + cep
	resp, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return respBody, resp.StatusCode, nil
}

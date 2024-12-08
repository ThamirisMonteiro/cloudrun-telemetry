package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

func initTracer(serviceName, zipkinURL string) (func(ctx context.Context) error, error) {
	exporter, err := zipkin.New(zipkinURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Zipkin exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(serviceName))),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp.Shutdown, nil
}

func main() {
	zipkinURL := os.Getenv("ZIPKIN_URL")
	if zipkinURL == "" {
		zipkinURL = "http://zipkin:9411/api/v2/spans"
	}

	shutdown, err := initTracer("service-a", zipkinURL)
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer func() { _ = shutdown(context.Background()) }()

	http.HandleFunc("/", CEPHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func CEPHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.GetTracerProvider().Tracer("service-a")
	ctx, span := tracer.Start(r.Context(), "CEPHandler")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"message": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"message": "failed to read request body"}`, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var req struct {
		CEP string `json:"cep"`
	}
	err = json.Unmarshal(body, &req)
	if err != nil || !validateCEP(req.CEP) {
		http.Error(w, `{"message": "invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}

	respBody, status, err := forwardToServiceB(ctx, req.CEP)
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

func forwardToServiceB(ctx context.Context, cep string) ([]byte, int, error) {
	tracer := otel.GetTracerProvider().Tracer("service-a")
	ctx, span := tracer.Start(ctx, "forwardToServiceB")
	defer span.End()

	url := "http://service-b:8082/" + cep
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

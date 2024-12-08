package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
)

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	Estado     string `json:"estado"`
}

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

	shutdown, err := initTracer("service-b", zipkinURL)
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer func() { _ = shutdown(context.Background()) }()

	http.HandleFunc("/", cepHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func cepHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.GetTracerProvider().Tracer("service-b")
	ctx, span := tracer.Start(r.Context(), "cepHandler")
	defer span.End()

	cep := r.URL.Path[1:]

	if err := validateCEP(cep); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	location, err, status := getLocationFromCEP(ctx, cep)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	temp, err := getTemperature(ctx, location.Localidade)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(temp))
}

func validateCEP(cep string) error {
	re := regexp.MustCompile(`^\d{8}$`)
	if !re.MatchString(cep) {
		return errors.New("invalid zipcode")
	}
	return nil
}

func getLocationFromCEP(ctx context.Context, cep string) (*ViaCEPResponse, error, int) {
	tracer := otel.GetTracerProvider().Tracer("service-b")
	_, span := tracer.Start(ctx, "getLocation")
	defer span.End()

	url := "http://viacep.com.br/ws/" + cep + "/json/"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err, http.StatusInternalServerError
	}

	var response ViaCEPResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	if response.Localidade == "" {
		return nil, errors.New("invalid zipcode"), http.StatusNotFound
	}

	return &response, nil, http.StatusOK
}

func getTemperature(ctx context.Context, location string) (string, error) {
	tracer := otel.GetTracerProvider().Tracer("service-b")
	_, span := tracer.Start(ctx, "getTemperature")
	defer span.End()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		log.Fatal("A chave de API não foi definida em WEATHER_API_KEY")
	}

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", weatherAPIKey, location)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching weather data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("received non-OK status code %d from weather API", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading weather response body: %v", err)
	}

	var response WeatherAPIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling weather response: %v", err)
	}

	tempF := response.Current.TempC*1.8 + 32
	tempK := response.Current.TempC + 273.15

	tempC := math.Round(response.Current.TempC*10) / 10
	tempF = math.Round(tempF*10) / 10
	tempK = math.Round(tempK*10) / 10

	data := map[string]interface{}{
		"city":   location,
		"temp_C": tempC,
		"temp_F": tempF,
		"temp_K": tempK,
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshalling temperature data: %v", err)
	}

	return string(dataJSON), nil
}

# Projeto de Observabilidade e Open Telemetry

## Como Rodar

### 1. Configurar Variáveis de Ambiente

Crie um arquivo `.env` com a chave da WeatherAPI:

```env
WEATHER_API_KEY=SuaChaveDeAPI
ZIPKIN_URL=http://zipkin:9411/api/v2/spans
```

### 2. Subir os Contêineres

Execute o seguinte comando para subir os contêineres:

```bash
docker-compose up --build
```

### 3. Testar

Envie uma requisição para o `service-a` com um CEP:

```bash
curl http://localhost:8081/80035050
```

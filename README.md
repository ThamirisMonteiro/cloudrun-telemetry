# Passos para Implementação do Serviço A

## 1. Planejamento e Setup do Projeto

1. Inicialize o projeto Go com o comando `go mod init`.
2. Estruture o projeto com pastas recomendadas, como `cmd`, `internal`, e `pkg`.
3. Adicione as dependências necessárias:
	- Framework HTTP (ex.: `github.com/gin-gonic/gin` ou `net/http` nativo).
	- OpenTelemetry SDK e Zipkin (ex.: `go.opentelemetry.io/otel` e `go.opentelemetry.io/otel/exporters/zipkin`).
	- Biblioteca para validação de entradas (ex.: `github.com/go-playground/validator/v10`).
4. Crie um Dockerfile para incluir o binário do serviço e suas dependências.
5. Configure o arquivo `docker-compose.yml` para o serviço e outros componentes, como Zipkin e OTEL Collector.
6. Defina variáveis de ambiente, como a URL do Serviço B, URLs do OTEL e Zipkin, entre outras.

## 2. Implementação do Serviço A

1. Crie um servidor HTTP e implemente um endpoint POST `/cep`.
2. Valide o input recebido:
	- Certifique-se de que o `cep` possui exatamente 8 dígitos e é uma string.
	- Retorne `422 - invalid zipcode` em caso de falha.
3. Encaminhe o CEP válido para o Serviço B via requisição HTTP.
4. Trate as respostas do Serviço B:
	- Para `422`, encaminhe a mensagem de erro.
	- Para `404`, encaminhe a mensagem de erro.
	- Para outros erros inesperados, retorne `500 - erro interno do servidor`.

## 3. Integração com OpenTelemetry e Zipkin

1. Configure o OTEL Collector com um exporter Zipkin.
2. Inicialize o OpenTelemetry no Serviço A:
	- Capture spans para medir o tempo de execução do endpoint `/cep`.
	- Adicione spans para a chamada HTTP ao Serviço B.
3. Propague o contexto de tracing distribuído entre os serviços.
4. Garanta que os traces são visíveis no Zipkin.

## 4. Testes e Documentação

1. Implemente testes unitários e de integração:
	- Teste a validação de CEPs.
	- Verifique a comunicação com o Serviço B em diferentes cenários (`200`, `422`, `404`).
	- Teste a integração do tracing com OTEL e Zipkin.
2. Documente o projeto em um arquivo `README.md`, incluindo:
	- Passos para rodar o serviço localmente.
	- Configuração do OTEL Collector e Zipkin.
	- Execução dos testes.

## 5. Deploy e Testes Finais

1. Suba os serviços com `docker-compose`:
	- Inclua no arquivo `docker-compose.yml` o Serviço A, o Serviço B, o OTEL Collector e o Zipkin.
2. Realize testes de fluxo completo:
	- Teste com CEPs válidos e inválidos.
	- Verifique os spans e traces no Zipkin.
3. Ajuste logs, traces e métricas para ambientes de desenvolvimento e produção.

## Estrutura Resumida do `docker-compose.yml`

version: "3.8"  
services:  
service-a:  
build:  
context: .  
ports:  
- "8080:8080"  
environment:  
- SERVICE_B_URL=http://service-b:8080  
- OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317  
- ZIPKIN_ENDPOINT=http://zipkin:9411/api/v2/spans  
depends_on:  
- service-b  
- otel-collector  
- zipkin

service-b:  
image: service-b-image  
ports:  
- "8081:8080"

otel-collector:  
image: otel/opentelemetry-collector:latest  
ports:  
- "4317:4317"

zipkin:  
image: openzipkin/zipkin:latest  
ports:  
- "9411:9411"  

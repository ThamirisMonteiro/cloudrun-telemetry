# Checklist para Implementação de OTEL, Zipkin e Docker

## 1. **Adicionar OTEL (OpenTelemetry) ao Serviço A e Serviço B**
- [ ] Instalar a biblioteca `otel` no **Serviço A** e **Serviço B**.
	- Adicionar as dependências de OpenTelemetry (Go SDK) e Zipkin.
- [ ] Inicializar o tracer do OpenTelemetry em ambos os serviços.
- [ ] Criar spans para medir a execução de operações importantes:
	- **Serviço A**: Criação do span para a validação do CEP e a comunicação com o Serviço B.
	- **Serviço B**: Criação do span para a validação do CEP e a comunicação com a API de temperatura.
- [ ] Iniciar o collector do OpenTelemetry ou configurar o Zipkin como coletor.
- [ ] Configurar os spans para serem enviados ao Zipkin.

## 2. **Integrar o Zipkin para Tracing**
- [ ] Adicionar as dependências do Zipkin no projeto (e.g., `otel-zipkin`).
- [ ] Configurar o serviço para enviar os dados de tracing para o Zipkin.
- [ ] Testar o tracing visualizando no painel do Zipkin (geralmente acessível via `http://localhost:9411`).

## 3. **Criar os Dockerfiles para o Serviço A e Serviço B**
- [ ] Criar um `Dockerfile` para o **Serviço A**:
	- Definir a imagem base (e.g., `golang:1.20`).
	- Copiar o código para dentro do contêiner.
	- Rodar os comandos necessários para compilar o código (e.g., `go build`).
	- Expor a porta em que o serviço vai rodar (e.g., `8081`).
	- Configurar o comando de inicialização do serviço (e.g., `./service-a`).

- [ ] Criar um `Dockerfile` para o **Serviço B**:
	- Semelhante ao do Serviço A, com a configuração da imagem, build e comando de inicialização.

## 4. **Criar o Docker Compose**
- [ ] Criar um `docker-compose.yml` para orquestrar ambos os serviços e o Zipkin.
	- Definir os contêineres para o **Serviço A** e **Serviço B**.
	- Definir um contêiner para o **Zipkin**:
		- Usar uma imagem oficial do Zipkin (e.g., `openzipkin/zipkin`).
		- Mapear as portas necessárias para acessar o painel de Zipkin.
- [ ] Definir as variáveis de ambiente e volumes necessários.
- [ ] Configurar links entre os contêineres para que os serviços possam se comunicar.

## 5. **Testar a Configuração Localmente**
- [ ] Rodar o `docker-compose up` para iniciar todos os contêineres (Serviço A, Serviço B e Zipkin).
- [ ] Verificar se os serviços estão funcionando corretamente e se estão se comunicando.
- [ ] Verificar se os traces estão sendo enviados corretamente ao Zipkin.
- [ ] Acessar o painel do Zipkin para ver os traces das requisições entre os serviços.

## 6. **Documentação**
- [ ] Documentar o processo de execução dos serviços e do Docker Compose.
- [ ] Instruir sobre como rodar os serviços localmente usando o `docker-compose`:
	- Como executar `docker-compose up`.
	- Como visualizar o painel do Zipkin.

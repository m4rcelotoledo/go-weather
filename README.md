# GoWeather

Uma API em Go que recebe um CEP brasileiro, busca a cidade correspondente e retorna a temperatura atual em Celsius, Fahrenheit e Kelvin.

## Funcionalidades

- Validação de CEP (8 dígitos numéricos)
- Integração com ViaCEP para buscar cidade pelo CEP
- Integração com WeatherAPI para obter temperatura atual
- Conversão automática para Fahrenheit e Kelvin
- Tratamento de erros com códigos HTTP apropriados
- Testes unitários abrangentes

## Tecnologias Utilizadas

- Go 1.24
- Chi Router (github.com/go-chi/chi/v5)
- Testify (github.com/stretchr/testify)
- Docker & Docker Compose

## Estrutura do Projeto

```
go-weather/
├── internal/
│   ├── clients/
│   │   ├── viacep/
│   │   │   ├── client.go      # Cliente ViaCEP
│   │   │   └── client_test.go
│   │   └── weather/
│   │       ├── client.go      # Cliente WeatherAPI
│   │       └── client_test.go
│   ├── config/
│   │   ├── env.go             # Configuração de ambiente
│   │   └── env_test.go
│   ├── handlers/
│   │   ├── weather.go         # Handler HTTP
│   │   └── weather_test.go
│   └── middleware/
│       ├── metrics.go         # Middleware de métricas
│       └── metrics_test.go
├── main.go                    # Servidor HTTP
├── main_test.go               # Testes de integração
├── go.mod
├── go.sum
├── Dockerfile                 # Multi-stage build
├── docker-compose.yml
├── docs/                      # Documentação técnica
└── README.md
```

## Configuração

### 1. Obter API Key do WeatherAPI

1. Acesse https://www.weatherapi.com/
2. Crie uma conta gratuita
3. Copie sua API Key

### 2. Configurar Variáveis de Ambiente

Copie o arquivo de exemplo e configure sua API key:

```bash
cp .env.example .env
```

Edite o arquivo `.env` e substitua `your_weather_api_key_here` pela sua API key:

```
WEATHER_API_KEY=sua_api_key_aqui
```

## Execução Local

### Usando Docker Compose (Recomendado)

```bash
# Construir e iniciar o serviço
docker-compose up --build

# Executar em background
docker-compose up -d --build

# Parar o serviço
docker-compose down
```

### Usando Go Diretamente

```bash
# Instalar dependências
go mod tidy

# Definir variável de ambiente
export WEATHER_API_KEY=sua_api_key_aqui

# Executar
go run main.go
```

## Uso da API

### Endpoints

#### Health Check
```
GET /health
```

#### Buscar Clima por CEP
```
GET /cep/{cep}
```

### Exemplo de Requisição

```bash
# Health check
curl http://localhost:8080/health

# Buscar clima por CEP
curl http://localhost:8080/cep/01001000
```

### Respostas

#### Sucesso (HTTP 200)
```json
{
    "temp_C": 28.5,
    "temp_F": 83.3,
    "temp_K": 301.5
}
```

#### CEP Inválido (HTTP 422)
```json
"invalid zipcode"
```

#### CEP Não Encontrado (HTTP 404)
```json
"can not find zipcode"
```

## Acesso à Aplicação em Produção

A aplicação está disponível no Google Cloud Run e pode ser acessada através dos seguintes endpoints:

### Health Check
```
GET https://go-weather-i6zfts7xxa-rj.a.run.app/health
```

**Resposta esperada:**
```json
{
    "status": "healthy",
    "time": "2025-08-06T22:41:24Z"
}
```

### Buscar Clima por CEP
```
GET https://go-weather-i6zfts7xxa-rj.a.run.app/cep/{cep}
```

**Exemplos de uso:**

```bash
# Health check
curl https://go-weather-i6zfts7xxa-rj.a.run.app/health

# Buscar clima para São Paulo (CEP: 01001000)
curl https://go-weather-i6zfts7xxa-rj.a.run.app/cep/01001000

# Buscar clima para Rio de Janeiro (CEP: 20040020)
curl https://go-weather-i6zfts7xxa-rj.a.run.app/cep/20040020

# Buscar clima para Belo Horizonte (CEP: 30112000)
curl https://go-weather-i6zfts7xxa-rj.a.run.app/cep/30112000
```

**Resposta de sucesso:**
```json
{
    "temp_C": 25.0,
    "temp_F": 77.0,
    "temp_K": 298.15
}
```

### Status da Aplicação

- ✅ **Deploy**: Concluído com sucesso
- ✅ **Health Check**: Funcionando
- ✅ **API de Clima**: Integrada e operacional
- ✅ **ViaCEP**: Integrado e operacional
- ✅ **HTTPS**: Configurado automaticamente pelo Cloud Run
- ✅ **Escalabilidade**: Configurada para 0-10 instâncias

### Informações Técnicas

- **Plataforma**: Google Cloud Run
- **Região**: southamerica-east1 (São Paulo)
- **URL Base**: https://go-weather-i6zfts7xxa-rj.a.run.app
- **Porta**: 8080 (configurada automaticamente)
- **Timeout**: 180 segundos
- **Memória**: 128Mi
- **CPU**: 1 vCPU

## Documentação e Scripts

### 📚 Documentação Técnica

A documentação detalhada do projeto está disponível na pasta `docs/`:

- **[DEPLOY_GUIDE.md](docs/DEPLOY_GUIDE.md)** - Guia completo para deploy no Google Cloud Run

### 🛠️ Scripts de Automação

Scripts úteis para desenvolvimento e deploy na pasta `scripts/`:

- **[deploy.sh](scripts/deploy.sh)** - Script automatizado para deploy no Google Cloud Run

**Uso do script de deploy:**
```bash
# Dar permissão de execução
chmod +x scripts/deploy.sh

# Executar deploy
./scripts/deploy.sh SEU_PROJETO_ID go-weather
```

### 📖 Estrutura de Documentação

```
docs/
└──  DEPLOY_GUIDE.md            # Guia de deploy no Cloud Run

scripts/
└── deploy.sh                  # Script de deploy automatizado
```

## Testes

### Executar Todos os Testes

```bash
go test ./...
```

### Executar Testes com Cobertura

```bash
go test -cover ./...
```

### Executar Testes Verbosos

```bash
go test -v ./...
```

### Cobertura de Testes

A aplicação possui cobertura abrangente de testes:

- **Testes Unitários**: Handlers, clients, config e middleware
- **Testes de Integração**: Roteamento, middleware e graceful shutdown
- **Testes de Cenários de Erro**: CEPs inválidos, APIs indisponíveis, timeouts
- **Testes de Configuração**: Variáveis de ambiente, valores padrão
- **Testes de Middleware**: Logging, timing, recuperação de pânico

## Deploy no Google Cloud Run

### 1. Preparar o Ambiente

```bash
# Instalar Google Cloud CLI
# https://cloud.google.com/sdk/docs/install

# Autenticar
gcloud auth login

# Configurar projeto
gcloud config set project SEU_PROJECT_ID

# Habilitar APIs necessárias
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable artifactregistry.googleapis.com
```

### 2. Criar Repositório no Artifact Registry

```bash
# Criar repositório
gcloud artifacts repositories create go-weather \
    --repository-format=docker \
    --location=us-central1 \
    --description="GoWeather API repository"

# Configurar Docker para usar o repositório
gcloud auth configure-docker us-central1-docker.pkg.dev
```

### 3. Build e Push da Imagem

```bash
# Definir variáveis
export PROJECT_ID=$(gcloud config get-value project)
export REGION=us-central1
export REPOSITORY=go-weather
export IMAGE=go-weather

# Build da imagem
docker build -t $REGION-docker.pkg.dev/$PROJECT_ID/$REPOSITORY/$IMAGE:latest .

# Push da imagem
docker push $REGION-docker.pkg.dev/$PROJECT_ID/$REPOSITORY/$IMAGE:latest
```

### 4. Deploy no Cloud Run

```bash
# Deploy com variável de ambiente
gcloud run deploy go-weather \
    --image=$REGION-docker.pkg.dev/$PROJECT_ID/$REPOSITORY/$IMAGE:latest \
    --platform=managed \
    --region=$REGION \
    --allow-unauthenticated \
    --port=8080 \
    --set-env-vars=WEATHER_API_KEY=SUA_API_KEY_AQUI

# Ou usando arquivo de ambiente (recomendado para produção)
gcloud run deploy go-weather \
    --image=$REGION-docker.pkg.dev/$PROJECT_ID/$REPOSITORY/$IMAGE:latest \
    --platform=managed \
    --region=$REGION \
    --allow-unauthenticated \
    --port=8080 \
    --env-vars-file=.env
```

### 5. Obter URL do Serviço

```bash
gcloud run services describe go-weather --region=$REGION --format="value(status.url)"
```

## Comandos Úteis

### Desenvolvimento

```bash
# Formatar código
go fmt ./...

# Executar linter
go vet ./...

# Instalar dependências
go mod tidy

# Verificar dependências
go mod verify
```

### Docker

```bash
# Build manual da imagem
docker build -t go-weather .

# Executar container localmente
docker run -p 8080:8080 -e WEATHER_API_KEY=sua_api_key go-weather

# Ver logs do container
docker-compose logs -f go-weather
```

### Testes de Integração

```bash
# Testar endpoint com CEP válido
curl http://localhost:8080/01001000

# Testar endpoint com CEP inválido
curl http://localhost:8080/123

# Testar endpoint com CEP não encontrado
curl http://localhost:8080/00000000
```

## Arquitetura

A aplicação segue uma arquitetura limpa com separação de responsabilidades:

- **main.go**: Ponto de entrada e configuração do servidor
- **internal/handlers**: Lógica de apresentação (HTTP handlers)
- **internal/clients**: Clientes para APIs externas (ViaCEP e WeatherAPI)

## Contribuição

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## Licença

Este projeto é parte de um desafio acadêmico do curso Go Expert da Full Cycle.

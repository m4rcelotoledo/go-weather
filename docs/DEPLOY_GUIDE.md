# Guia de Deploy - Google Cloud Run

Este guia te ajudará a fazer o deploy da aplicação Go Weather no Google Cloud Run.

## Pré-requisitos

1. **Conta Google Cloud**: Tenha uma conta Google Cloud ativa
2. **Google Cloud SDK**: Instale o Google Cloud SDK
3. **Projeto Google Cloud**: Crie um projeto no Google Cloud Console
4. **API Key do Weather**: Tenha sua chave da API de clima

## Passo a Passo

### 1. Instalar Google Cloud SDK

```bash
# No Manjaro/Arch Linux
sudo pacman -S google-cloud-cli

# Ou baixe de: https://cloud.google.com/sdk/docs/install
```

### 2. Configurar o Google Cloud

```bash
# Fazer login
gcloud auth login

# Listar projetos disponíveis
gcloud projects list

# Configurar projeto (substitua SEU_PROJETO_ID)
gcloud config set project SEU_PROJETO_ID
```

### 3. Habilitar APIs Necessárias

```bash
# Habilitar Cloud Run API
gcloud services enable run.googleapis.com

# Habilitar Cloud Build API
gcloud services enable cloudbuild.googleapis.com
```

### 4. Configurar Secret Manager (Recomendado)

Para armazenar sua API key de forma segura:

```bash
# Criar secret
echo -n "SUA_API_KEY_AQUI" | gcloud secrets create weather-api-key --data-file=-

# Ou atualizar secret existente
echo -n "SUA_API_KEY_AQUI" | gcloud secrets versions add weather-api-key --data-file=-
```

### 5. Fazer o Deploy

#### Opção A: Usando o script automatizado

```bash
# Dar permissão de execução (se necessário)
chmod +x deploy.sh

# Executar deploy
./deploy.sh SEU_PROJETO_ID go-weather
```

#### Opção B: Comando manual

```bash
gcloud run deploy go-weather \
    --source . \
    --region us-central1 \
    --platform managed \
    --allow-unauthenticated \
    --set-env-vars="WEATHER_API_KEY=SUA_API_KEY_AQUI" \
    --memory 512Mi \
    --cpu 1 \
    --max-instances 10 \
    --timeout 300
```

### 6. Verificar o Deploy

```bash
# Ver URL do serviço
gcloud run services describe go-weather --region=us-central1 --format="value(status.url)"

# Ver logs
gcloud run services logs read go-weather --region=us-central1

# Testar a aplicação
curl https://SEU_SERVICO_URL/health
```

## Configurações do Deploy

- **Região**: us-central1 (você pode alterar)
- **Memória**: 512Mi
- **CPU**: 1 vCPU
- **Máximo de instâncias**: 10
- **Timeout**: 300 segundos
- **Acesso**: Público (--allow-unauthenticated)

## Variáveis de Ambiente

A aplicação espera as seguintes variáveis:

- `WEATHER_API_KEY`: Chave da API de clima
- `PORT`: Porta (padrão: 8080)

## Troubleshooting

### Erro: "Permission denied"
```bash
# Verificar permissões
gcloud auth list
gcloud config get-value project
```

### Erro: "API not enabled"
```bash
# Habilitar APIs necessárias
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
```

### Erro: "Quota exceeded"
- Verifique os limites do seu projeto no Google Cloud Console
- Considere usar uma região diferente

### Logs e Debugging
```bash
# Ver logs em tempo real
gcloud run services logs tail go-weather --region=us-central1

# Ver detalhes do serviço
gcloud run services describe go-weather --region=us-central1
```

## Custos

O Google Cloud Run cobra apenas pelo tempo de execução:
- **CPU**: ~$0.00002400 por vCPU-segundo
- **Memória**: ~$0.00000250 por GB-segundo
- **Requests**: $0.40 por milhão de requests

Para uma aplicação pequena, os custos são muito baixos (geralmente < $5/mês).

## Próximos Passos

1. **Configurar domínio personalizado**
2. **Configurar HTTPS**
3. **Implementar CI/CD**
4. **Configurar monitoramento**
5. **Configurar alertas**

## Comandos Úteis

```bash
# Atualizar serviço
gcloud run services update go-weather --region=us-central1

# Deletar serviço
gcloud run services delete go-weather --region=us-central1

# Listar serviços
gcloud run services list --region=us-central1

# Ver métricas
gcloud run services describe go-weather --region=us-central1
```

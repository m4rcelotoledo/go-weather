#!/bin/bash

# Script de deploy para Google Cloud Run
# Uso: ./deploy.sh [PROJECT_ID] [SERVICE_NAME]

set -e

# Configurações padrão
PROJECT_ID=${1:-"seu-projeto-id"}
SERVICE_NAME=${2:-"go-weather"}
REGION="us-central1"

echo "🚀 Iniciando deploy no Google Cloud Run..."
echo "Projeto: $PROJECT_ID"
echo "Serviço: $SERVICE_NAME"
echo "Região: $REGION"

# Verificar se o gcloud está instalado
if ! command -v gcloud &> /dev/null; then
    echo "❌ Google Cloud SDK não está instalado. Instale em: https://cloud.google.com/sdk/docs/install"
    exit 1
fi

# Verificar se está autenticado
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
    echo "🔐 Fazendo login no Google Cloud..."
    gcloud auth login
fi

# Configurar projeto
echo "📋 Configurando projeto..."
gcloud config set project $PROJECT_ID

# Habilitar APIs necessárias
echo "🔧 Habilitando APIs necessárias..."
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com

# Fazer deploy
echo "🏗️ Fazendo deploy..."
gcloud run deploy $SERVICE_NAME \
    --source . \
    --region $REGION \
    --platform managed \
    --allow-unauthenticated \
    --set-env-vars="WEATHER_API_KEY=$(gcloud secrets versions access latest --secret="weather-api-key" --project=$PROJECT_ID 2>/dev/null || echo '')" \
    --memory 512Mi \
    --cpu 1 \
    --max-instances 10 \
    --timeout 300

echo "✅ Deploy concluído!"
echo "🌐 URL do serviço:"
gcloud run services describe $SERVICE_NAME --region=$REGION --format="value(status.url)"

echo ""
echo "📝 Para verificar os logs:"
echo "gcloud run services logs read $SERVICE_NAME --region=$REGION"

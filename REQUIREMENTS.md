# Requisitos do Projeto - Sistema de Temperatura por CEP com OTEL + Zipkin

## 📋 Objetivo

Desenvolver um sistema distribuído em Go que:
- Recebe um CEP
- Identifica a cidade
- Retorna o clima atual (temperatura em Celsius, Fahrenheit e Kelvin) com o nome da cidade
- Implementa **OTEL (Open Telemetry)** e **Zipkin** para tracing distribuído

## 🏗️ Arquitetura

O sistema será composto por **2 serviços**:

```
┌─────────────┐         HTTP          ┌─────────────┐
│  Serviço A  │ ──────────────────>   │  Serviço B  │
│   (Input)   │      POST /cep        │ (Orquestração)│
└─────────────┘                       └─────────────┘
      │                                      │
      │                                      ├─> ViaCEP API
      │                                      └─> WeatherAPI
      │
      └────────> OTEL Collector ←───────────┘
                        │
                        ▼
                    Zipkin UI
```

---

## 📦 Serviço A - Input e Validação

### Responsabilidades
- Receber requisições POST com CEP
- Validar o formato do CEP (8 dígitos numéricos)
- Encaminhar para o Serviço B via HTTP

### ✅ Requisitos Funcionais

#### Endpoint
- **Método:** `POST`
- **Path:** `/cep` (ou similar)
- **Content-Type:** `application/json`

#### Request Body Schema
```json
{
  "cep": "29902555"
}
```

#### Validações
1. **CEP deve ser uma STRING** ✓
2. **CEP deve conter exatamente 8 dígitos** ✓
3. **CEP deve conter apenas números** ✓

#### Respostas

| Cenário | HTTP Code | Response Body |
|---------|-----------|---------------|
| ✅ CEP válido | 200 | Resposta do Serviço B (delegada) |
| ❌ CEP inválido | 422 | `{ "message": "invalid zipcode" }` |

#### Fluxo
1. Recebe POST com CEP
2. Valida formato (8 dígitos, string numérica)
3. Se válido → Encaminha para Serviço B via HTTP
4. Se inválido → Retorna erro 422

---

## 🌐 Serviço B - Orquestração e Busca

### Responsabilidades
- Receber CEP do Serviço A
- Buscar localização via ViaCEP
- Buscar temperatura via WeatherAPI
- Converter temperaturas (C → F, K)
- Retornar dados formatados

### ✅ Requisitos Funcionais

#### Endpoint
- **Método:** `GET` ou `POST` (recebe do Serviço A)
- **Path:** `/weather` (ou similar)

#### Validações
1. **CEP deve ter 8 dígitos** ✓
2. **CEP deve existir no ViaCEP** ✓

#### Respostas

| Cenário | HTTP Code | Response Body |
|---------|-----------|---------------|
| ✅ Sucesso | 200 | `{ "city": "São Paulo", "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5 }` |
| ❌ Formato inválido | 422 | `{ "message": "invalid zipcode" }` |
| ❌ CEP não encontrado | 404 | `{ "message": "can not find zipcode" }` |

#### Conversões de Temperatura
- **Fahrenheit:** `F = C × 1.8 + 32`
- **Kelvin:** `K = C + 273`

#### APIs Externas
- **ViaCEP:** `https://viacep.com.br/ws/{cep}/json/`
- **WeatherAPI:** `https://api.weatherapi.com/v1/current.json?key={key}&q={city}`

---

## 🔍 Observabilidade - OTEL + Zipkin

### ✅ Requisitos de Tracing

#### 1. Tracing Distribuído
- [ ] Implementar propagação de contexto entre Serviço A → Serviço B
- [ ] Usar OTEL SDK para Go
- [ ] Configurar trace exporter para Zipkin

#### 2. Spans Obrigatórios

**Serviço A:**
- [ ] Span principal: `POST /cep`
- [ ] Span: `validate_cep`
- [ ] Span: `call_service_b` (medir tempo de chamada HTTP)

**Serviço B:**
- [ ] Span principal: `GET /weather` (ou equivalente)
- [ ] Span: `fetch_cep_location` (medir tempo ViaCEP)
- [ ] Span: `fetch_weather` (medir tempo WeatherAPI)
- [ ] Span: `convert_temperatures`

#### 3. Atributos dos Spans
- `cep`: CEP consultado
- `city`: Cidade encontrada
- `http.status_code`: Código HTTP de resposta
- `error`: Booleano indicando erro
- `error.message`: Mensagem de erro (se houver)

#### 4. Infraestrutura
- [ ] OTEL Collector configurado
- [ ] Zipkin rodando (porta padrão: 9411)
- [ ] Exportar traces via OTLP ou HTTP

---

## 🐳 Docker & Docker Compose

### ✅ Requisitos de Containerização

#### Serviços no Docker Compose
1. **service-a** - Serviço A (Input)
2. **service-b** - Serviço B (Orquestração)
3. **otel-collector** - OpenTelemetry Collector
4. **zipkin** - Zipkin UI

#### Configuração Esperada
```yaml
version: '3.8'

services:
  service-a:
    build: ./service-a
    ports:
      - "8080:8080"
    environment:
      - SERVICE_B_URL=http://service-b:8081
      - OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317
    depends_on:
      - service-b
      - otel-collector

  service-b:
    build: ./service-b
    ports:
      - "8081:8081"
    environment:
      - WEATHERAPI_KEY=${WEATHERAPI_KEY}
      - OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317
    depends_on:
      - otel-collector

  otel-collector:
    image: otel/opentelemetry-collector:latest
    command: ["--config=/etc/otel-collector-config.yaml"]
    volumes:
      - ./otel-collector-config.yaml:/etc/otel-collector-config.yaml
    ports:
      - "4317:4317"   # OTLP gRPC
      - "4318:4318"   # OTLP HTTP

  zipkin:
    image: openzipkin/zipkin:latest
    ports:
      - "9411:9411"
```

---

## 📚 Recursos e Documentação

### APIs Externas
- [ViaCEP](https://viacep.com.br/) - Consulta de CEP brasileiro
- [WeatherAPI](https://www.weatherapi.com/) - Dados meteorológicos (requer API key gratuita)

### OpenTelemetry
- [Getting Started - Go](https://opentelemetry.io/docs/languages/go/getting-started/)
- [Creating Spans](https://opentelemetry.io/docs/languages/go/instrumentation/#creating-spans)
- [OTEL Collector Quick Start](https://opentelemetry.io/docs/collector/quick-start/)

### Zipkin
- [Zipkin Quickstart](https://zipkin.io/pages/quickstart.html)

---

## 📋 Checklist de Implementação

### Fase 1: Serviço B (Orquestração)
- [ ] Criar projeto Go para Serviço B
- [ ] Implementar endpoint de recebimento de CEP
- [ ] Integrar com ViaCEP API
- [ ] Integrar com WeatherAPI
- [ ] Implementar conversões de temperatura
- [ ] Implementar tratamento de erros (422, 404)
- [ ] Criar testes unitários
- [ ] Criar Dockerfile

### Fase 2: Serviço A (Input)
- [ ] Criar projeto Go para Serviço A
- [ ] Implementar endpoint POST /cep
- [ ] Implementar validação de CEP (8 dígitos, string)
- [ ] Implementar chamada HTTP para Serviço B
- [ ] Implementar tratamento de erro 422
- [ ] Criar testes unitários
- [ ] Criar Dockerfile

### Fase 3: Observabilidade
- [ ] Instalar dependências OTEL Go SDK
- [ ] Configurar OTEL no Serviço A
- [ ] Configurar OTEL no Serviço B
- [ ] Implementar propagação de contexto (trace headers)
- [ ] Criar spans para validação de CEP
- [ ] Criar spans para chamada ao Serviço B
- [ ] Criar spans para ViaCEP
- [ ] Criar spans para WeatherAPI
- [ ] Adicionar atributos aos spans (cep, city, status)
- [ ] Configurar OTEL Collector
- [ ] Configurar export para Zipkin

### Fase 4: Containerização
- [ ] Criar docker-compose.yaml
- [ ] Configurar service-a no compose
- [ ] Configurar service-b no compose
- [ ] Configurar otel-collector no compose
- [ ] Configurar zipkin no compose
- [ ] Criar otel-collector-config.yaml
- [ ] Testar comunicação entre serviços
- [ ] Validar traces no Zipkin UI

### Fase 5: Documentação
- [ ] README.md com instruções de execução
- [ ] Documentar variáveis de ambiente necessárias
- [ ] Documentar endpoints da API
- [ ] Documentar como acessar Zipkin UI
- [ ] Documentar exemplos de requisições
- [ ] Adicionar screenshots do Zipkin (opcional)

### Fase 6: Testes e Validação
- [ ] Testar fluxo completo com CEP válido
- [ ] Testar CEP inválido (formato)
- [ ] Testar CEP não encontrado
- [ ] Validar traces no Zipkin
- [ ] Validar tempos de resposta nos spans
- [ ] Validar propagação de trace ID
- [ ] Testar em ambiente Docker

---

## 🎯 Critérios de Aceitação

### Funcional
✅ Serviço A valida CEP e encaminha para Serviço B  
✅ Serviço B busca localização e temperatura  
✅ Conversões de temperatura corretas  
✅ Códigos HTTP corretos (200, 404, 422)  
✅ Formato de resposta JSON conforme especificado  

### Observabilidade
✅ Traces distribuídos visíveis no Zipkin  
✅ Spans medem tempo de APIs externas  
✅ Propagação de contexto funcionando  
✅ Atributos dos spans incluem CEP e cidade  

### Infraestrutura
✅ Docker Compose sobe todos os serviços  
✅ OTEL Collector recebe e exporta traces  
✅ Zipkin UI acessível e funcional  
✅ Documentação clara de como executar  

---

## 🚀 Como Executar (Exemplo)

```bash
# 1. Configurar API Key da WeatherAPI
export WEATHERAPI_KEY=sua_chave_aqui

# 2. Subir infraestrutura
docker-compose up -d

# 3. Testar Serviço A
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "01001000"}'

# 4. Acessar Zipkin UI
# http://localhost:9411

# 5. Derrubar infraestrutura
docker-compose down
```

---

## 📦 Entregáveis

1. **Código-fonte completo**
   - Serviço A (Go)
   - Serviço B (Go)
   - Configurações OTEL
   - Docker Compose

2. **Documentação**
   - README.md com instruções de execução
   - Variáveis de ambiente necessárias
   - Exemplos de uso da API
   - Como visualizar traces no Zipkin

3. **Containerização**
   - Dockerfiles para ambos os serviços
   - docker-compose.yaml funcional
   - otel-collector-config.yaml

4. **Testes** (Opcional, mas recomendado)
   - Testes unitários
   - Testes de integração
   - Exemplos de requisições (.http files)

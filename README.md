# API de Cotação de Frete (Frete Rápido)

API REST em **Golang** para cotação de frete via Frete Rápido e consulta de métricas das cotações armazenadas. Aplicação containerizada com Docker. Desenvolvida com **TDD** (Test-Driven Development): testes escritos para definir o comportamento antes da implementação, seguindo o ciclo red-green-refactor.

## Requisitos

- **Docker** e **Docker Compose** (recomendado), ou
- **Go 1.21+** e **PostgreSQL 14+** para execução local

## Como executar com Docker (recomendado)

1. Na raiz do projeto, suba os containers:

```bash
docker-compose up --build
```

2. A API ficará disponível em **http://localhost:8080**.

3. O PostgreSQL fica acessível na porta **5433** no host (interno ao Docker usa 5432). A API conecta ao banco pela rede interna; não é necessário expor a porta para uso normal.

**Executar em segundo plano:** use `docker-compose up -d` (containers em background).

### Parar os serviços

```bash
docker-compose down
```

Para remover também o volume do banco:

```bash
docker-compose down -v
```

## Como executar localmente (sem Docker)

1. Tenha um PostgreSQL rodando e crie um banco (ex.: `quote_api`).

2. Copie as variáveis de ambiente:

```bash
cp .env.example .env
```

Ajuste `DB_*` no `.env` se necessário.

3. Baixe as dependências e rode a API:

```bash
go mod tidy
go run ./cmd/api
```

A API usa a porta definida em `SERVER_PORT` (padrão **8080**).

## Variáveis de ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `SERVER_PORT` | Porta HTTP da API | `8080` |
| `DB_HOST` | Host do PostgreSQL | `localhost` |
| `DB_PORT` | Porta do PostgreSQL | `5432` |
| `DB_USER` | Usuário do banco | `postgres` |
| `DB_PASSWORD` | Senha do banco | `postgres` |
| `DB_NAME` | Nome do banco | `quote_api` |
| `DB_SSLMODE` | SSL do PostgreSQL | `disable` |
| `FRETE_RAPIDO_BASE_URL` | URL base da API Frete Rápido | `https://sp.freterapido.com` |
| `FRETE_RAPIDO_TOKEN` | Token de autenticação | (valor do desafio) |
| `FRETE_RAPIDO_PLATFORM_CODE` | Código da plataforma | (valor do desafio) |
| `FRETE_RAPIDO_SHIPPER_CNPJ` | CNPJ remetente (apenas números) | `25438296000158` |
| `FRETE_RAPIDO_DISPATCHER_CEP` | CEP do expedidor (apenas números) | `29161376` |

## Endpoints

### 1. POST /quote

Cria uma cotação fictícia com a API Frete Rápido e persiste o resultado no banco.

**Corpo da requisição (JSON):**

```json
{
  "recipient": {
    "address": {
      "zipcode": "01311000"
    }
  },
  "volumes": [
    {
      "category": 7,
      "amount": 1,
      "unitary_weight": 5,
      "price": 349,
      "sku": "abc-teste-123",
      "height": 0.2,
      "width": 0.2,
      "length": 0.2
    }
  ]
}
```

**Regras de validação:**

- `recipient.address.zipcode`: obrigatório, exatamente 8 caracteres numéricos.
- `volumes`: obrigatório, pelo menos 1 item.
- Cada volume: `category` (≥ 1), `amount` (≥ 1), `unitary_weight` (> 0), `price` (≥ 0), `height`, `width`, `length` (> 0). `sku` opcional.

**Resposta de sucesso (200):**

```json
{
  "carrier": [
    {
      "name": "EXPRESSO FR",
      "service": "Rodoviário",
      "deadline": "3",
      "price": 17
    },
    {
      "name": "Correios",
      "service": "SEDEX",
      "deadline": "1",
      "price": 20.99
    }
  ]
}
```

**Exemplos de erro:**

- **400** – Dados inválidos (ex.: zipcode com menos de 8 caracteres, volumes vazios).
- **502** – Falha ao chamar a API Frete Rápido.
- **500** – Erro ao salvar cotação no banco.

---

### 2. GET /metrics?last_quotes={?}

Retorna métricas das cotações armazenadas. O parâmetro **last_quotes** é opcional e indica a quantidade de cotações a considerar (ordem decrescente de criação). Se omitido, considera todas as cotações.

**Parâmetros:**

- `last_quotes` (opcional): inteiro positivo (ex.: `10` para as últimas 10 cotações).

**Resposta de sucesso (200):**

```json
{
  "by_carrier": [
    {
      "carrier_name": "Correios",
      "total_quotes": 5,
      "total_freight": 104.95,
      "average_freight": 20.99
    },
    {
      "carrier_name": "EXPRESSO FR",
      "total_quotes": 5,
      "total_freight": 85,
      "average_freight": 17
    }
  ],
  "cheapest_overall": 17,
  "most_expensive_overall": 20.99
}
```

**Exemplos de erro:**

- **400** – `last_quotes` informado mas não é um inteiro positivo.
- **500** – Erro ao consultar o banco.

## Exemplos de requisição (curl)

### POST /quote

```bash
curl -X POST http://localhost:8080/quote \
  -H "Content-Type: application/json" \
  -d '{
    "recipient": {
      "address": {
        "zipcode": "01311000"
      }
    },
    "volumes": [
      {
        "category": 7,
        "amount": 1,
        "unitary_weight": 5,
        "price": 349,
        "sku": "abc-teste-123",
        "height": 0.2,
        "width": 0.2,
        "length": 0.2
      }
    ]
  }'
```

### GET /metrics (todas as cotações)

```bash
curl http://localhost:8080/metrics
```

### GET /metrics (últimas 5 cotações)

```bash
curl "http://localhost:8080/metrics?last_quotes=5"
```

## Como testar a API

Após subir os containers, use os exemplos de curl abaixo ou o guia **[COMO_TESTAR.md](COMO_TESTAR.md)** (inclui PowerShell e testes de validação).

## Testes automatizados (go test)

Na raiz do projeto:

```bash
go test ./...
```

Para testes com cobertura:

```bash
go test -cover ./...
```

### Aplicação de TDD (exigência do projeto)

O desenvolvimento seguiu **TDD**: para cada funcionalidade, o **teste foi escrito primeiro** (ou em paridade) para definir o comportamento esperado; em seguida a implementação foi feita para fazer o teste passar (ciclo **red → green → refactor**).

**Ciclo utilizado:**  
1. **Red** – Escrever um teste que falha (comportamento desejado).  
2. **Green** – Implementar o mínimo necessário para o teste passar.  
3. **Refactor** – Melhorar o código mantendo os testes verdes.

**Mapeamento comportamento ↔ teste:**

| Comportamento | Teste que define |
|---------------|------------------|
| POST /quote com CEP e volumes válidos retorna ofertas e persiste | `TestQuoteService_CreateQuote_ValidZipcode` |
| CEP com menos de 8 caracteres → erro | `TestQuoteService_CreateQuote_InvalidZipcode_Length` |
| CEP com letras → erro | `TestQuoteService_CreateQuote_InvalidZipcode_NonNumeric` |
| JSON inválido no POST /quote → 400 | `TestQuoteHandler_CreateQuote_InvalidJSON` |
| Zipcode ausente no body → 400 | `TestQuoteHandler_CreateQuote_ValidationError_MissingZipcode` |
| GET /metrics com last_quotes inválido (abc, -1, 0) → 400 | `TestMetricsService_GetMetrics_InvalidLastQuotes`, `TestMetricsHandler_GetMetrics_InvalidLastQuotes` |
| GET /metrics com last_quotes válido retorna métricas | `TestMetricsService_GetMetrics_ValidLastQuotes` |

Os testes usam **AAA** (Arrange-Act-Assert), nomes descritivos e **mocks** (repositório, cliente HTTP) para isolar a unidade testada.

## Estrutura do projeto

```
.
├── cmd/api/main.go          # Entrada da aplicação
├── internal/
│   ├── config/               # Configuração (env)
│   ├── domain/               # Entidades e DTOs
│   ├── client/               # Cliente HTTP Frete Rápido
│   ├── repository/           # Persistência (PostgreSQL)
│   ├── service/              # Regras de negócio
│   └── handler/              # Handlers HTTP (Gin)
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## Banco de dados

As tabelas são criadas automaticamente na subida da API (se não existirem):

- **quotes**: id (UUID), zipcode, created_at
- **quote_offers**: id (UUID), quote_id (FK), carrier_name, service, deadline_days, final_price

As cotações retornadas pelo POST /quote são gravadas em `quotes` e `quote_offers` e usadas pelo GET /metrics.

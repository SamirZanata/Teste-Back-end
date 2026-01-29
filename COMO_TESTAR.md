# Como testar a API

Com os containers rodando (`docker-compose up -d`), a API fica em **http://localhost:8080**.

---

## 1. POST /quote – Criar cotação

Envia o JSON de entrada e recebe as ofertas das transportadoras (dados gravados no banco).

**PowerShell:**

```powershell
curl -X POST http://localhost:8080/quote `
  -H "Content-Type: application/json" `
  -d "@curl_data.json"
```

**Ou com o JSON inline:**

```powershell
curl -X POST http://localhost:8080/quote -H "Content-Type: application/json" -d "{\"recipient\":{\"address\":{\"zipcode\":\"01311000\"}},\"volumes\":[{\"category\":7,\"amount\":1,\"unitary_weight\":5,\"price\":349,\"sku\":\"abc-teste-123\",\"height\":0.2,\"width\":0.2,\"length\":0.2}]}"
```

**Resposta esperada (200):** JSON com `carrier` (lista de transportadoras com `name`, `service`, `deadline`, `price`).

Rode algumas vezes para gerar cotações e poder testar as métricas.

---

## 2. GET /metrics – Consultar métricas

Retorna métricas das cotações gravadas.

**Todas as cotações:**

```powershell
curl http://localhost:8080/metrics
```

**Últimas N cotações (ex.: 5):**

```powershell
curl "http://localhost:8080/metrics?last_quotes=5"
```

**Resposta esperada (200):** JSON com `by_carrier`, `cheapest_overall` e `most_expensive_overall`.

---

## 3. Testes de validação (POST /quote)

**Zipcode inválido (deve retornar 400):**

```powershell
curl -X POST http://localhost:8080/quote -H "Content-Type: application/json" -d "{\"recipient\":{\"address\":{\"zipcode\":\"123\"}},\"volumes\":[{\"category\":7,\"amount\":1,\"unitary_weight\":5,\"price\":349,\"height\":0.2,\"width\":0.2,\"length\":0.2}]}"
```

**last_quotes inválido (deve retornar 400):**

```powershell
curl "http://localhost:8080/metrics?last_quotes=abc"
```

---

## 4. Postman – JSON dinâmico

### Opção A: Variáveis do Postman (recomendado)

Você altera os valores na aba **Variables** (Environment ou Collection) e o body usa `{{nome}}`.

1. Crie um **Environment** (ícone engrenagem → Environments → Add → ex.: "Quote API").
2. Adicione as variáveis (exemplo):

| Variable    | Initial Value | Current Value |
|------------|---------------|---------------|
| zipcode    | 01311000      | 01311000      |
| category   | 7             | 7             |
| amount     | 1             | 1             |
| weight     | 5             | 5             |
| price      | 349           | 349           |
| sku        | abc-teste-123 | abc-teste-123 |
| height     | 0.2           | 0.2           |
| width      | 0.2           | 0.2           |
| length     | 0.2           | 0.2           |

3. Selecione esse Environment no canto superior direito do Postman.
4. Nova request **POST** → URL: `http://localhost:8080/quote`
5. **Headers**: `Content-Type` = `application/json`
6. **Body** → raw → JSON; use o body abaixo (as variáveis são substituídas ao enviar):

```json
{
  "recipient": {
    "address": {
      "zipcode": "{{zipcode}}"
    }
  },
  "volumes": [
    {
      "category": {{category}},
      "amount": {{amount}},
      "unitary_weight": {{weight}},
      "price": {{price}},
      "sku": "{{sku}}",
      "height": {{height}},
      "width": {{width}},
      "length": {{length}}
    }
  ]
}
```

Para testar outro CEP ou outro volume, altere só os valores em **Variables** e clique em **Send**. Não precisa editar o JSON.

---

### Opção B: Editar o JSON direto no Body

1. Método: **POST**
2. URL: `http://localhost:8080/quote`
3. Aba **Headers**: adicione `Content-Type` = `application/json`
4. Aba **Body**: escolha **raw** → tipo **JSON**
5. Cole o JSON abaixo e altere o que quiser (CEP, volumes, preços, etc.):

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
    },
    {
      "category": 7,
      "amount": 2,
      "unitary_weight": 4,
      "price": 556,
      "sku": "abc-teste-527",
      "height": 0.4,
      "width": 0.6,
      "length": 0.15
    }
  ]
}
```

**O que você pode alterar no Postman:**

| Campo | Regra | Exemplo |
|-------|--------|---------|
| `recipient.address.zipcode` | 8 dígitos numéricos | `"01311000"`, `"22041080"` |
| `volumes` | Array com pelo menos 1 item | Adicione ou remova itens |
| Cada volume: `category` | Número ≥ 1 | `7` |
| `amount` | Número ≥ 1 | `1`, `3` |
| `unitary_weight` | Número > 0 (kg) | `5`, `2.5` |
| `price` | Número ≥ 0 (preço unitário) | `349`, `100` |
| `sku` | Opcional | `"meu-sku"` ou `""` |
| `height`, `width`, `length` | Números > 0 (metros) | `0.2`, `0.5` |

Clique em **Send** para enviar. A resposta traz `carrier` (ofertas das transportadoras).

### GET /metrics (métricas)

1. Método: **GET**
2. URL: `http://localhost:8080/metrics`
3. Aba **Params** (opcional): adicione `last_quotes` = `5` (ou outro número; deixe em branco para todas as cotações)
4. **Send**

---

## Ordem sugerida

1. `docker-compose up -d`
2. `curl -X POST ... -d "@curl_data.json"` (2 ou 3 vezes)
3. `curl http://localhost:8080/metrics`
4. `curl "http://localhost:8080/metrics?last_quotes=2"`

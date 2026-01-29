# Enviar o projeto ao GitHub em vários commits

Use este guia **dentro da pasta do projeto Back-end**, com um repositório Git inicializado só para ela (não o repositório da pasta pai).

---

## 1. Garantir que o repo é só do Back-end

Se o Git estiver na pasta pai (ex.: Desktop/Testes), crie um repo novo só para a API:

```powershell
cd "c:\Users\Samir\OneDrive\Desktop\Testes\Back-end"
git init
```

Se já tiver feito commits antes e quiser **desfazer o último commit** e refazer em vários:

```powershell
git reset --soft HEAD~1
```

Assim tudo que estava no último commit volta para "staged"; daí você dá `git restore --staged .` e vai adicionando em partes (passos abaixo).

---

## 2. Adicionar o remote do GitHub

Crie um repositório vazio no GitHub (sem README, sem .gitignore). Depois:

```powershell
git remote add origin https://github.com/SEU_USUARIO/SEU_REPO.git
```

Troque `SEU_USUARIO` e `SEU_REPO` pelo seu usuário e nome do repositório.

---

## 3. Vários commits em sequência

Rode cada bloco **na ordem**. Assim o histórico fica com commits pequenos e lógicos.

**Commit 1 – Estrutura e dependências**
```powershell
git add .gitignore go.mod go.sum
git commit -m "chore: add go module and gitignore"
```

**Commit 2 – Domain**
```powershell
git add internal/domain/
git commit -m "feat(domain): add quote and metrics models"
```

**Commit 3 – Config**
```powershell
git add internal/config/
git commit -m "feat(config): add configuration loader"
```

**Commit 4 – Repository**
```powershell
git add internal/repository/
git commit -m "feat(repository): add quote repository and postgres implementation"
```

**Commit 5 – Client Frete Rápido**
```powershell
git add internal/client/
git commit -m "feat(client): add Frete Rápido API client"
```

**Commit 6 – Services**
```powershell
git add internal/service/quote_service.go internal/service/metrics_service.go
git commit -m "feat(service): add quote and metrics services"
```

**Commit 7 – Handlers**
```powershell
git add internal/handler/quote_handler.go internal/handler/metrics_handler.go
git commit -m "feat(handler): add quote and metrics HTTP handlers"
```

**Commit 8 – Main e rotas**
```powershell
git add cmd/
git commit -m "feat(api): add main and routes"
```

**Commit 9 – Docker**
```powershell
git add Dockerfile docker-compose.yml .dockerignore
git commit -m "chore(docker): add Dockerfile and docker-compose"
```

**Commit 10 – Testes**
```powershell
git add internal/service/*_test.go internal/handler/*_test.go
git commit -m "test: add unit tests for services and handlers"
```

**Commit 11 – Documentação e dados de exemplo**
```powershell
git add README.md COMO_TESTAR.md GUIA_GIT_COMMITS.md .env.example curl_data.json
git commit -m "docs: add README, testing guide and sample data"
```

---

## 4. Enviar para o GitHub

```powershell
git branch -M main
git push -u origin main
```

Se o repo no GitHub já tiver sido criado com outra branch (ex.: `master`), use o nome dela em vez de `main`.

---

## Resumo

- 11 commits com mensagens claras (chore / feat / test / docs).
- Nada de “subir tudo de uma vez”; o histórico fica natural.
- O `.gitignore` evita subir pastas como `go-build`, `go-mod` e arquivos `.env`.

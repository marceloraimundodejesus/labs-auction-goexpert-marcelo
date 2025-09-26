````markdown
# Labs Auction (Go Expert)

API de leilões em Go. Destaques:

- Fechamento automático de leilões via **goroutine**, configurado por `AUCTION_DURATION` e `AUCTION_INTERVAL`;
- Lances em **batch** com `BATCH_INSERT_INTERVAL` e `MAX_BATCH_SIZE`;
- MongoDB como persistência;
- Endpoints REST com Gin.

## Requisitos

- Go 1.20+
- Docker e Docker Compose (opcional para subir MongoDB)
- PowerShell / curl / Postman para testar

## Variáveis de ambiente

| Variável                | Exemplo                                                           | Observação                                |
| ----------------------- | ----------------------------------------------------------------- | ----------------------------------------- |
| `MONGODB_URL`           | `mongodb://admin:admin@127.0.0.1:27017/auctions?authSource=admin` | Conexão com o Mongo                       |
| `BATCH_INSERT_INTERVAL` | `20s`                                                             | Tempo p/ flush do batch de lances         |
| `MAX_BATCH_SIZE`        | `4`                                                               | Qtde máxima no buffer antes do flush      |
| `AUCTION_DURATION`      | `5m`                                                              | Duração total do leilão                   |
| `AUCTION_INTERVAL`      | `20s`                                                             | Intervalo entre checagens do encerramento |

Você pode usar um `.env` em `cmd/auction/.env` (carregado em runtime) ou exportar as variáveis no terminal.

## Como rodar (local, sem Docker para app)

Suba o MongoDB (pode ser via Docker) e rode a API localmente:

```powershell
# Exemplo de envs para desenvolvimento
$env:MONGODB_URL="mongodb://admin:admin@127.0.0.1:27017/auctions?authSource=admin"
$env:BATCH_INSERT_INTERVAL="20s"
$env:MAX_BATCH_SIZE="4"
$env:AUCTION_DURATION="5m"
$env:AUCTION_INTERVAL="20s"

go run cmd/auction/main.go
```
````

### Subindo MongoDB via Docker

```bash
docker compose up -d mongodb
# ou: docker run -d -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=admin --name mongodb mongo:latest
```

> Se a porta `:8080` estiver ocupada, finalize o processo em uso (`netstat -ano | findstr :8080` → `taskkill /PID <pid> /F`) ou ajuste a porta no `main.go`/proxy.

## Endpoints

- `POST /auction` — cria leilão
- `GET /auction?status={0|1}&category={opc}&productName={opc}` — lista leilões

  - `status=0` → **Active**, `status=1` → **Completed**

- `GET /auction/:auctionId` — detalhe do leilão
- `GET /auction/winner/:auctionId` — leilão + lance vencedor (se houver)
- `POST /bid` — cria lance (entra em batch)
- `GET /bid/:auctionId` — lances de um leilão

### Exemplos (PowerShell)

```powershell
# Criar leilão
$body = @{
  product_name = "iPhone 13"
  category     = "electronics"
  description  = "iPhone 13 128GB em ótimo estado"
  condition    = 2
} | ConvertTo-Json -Compress
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/auction" -ContentType "application/json" -Body $body

# Listar abertos
Invoke-RestMethod "http://localhost:8080/auction?status=0"

# Criar lance
$auctionId = (Invoke-RestMethod "http://localhost:8080/auction?status=0")[-1].id
$bid = @{ user_id = ([guid]::NewGuid().Guid); auction_id = $auctionId; amount = 999.99 } | ConvertTo-Json -Compress
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/bid" -ContentType "application/json" -Body $bid
```

## Smoke test (fechamento automático)

Para facilitar a validação do winner com lance, use `BATCH_INSERT_INTERVAL` **menor** que `AUCTION_DURATION`:

```powershell
$env:BATCH_INSERT_INTERVAL="2s"
$env:MAX_BATCH_SIZE="4"
$env:AUCTION_DURATION="20s"
$env:AUCTION_INTERVAL="2s"
$env:MONGODB_URL="mongodb://admin:admin@127.0.0.1:27017/auctions?authSource=admin"
go run cmd/auction/main.go
```

Em outro terminal:

```powershell
# Cria leilão
$body = @{ product_name="AutoClose With Bid"; category="electronics"; description="desc 1234567890"; condition=2 } | ConvertTo-Json -Compress
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/auction" -ContentType "application/json" -Body $body
$auctionId = (Invoke-RestMethod "http://localhost:8080/auction?status=0")[-1].id

# Dá 1 lance (flush rápido em ~2s)
$bid = @{ user_id = ([guid]::NewGuid().Guid); auction_id = $auctionId; amount = 777.77 } | ConvertTo-Json -Compress
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/bid" -ContentType "application/json" -Body $bid

Start-Sleep -Seconds 5
Invoke-RestMethod "http://localhost:8080/bid/$auctionId" | ConvertTo-Json -Compress   # deve listar o lance

Start-Sleep -Seconds 20
Invoke-RestMethod "http://localhost:8080/auction/winner/$auctionId" | ConvertTo-Json -Compress  # auction + bid
```

## Desenvolvimento

- `go build ./...`
- `go test ./...`

## Solução de problemas

- **Porta 8080 ocupada**: identifique e encerre o processo (`netstat -ano | findstr :8080` → `taskkill /PID <pid> /F`).
- **Falha ao conectar no Mongo (localhost vs 127.0.0.1)**: em Windows, prefira `127.0.0.1` para evitar `::1` (IPv6).
- **Winner sem bid**: se `BATCH_INSERT_INTERVAL` for igual/maior que `AUCTION_DURATION`, o flush pode ocorrer após o fechamento. Reduza o intervalo ou envie lances suficientes para estourar `MAX_BATCH_SIZE`.

---

```

```

# AdotaPet

Backend e banco da Fase 1 do AdotaPet, derivados dos artefatos em `../docs/fase-1`.

## Desenvolvimento Local

1. Copie `.env.example` para `.env` e preencha os valores locais obrigatorios:
   `POSTGRES_PASSWORD`, `DATABASE_URL` e `JWT_ACCESS_SECRET`.
   Os TTLs dos tokens podem ser ajustados com `JWT_ACCESS_TTL` e
   `JWT_REFRESH_TTL`, usando formato do Go como `15m`, `1h` ou `720h`.
   O TTL dos codigos de verificacao usa `VERIFICATION_CODE_TTL`.
2. Entre no diretorio do backend:

```sh
cd backend
```

3. Suba PostgreSQL/PostGIS e Redis:

```sh
make db-up
```

4. Execute as migrations:

```sh
make migrate
```

5. Inicie a API:

```sh
make run
```

## Endpoints iniciais

- `GET /health/live`
- `GET /health/ready`
- Rotas de `/api/v1` registradas como placeholders para a primeira fatia vertical.

## Estrutura

- `cmd/api`: entrypoint HTTP.
- `internal/adapters/inbound/http`: handlers, middleware e mapeamento de erros.
- `internal/app`: use cases e portas.
- `internal/domain`: entidades e tipos de dominio.
- `internal/adapters/outbound`: adapters de PostgreSQL, Redis, storage e notificacoes.
- `db/migration`: migrations Flyway.
- `deployments/local`: infraestrutura local.

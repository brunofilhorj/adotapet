# AdotaPet

Repositorio organizado por areas do produto.

## Estrutura

```text
docs/
  Arquitetura, HLD, LLD e diagramas por fase.

backend/
  API Go, migrations, infraestrutura local e adapters.

frontend/
  mobile/
    Aplicativo mobile.
  web/
    Aplicacao web futura.
  extranet/
    Portal externo/futuro para parceiros, abrigos ou operacao.
```

## Desenvolvimento

O backend da Fase 1 esta em `backend/`.

```sh
cd backend
make db-up
make migrate
make run
```

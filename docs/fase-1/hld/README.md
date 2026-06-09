# AdotaPet — High Level Design (HLD)

Documento de arquitetura de alto nível para o aplicativo de busca e adoção de filhotes.

## Como abrir no draw.io

1. Acesse [https://app.diagrams.net](https://app.diagrams.net) (draw.io)
2. **File → Open from → Device**
3. Selecione o arquivo `adotapet-hld.drawio`
4. Navegue pelas abas na parte inferior do editor

Também é possível abrir diretamente no VS Code/Cursor com a extensão **Draw.io Integration**.

### Problemas ao abrir?

Se aparecer erro de parsing, tente:
- Arrastar o arquivo diretamente para a janela do draw.io
- Usar **File → Import from → Device** em vez de Open
- Abrir via extensão Draw.io Integration no Cursor (geralmente mais tolerante)

## Diagramas incluídos

| Aba | Conteúdo |
|-----|----------|
| **1. Contexto do Sistema** | Visão C4 Nível 1 — atores (adotante, doador) e integrações externas |
| **2. Contêineres** | C4 Nível 2 — app mobile, backend Go, PostgreSQL, Redis, S3 |
| **3. Arquitetura Hexagonal** | Ports & Adapters — camadas domain, application, infrastructure |
| **4. Modelo de Dados** | Entidades PostgreSQL + PostGIS, relacionamentos e índices |
| **5. Fluxos Principais** | Busca geolocalizada, cadastro de filhote, chat, autenticação, API REST |
| **6. Telas Mobile** | Wireframe de navegação e mapa de telas do app |

## Decisões arquiteturais (ADRs resumidas)

### Stack

| Camada | Tecnologia | Justificativa |
|--------|-----------|---------------|
| Mobile | Flutter (recomendado) | Uma base de código para iOS e Android, UI rica, bom suporte a mapas/GPS |
| Backend | Go | Backend simples de operar, binário único, boa performance, concorrência nativa para chat e notificações |
| Banco | PostgreSQL 16 + PostGIS | Busca geoespacial nativa, JSONB para filtros salvos, ACID |
| Cache | Redis | Cache de buscas frequentes, sessões WebSocket |
| Mídia | S3 + CDN | Upload via presigned URL, escalável |
| Auth | JWT (access + refresh) | Stateless, padrão para mobile |
| Chat | WebSocket + persistência | Tempo real com histórico em PostgreSQL |
| Push | FCM + APNs | Notificações em ambas plataformas |

### Arquitetura Hexagonal (Backend)

```
cmd/api/                   → entrypoint HTTP/WebSocket
internal/adapters/inbound/ → driving adapters (REST handlers, WebSocket)
internal/app/              → use cases + port interfaces
internal/domain/           → entities, value objects, domain services (Go puro)
internal/adapters/outbound/→ driven adapters (PostgreSQL, S3, Redis, FCM, Geocoding)
```

**Regra de dependência:** `cmd/api → inbound adapters → app → domain`; adapters outbound implementam portas da camada `app`.

O pacote `internal/domain` não depende de frameworks, banco, Redis, S3 ou HTTP.

### Módulos de domínio (monólito modular)

- **Usuários** — cadastro, perfis, roles (ADOPTER, DONOR, SHELTER)
- **Anúncios** — CRUD de filhotes, upload de fotos, status de adoção
- **Busca** — filtros avançados + geolocalização via PostGIS
- **Mensagens** — conversas 1:1, WebSocket, notificações
- **Notificações** — push, e-mail, alertas de busca salva

## Próximos passos

1. **Scaffolding** — criar módulo Go do backend e projeto Flutter
2. **MVP** — implementar autenticação, anúncios, busca geolocalizada, chat e notificações básicas
3. **Fase 2** — implementar avaliações de tutores/doadores e reputação pública
4. **Fase 3** — implementar perguntas e respostas públicas nos anúncios

# AdotaPet - Documentação

Documentação arquitetural organizada por fases de evolução do produto.

## Fases

| Fase | Objetivo | HLD | LLD |
|------|----------|-----|-----|
| Fase 1 | MVP: cadastro, anúncios, busca geolocalizada, chat e notificações básicas | `fase-1/hld` | `fase-1/lld` |
| Fase 2 | Rede de confiança: avaliações de tutores/doadores e reputação pública | `fase-2/hld` | `fase-2/lld` |
| Fase 3 | Transparência nos anúncios: perguntas e respostas públicas por filhote | `fase-3/hld` | `fase-3/lld` |

## Convenção

Cada fase possui:

- `hld/README.md`: arquitetura de alto nível e decisões.
- `hld/*.drawio`: diagrama visual de alto nível.
- `lld/README.md`: design de baixo nível, contratos, dados e regras.
- `lld/*.drawio`: diagrama visual de baixo nível.

## Dependência Entre Fases

```text
Fase 1: base operacional do app
  ↓
Fase 2: reputação de tutores/doadores
  ↓
Fase 3: perguntas públicas nos anúncios
```

A Fase 2 depende de conversas, usuários e anúncios da Fase 1. A Fase 3 depende de usuários, anúncios e moderação básica, e pode reutilizar sinais de reputação da Fase 2.

# AdotaPet - Fase 3 HLD

High Level Design da Fase 3: perguntas e respostas públicas nos anúncios.

## Objetivo

Aumentar transparência e reduzir perguntas repetidas no chat privado permitindo perguntas públicas no perfil do filhote, com respostas oficiais do tutor/doador.

## Escopo

- Usuários autenticados podem perguntar em anúncios disponíveis.
- Dono do anúncio responde publicamente.
- Perguntas e respostas ficam visíveis no perfil do filhote.
- Usuários podem denunciar perguntas ou respostas.
- Owner pode ocultar pergunta inadequada com trilha de auditoria.

Fora de escopo:

- Comentários livres em formato de mural social.
- Respostas por usuários que não sejam o dono do anúncio.
- Reações, likes ou thread aninhada.
- Chat público em tempo real.

## Componentes

```text
Mobile Flutter
  - Seção "Perguntas" no perfil do filhote
  - Formulário de nova pergunta
  - Tela de resposta para owner

Backend Go
  - Puppy question handlers
  - Question and answer use cases
  - Moderation/report hooks

PostgreSQL
  - puppy_questions

Notificações
  - owner recebe aviso de nova pergunta
  - autor recebe aviso quando pergunta é respondida
```

## Decisões

| Tema | Decisão |
|------|---------|
| Formato | Perguntas e respostas, não comentários livres |
| Autoridade | Apenas dono do anúncio pode responder |
| Visibilidade | Perguntas `VISIBLE` aparecem publicamente |
| Moderação | Status `VISIBLE`, `HIDDEN`, `REPORTED` |
| Segurança | Bloquear dados pessoais e permitir denúncia |

## Relação Com Outras Fases

- Usa usuários, anúncios e notificações da Fase 1.
- Pode usar reputação da Fase 2 para destacar tutores confiáveis.
- Complementa o chat privado, mas não substitui conversa sensível.

## Diagrama

Abra `adotapet-fase-3-hld.drawio` no diagrams.net ou na extensão Draw.io.

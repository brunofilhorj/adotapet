# AdotaPet - Fase 2 HLD

High Level Design da Fase 2: reputação e avaliações de tutores/doadores.

## Objetivo

Criar uma rede de confiança dentro do AdotaPet permitindo que usuários avaliem tutores, doadores e abrigos após uma interação qualificada.

## Escopo

- Avaliação de tutor/doador após conversa ou adoção marcada.
- Nota de 1 a 5.
- Comentário opcional.
- Tags rápidas de experiência.
- Reputação agregada no perfil público.
- Denúncia e ocultação de avaliações inadequadas.

Fora de escopo:

- Avaliação livre sem vínculo entre usuários.
- Ranking público global de tutores.
- Resposta pública do avaliado, inicialmente.
- Moderação automatizada por IA.

## Componentes

```text
Mobile Flutter
  - Tela de avaliar tutor
  - Card de reputação no perfil
  - Lista de avaliações recebidas

Backend Go
  - Review handlers
  - Review use cases
  - Reputation summary use cases
  - Moderation hooks

PostgreSQL
  - tutor_reviews
  - tutor_reputation_summaries

Notificações
  - aviso de nova avaliação recebida
  - aviso de avaliação ocultada/reportada
```

## Regras De Alto Nível

- Um usuário só pode avaliar outro se houver interação qualificada.
- Interação qualificada pode ser conversa criada, mensagem trocada ou adoção marcada.
- O mesmo avaliador só pode criar uma avaliação por contexto de adoção/conversa.
- A média pública aparece apenas após quantidade mínima de avaliações.
- Avaliações reportadas podem ficar ocultas até revisão.

## Experiência Mobile

- Após uma adoção marcada ou conversa relevante, o app sugere avaliar o tutor.
- Perfil público mostra nota média, quantidade de avaliações e tags mais recorrentes.
- Avaliações textuais aparecem em lista paginada.
- Usuários podem denunciar uma avaliação.

## Decisões

| Tema | Decisão |
|------|---------|
| Modelo | Avaliação do usuário/tutor, com vínculo opcional ao filhote e conversa |
| Visibilidade | Média pública apenas com mínimo de 3 avaliações |
| Segurança | Avaliação exige vínculo real entre usuários |
| Moderação | Status `PUBLISHED`, `HIDDEN`, `REPORTED` |
| Agregação | Sumário materializado para evitar cálculo caro no perfil |

## Diagrama

Abra `adotapet-fase-2-hld.drawio` no diagrams.net ou na extensão Draw.io.

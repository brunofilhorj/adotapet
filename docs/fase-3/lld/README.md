# AdotaPet - Fase 3 LLD

Low Level Design da Fase 3: perguntas e respostas públicas por anúncio.

## 1. Domínio

### PuppyQuestion

```text
id: UUID
puppyId: UUID
authorId: UUID
question: String
answer: String?
answeredByOwnerId: UUID?
status: QuestionStatus
createdAt: Instant
answeredAt: Instant?
reportedAt: Instant?
hiddenAt: Instant?
```

Enums:

```text
QuestionStatus = VISIBLE | HIDDEN | REPORTED
```

Regras:

- `question` deve ter entre 5 e 500 caracteres.
- `answer`, quando preenchida, deve ter entre 1 e 1000 caracteres.
- Apenas usuário autenticado pode perguntar.
- Apenas o dono do anúncio pode responder.
- Autor da pergunta não pode responder sua própria pergunta, exceto se também for o dono, cenário que deve ser bloqueado para evitar autopromoção.
- Perguntas em anúncios `REMOVED` não são permitidas.
- Perguntas reportadas podem ficar ocultas até revisão.

## 2. Use Cases

```text
CreatePuppyQuestionUseCase
  input: puppyId, authorId, question
  output: questionId, status
  ports: PuppyRepository, PuppyQuestionRepository, NotificationPort

AnswerPuppyQuestionUseCase
  input: questionId, ownerId, answer
  output: questionId, answeredAt
  ports: PuppyRepository, PuppyQuestionRepository, NotificationPort

ListPuppyQuestionsUseCase
  input: puppyId, page, size
  output: perguntas visiveis paginadas

ReportPuppyQuestionUseCase
  input: questionId, reporterId, reason
  output: status

HidePuppyQuestionUseCase
  input: questionId, actorId
  output: status
```

## 3. Portas Go

```go
type PuppyQuestionRepository interface {
	Save(ctx context.Context, question PuppyQuestion) (PuppyQuestion, error)
	FindByID(ctx context.Context, id uuid.UUID) (*PuppyQuestion, error)
	FindByPuppyID(ctx context.Context, puppyID uuid.UUID, page PageRequest) (Page[PuppyQuestion], error)
	UpdateAnswer(ctx context.Context, id uuid.UUID, ownerID uuid.UUID, answer string) (PuppyQuestion, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status QuestionStatus) (PuppyQuestion, error)
}
```

## 4. API REST

Base path: `/api/v1`

```text
POST   /puppies/{id}/questions
GET    /puppies/{id}/questions
POST   /puppy-questions/{id}/answer
POST   /puppy-questions/{id}/report
PATCH  /puppy-questions/{id}/status
```

### POST `/puppies/{id}/questions`

Request:

```json
{
  "question": "Ela convive bem com gatos?"
}
```

Response `201`:

```json
{
  "id": "ed332a51-ff69-4287-a6c9-7be771eabc55",
  "status": "VISIBLE"
}
```

### POST `/puppy-questions/{id}/answer`

Request:

```json
{
  "answer": "Sim, ela convive com dois gatos no lar temporario."
}
```

Response `200`:

```json
{
  "id": "ed332a51-ff69-4287-a6c9-7be771eabc55",
  "answeredAt": "2026-06-08T12:30:00Z"
}
```

### GET `/puppies/{id}/questions`

Response `200`:

```json
{
  "items": [
    {
      "id": "ed332a51-ff69-4287-a6c9-7be771eabc55",
      "author": {
        "id": "4a52c5cb-1cd6-49ed-9423-3c918bba6c13",
        "name": "Maria Souza"
      },
      "question": "Ela convive bem com gatos?",
      "answer": "Sim, ela convive com dois gatos no lar temporario.",
      "answeredAt": "2026-06-08T12:30:00Z",
      "createdAt": "2026-06-08T11:30:00Z"
    }
  ],
  "page": 0,
  "size": 20,
  "totalElements": 1
}
```

## 5. Banco De Dados

```sql
CREATE TABLE puppy_questions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  puppy_id UUID NOT NULL REFERENCES puppies(id) ON DELETE CASCADE,
  author_id UUID NOT NULL REFERENCES users(id),
  question TEXT NOT NULL,
  answer TEXT,
  answered_by_owner_id UUID REFERENCES users(id),
  status VARCHAR(20) NOT NULL CHECK (status IN ('VISIBLE', 'HIDDEN', 'REPORTED')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  answered_at TIMESTAMPTZ,
  reported_at TIMESTAMPTZ,
  hidden_at TIMESTAMPTZ
);

CREATE INDEX idx_puppy_questions_puppy_status_created
ON puppy_questions (puppy_id, status, created_at DESC);

CREATE INDEX idx_puppy_questions_author_id
ON puppy_questions (author_id);
```

## 6. Critérios De Aceite

- Usuário autenticado consegue perguntar em anúncio elegível.
- Dono do anúncio recebe notificação de nova pergunta.
- Apenas dono do anúncio consegue responder.
- Perguntas respondidas aparecem no perfil do filhote.
- Perguntas reportadas ou ocultas não aparecem publicamente.
- Comentários livres ou threads aninhadas não são criados nesta fase.

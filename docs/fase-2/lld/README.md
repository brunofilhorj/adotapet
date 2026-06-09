# AdotaPet - Fase 2 LLD

Low Level Design da Fase 2: avaliações de tutores/doadores e reputação.

## 1. Domínio

### TutorReview

```text
id: UUID
reviewerId: UUID
reviewedUserId: UUID
puppyId: UUID?
conversationId: UUID?
rating: Int
comment: String?
tags: ReviewTag[]
status: ReviewStatus
createdAt: Instant
updatedAt: Instant
reportedAt: Instant?
hiddenAt: Instant?
```

Enums:

```text
ReviewStatus = PUBLISHED | HIDDEN | REPORTED

ReviewTag =
  RESPONSIVE
  CAREFUL
  CLEAR_INFORMATION
  SAFE_ENVIRONMENT
  DID_NOT_RESPOND
  INCONSISTENT_PROFILE
```

Regras:

- `rating` deve estar entre 1 e 5.
- `reviewerId` e `reviewedUserId` devem ser diferentes.
- Deve existir vínculo qualificado entre avaliador e avaliado.
- `comment` é opcional, mas se existir deve ter no máximo 1000 caracteres.
- Uma avaliação reportada não entra no sumário público até revisão.

### TutorReputationSummary

```text
userId: UUID
averageRating: Decimal
reviewCount: Int
publishedReviewCount: Int
tagCounts: JSONB
lastReviewAt: Instant?
updatedAt: Instant
```

Regras:

- Sumário é recalculado após criação, ocultação ou republicação de avaliação.
- Média pública só aparece se `publishedReviewCount >= 3`.

## 2. Use Cases

```text
CreateTutorReviewUseCase
  input: reviewerId, reviewedUserId, puppyId?, conversationId?, rating, comment?, tags
  output: reviewId, status
  ports: TutorReviewRepository, ConversationRepository, PuppyRepository, ReputationSummaryRepository

ListTutorReviewsUseCase
  input: reviewedUserId, page, size
  output: reviews paginadas

GetTutorReputationUseCase
  input: userId
  output: averageRating?, reviewCount, tagCounts

ReportTutorReviewUseCase
  input: reviewId, reporterId, reason
  output: status

HideTutorReviewUseCase
  input: reviewId, moderatorId
  output: status
```

## 3. Portas Go

```go
type TutorReviewRepository interface {
	Save(ctx context.Context, review TutorReview) (TutorReview, error)
	FindByID(ctx context.Context, id uuid.UUID) (*TutorReview, error)
	FindByReviewedUserID(ctx context.Context, userID uuid.UUID, page PageRequest) (Page[TutorReview], error)
	ExistsByContext(ctx context.Context, reviewerID uuid.UUID, reviewedUserID uuid.UUID, puppyID *uuid.UUID, conversationID *uuid.UUID) (bool, error)
}

type ReputationSummaryRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) (*TutorReputationSummary, error)
	Recalculate(ctx context.Context, userID uuid.UUID) (TutorReputationSummary, error)
}
```

## 4. API REST

Base path: `/api/v1`

```text
POST   /tutor-reviews
GET    /users/{id}/reputation
GET    /users/{id}/reviews
POST   /tutor-reviews/{id}/report
PATCH  /admin/tutor-reviews/{id}/status
```

### POST `/tutor-reviews`

Request:

```json
{
  "reviewedUserId": "b3a2c3d4-80de-4efa-aabc-789d66545bd2",
  "puppyId": "7c0d20e4-d0bf-43ad-ac2b-548e43c8be99",
  "conversationId": "c610f228-8655-448d-a8b7-2e131c7da8a4",
  "rating": 5,
  "comment": "Muito atencioso e passou todas as informacoes.",
  "tags": ["RESPONSIVE", "CLEAR_INFORMATION"]
}
```

Response `201`:

```json
{
  "id": "1c0674b9-88e7-4fdd-b5bb-bf65c5567ac8",
  "status": "PUBLISHED"
}
```

### GET `/users/{id}/reputation`

Response `200`:

```json
{
  "userId": "b3a2c3d4-80de-4efa-aabc-789d66545bd2",
  "averageRating": 4.8,
  "reviewCount": 12,
  "publishedReviewCount": 12,
  "topTags": [
    {
      "tag": "RESPONSIVE",
      "count": 9
    }
  ]
}
```

## 5. Banco De Dados

```sql
CREATE TABLE tutor_reviews (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  reviewer_id UUID NOT NULL REFERENCES users(id),
  reviewed_user_id UUID NOT NULL REFERENCES users(id),
  puppy_id UUID REFERENCES puppies(id),
  conversation_id UUID REFERENCES conversations(id),
  rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
  comment TEXT,
  tags JSONB NOT NULL DEFAULT '[]'::jsonb,
  status VARCHAR(20) NOT NULL CHECK (status IN ('PUBLISHED', 'HIDDEN', 'REPORTED')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  reported_at TIMESTAMPTZ,
  hidden_at TIMESTAMPTZ,
  CHECK (reviewer_id <> reviewed_user_id)
);

CREATE UNIQUE INDEX idx_tutor_reviews_unique_context
ON tutor_reviews (reviewer_id, reviewed_user_id, puppy_id, conversation_id);

CREATE INDEX idx_tutor_reviews_reviewed_user_status
ON tutor_reviews (reviewed_user_id, status, created_at DESC);

CREATE TABLE tutor_reputation_summaries (
  user_id UUID PRIMARY KEY REFERENCES users(id),
  average_rating NUMERIC(3,2) NOT NULL DEFAULT 0,
  review_count INT NOT NULL DEFAULT 0,
  published_review_count INT NOT NULL DEFAULT 0,
  tag_counts JSONB NOT NULL DEFAULT '{}'::jsonb,
  last_review_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## 6. Critérios De Aceite

- Usuário consegue avaliar tutor apenas quando existe vínculo qualificado.
- Perfil público exibe reputação quando há no mínimo 3 avaliações publicadas.
- Avaliações reportadas deixam de aparecer publicamente até revisão.
- Sumário de reputação é atualizado após criação ou moderação.
- Não é possível avaliar o mesmo contexto duas vezes.

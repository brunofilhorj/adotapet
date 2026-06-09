# AdotaPet - Low Level Design (LLD)

Documento de baixo nível derivado do HLD em `docs/fase-1/hld/adotapet-hld.drawio`.

## 1. Escopo

Este LLD detalha a implementação do MVP da plataforma AdotaPet:

- Cadastro, autenticação e perfil de usuário.
- Cadastro e manutenção de anúncios de filhotes.
- Busca geolocalizada com filtros.
- Favoritos e buscas salvas.
- Conversas 1:1 e mensagens em tempo real.
- Notificações push e e-mail/SMS de verificação.
- Upload de fotos via object storage.

Fora do MVP inicial:

- Moderação automatizada de conteúdo.
- Pagamentos, planos ou monetização.
- Ranking avançado por machine learning.
- Integração direta com ONGs/abrigos externos.

## 2. Arquitetura De Implementação

### 2.1 Backend

Stack:

- Go.
- HTTP router idiomático, recomendado `chi` ou `Gin`.
- Middleware próprio para autenticação, autorização, CORS e rate limiting.
- Go modules.
- PostgreSQL 16 + PostGIS.
- Redis.
- Flyway.
- WebSocket handler dedicado, recomendado `gorilla/websocket` ou `nhooyr.io/websocket`.
- S3-compatible object storage.
- FCM/APNs para push.

Estrutura Go:

```text
cmd/api/
  entrypoint HTTP/WebSocket, wiring de dependencias, config

internal/adapters/inbound/http/
  REST handlers, request/response DTOs, middleware, error mapping

internal/adapters/inbound/ws/
  WebSocket handler e eventos de chat

internal/app/
  use cases, input ports, output ports, transaction boundaries

internal/domain/
  entities, value objects, domain services, domain errors

internal/adapters/outbound/
  PostgreSQL repositories, Redis adapter, S3 adapter, push adapter, geocoding adapter

pkg/
  bibliotecas reutilizaveis somente quando houver uso fora de internal
```

Regra de dependência:

```text
cmd/api -> inbound adapters -> app -> domain
outbound adapters -> app
outbound adapters -> domain
domain -> nenhuma dependencia de framework, HTTP, banco ou storage
```

### 2.2 Mobile

Stack:

- Flutter.
- REST sobre HTTPS.
- WebSocket sobre WSS.
- Storage local seguro para tokens.
- GPS nativo para busca por proximidade.
- Upload de fotos via URL pré-assinada.

Telas principais:

- Cadastro/Login.
- Home/Explorar.
- Busca avançada.
- Perfil do filhote.
- Favoritos.
- Mensagens.
- Chat.
- Meus anúncios.
- Meu perfil.

## 3. Pacotes Backend

```text
adotapet
  cmd
    api

  internal
    config
    platform
      clock
      logger
      tracing

    adapters
      inbound
        http
          middleware
          auth
          users
          puppies
          search
          favorites
          conversations
          notifications
          dto
          errors
        ws
          chat
      outbound
        postgres
          mapper
          repository
        redis
        storage
        push
        geocoding
        email

    app
      auth
      users
      puppies
      search
      favorites
      conversations
      notifications
      media
      port
        in
        out

    domain
      user
      puppy
      search
      conversation
      notification
      media
      common
```

## 4. Modelo De Domínio

### 4.1 User

Responsabilidades:

- Representar uma conta autenticável.
- Controlar papel e status.
- Vincular perfil público.

Campos:

```text
id: UUID
email: Email
passwordHash: String
role: UserRole
status: UserStatus
createdAt: Instant
updatedAt: Instant
```

Enums:

```text
UserRole = ADOPTER | DONOR | SHELTER
UserStatus = PENDING_VERIFICATION | ACTIVE | SUSPENDED | DELETED
```

Regras:

- E-mail deve ser único.
- Usuários `PENDING_VERIFICATION` podem autenticar apenas para reenviar/verificar código.
- Apenas `DONOR` e `SHELTER` podem criar anúncios.

### 4.2 Profile

Campos:

```text
userId: UUID
name: String
phone: String?
city: String
state: String
location: GeoPoint?
avatarUrl: String?
bio: String?
```

Regras:

- `name`, `city` e `state` são obrigatórios após onboarding.
- `location` é usada como fallback para anúncios sem localização própria.
- Telefone deve ser opcional no perfil público.

### 4.3 Puppy

Campos:

```text
id: UUID
ownerId: UUID
name: String
breed: String?
species: Species
ageMonths: Int
size: PuppySize
sex: Sex
description: String
location: GeoPoint
status: PuppyStatus
createdAt: Instant
updatedAt: Instant
adoptedAt: Instant?
```

Enums:

```text
Species = DOG | CAT | OTHER
PuppySize = SMALL | MEDIUM | LARGE
Sex = MALE | FEMALE | UNKNOWN
PuppyStatus = AVAILABLE | ADOPTED | PAUSED | REMOVED
```

Regras:

- `ageMonths` deve ser maior ou igual a zero.
- Anúncios `REMOVED` não aparecem em busca.
- Anúncios `PAUSED` aparecem apenas para o dono.
- Apenas o dono pode editar ou marcar como adotado.
- Ao marcar `ADOPTED`, `adoptedAt` deve ser preenchido.

### 4.4 PuppyPhoto

Campos:

```text
id: UUID
puppyId: UUID
url: String
sortOrder: Int
isPrimary: Boolean
createdAt: Instant
```

Regras:

- Um anúncio deve ter no mínimo uma foto para ser publicado.
- Apenas uma foto por anúncio pode ser `isPrimary = true`.
- `sortOrder` deve ser único por anúncio.

### 4.5 Conversation

Campos:

```text
id: UUID
puppyId: UUID
adopterId: UUID
donorId: UUID
status: ConversationStatus
createdAt: Instant
updatedAt: Instant
```

Enums:

```text
ConversationStatus = OPEN | CLOSED | BLOCKED
```

Regras:

- Uma conversa é sempre entre adotante e dono do anúncio.
- Deve existir no máximo uma conversa aberta por `puppyId + adopterId`.
- O dono não pode iniciar conversa consigo mesmo.

### 4.6 Message

Campos:

```text
id: UUID
conversationId: UUID
senderId: UUID
content: String
sentAt: Instant
readAt: Instant?
```

Regras:

- `content` deve ter entre 1 e 2000 caracteres.
- `senderId` deve pertencer à conversa.
- Mensagens em conversa `CLOSED` ou `BLOCKED` devem ser rejeitadas.

## 5. Use Cases

### 5.1 Autenticação

```text
RegisterUserUseCase
  input: email, password, role, name, city, state
  output: userId, status
  ports: UserRepository, PasswordHasher, VerificationCodePort, NotificationPort

LoginUseCase
  input: email, password
  output: accessToken, refreshToken, expiresIn
  ports: UserRepository, PasswordHasher, TokenPort

RefreshTokenUseCase
  input: refreshToken
  output: accessToken, refreshToken, expiresIn
  ports: TokenPort, UserRepository

VerifyAccountUseCase
  input: email, code
  output: userId, status
  ports: UserRepository, VerificationCodePort
```

### 5.2 Usuários

```text
GetMyProfileUseCase
UpdateProfileUseCase
UploadAvatarUseCase
DeactivateAccountUseCase
```

### 5.3 Anúncios

```text
CreatePuppyListingUseCase
  input: ownerId, puppy data, address/location, photo references
  output: puppyId, status
  ports: PuppyRepository, GeoLocationPort, MediaStoragePort

UpdatePuppyListingUseCase
PausePuppyListingUseCase
MarkPuppyAsAdoptedUseCase
RemovePuppyListingUseCase
GetPuppyDetailsUseCase
ListMyPuppyListingsUseCase
```

### 5.4 Busca

```text
SearchPuppiesNearbyUseCase
  input: latitude, longitude, radiusKm, filters, page, size, sort
  output: paginated puppy summaries
  ports: SearchRepository, CachePort

SaveSearchUseCase
ListSavedSearchesUseCase
DeleteSavedSearchUseCase
```

### 5.5 Favoritos

```text
AddFavoriteUseCase
RemoveFavoriteUseCase
ListFavoritesUseCase
```

### 5.6 Conversas E Mensagens

```text
StartConversationUseCase
  input: puppyId, adopterId
  output: conversationId
  ports: ConversationRepository, PuppyRepository, NotificationPort

SendMessageUseCase
  input: conversationId, senderId, content
  output: messageId, sentAt
  ports: ConversationRepository, MessageRepository, RealtimePort, NotificationPort

ListConversationsUseCase
ListMessagesUseCase
MarkConversationAsReadUseCase
CloseConversationUseCase
```

## 6. Portas

### 6.1 Input Ports

```go
type RegisterUserInputPort interface {
	Register(ctx context.Context, cmd RegisterUserCommand) (RegisteredUser, error)
}

type SearchPuppiesInputPort interface {
	Search(ctx context.Context, query PuppySearchQuery) (Page[PuppySummary], error)
}

type SendMessageInputPort interface {
	Send(ctx context.Context, cmd SendMessageCommand) (SentMessage, error)
}
```

### 6.2 Output Ports

```go
type UserRepository interface {
	Save(ctx context.Context, user User) (User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email Email) (*User, error)
	ExistsByEmail(ctx context.Context, email Email) (bool, error)
}

type PuppyRepository interface {
	Save(ctx context.Context, puppy Puppy) (Puppy, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Puppy, error)
	FindByOwnerID(ctx context.Context, ownerID uuid.UUID, page PageRequest) (Page[Puppy], error)
}

type SearchRepository interface {
	SearchNearby(ctx context.Context, query PuppySearchQuery) (Page[PuppySummary], error)
}

type MediaStoragePort interface {
	CreateUploadURL(ctx context.Context, cmd CreateUploadURLCommand) (PresignedUpload, error)
	DeleteObject(ctx context.Context, objectKey string) error
}

type NotificationPort interface {
	SendAccountVerification(ctx context.Context, target NotificationTarget, code string) error
	SendNewMessage(ctx context.Context, target NotificationTarget, message MessageNotification) error
	SendSavedSearchMatch(ctx context.Context, target NotificationTarget, match SavedSearchMatch) error
}

type GeoLocationPort interface {
	Geocode(ctx context.Context, address Address) (GeoPoint, error)
}

type CachePort interface {
	Get(ctx context.Context, key string, dest any) (bool, error)
	Put(ctx context.Context, key string, value any, ttl time.Duration) error
	Evict(ctx context.Context, pattern string) error
}
```

## 7. API REST

Base path: `/api/v1`

Todas as rotas protegidas exigem `Authorization: Bearer <access-token>`.

### 7.1 Auth

#### POST `/auth/register`

Request:

```json
{
  "email": "maria@example.com",
  "password": "SenhaForte123!",
  "role": "ADOPTER",
  "name": "Maria Souza",
  "city": "Sao Paulo",
  "state": "SP"
}
```

Response `201`:

```json
{
  "userId": "4a52c5cb-1cd6-49ed-9423-3c918bba6c13",
  "status": "PENDING_VERIFICATION"
}
```

#### POST `/auth/login`

Request:

```json
{
  "email": "maria@example.com",
  "password": "SenhaForte123!"
}
```

Response `200`:

```json
{
  "accessToken": "jwt",
  "refreshToken": "jwt",
  "expiresIn": 900
}
```

#### POST `/auth/verify`

Request:

```json
{
  "email": "maria@example.com",
  "code": "123456"
}
```

Response `200`:

```json
{
  "userId": "4a52c5cb-1cd6-49ed-9423-3c918bba6c13",
  "status": "ACTIVE"
}
```

#### POST `/auth/refresh`

Request:

```json
{
  "refreshToken": "jwt"
}
```

Response `200`: mesmo contrato de `/auth/login`.

### 7.2 Perfil

#### GET `/me`

Response `200`:

```json
{
  "id": "4a52c5cb-1cd6-49ed-9423-3c918bba6c13",
  "email": "maria@example.com",
  "role": "ADOPTER",
  "status": "ACTIVE",
  "profile": {
    "name": "Maria Souza",
    "phone": null,
    "city": "Sao Paulo",
    "state": "SP",
    "avatarUrl": null,
    "bio": null
  }
}
```

#### PATCH `/me/profile`

Request:

```json
{
  "name": "Maria Souza",
  "phone": "+5511999999999",
  "city": "Sao Paulo",
  "state": "SP",
  "bio": "Procuro um filhote calmo."
}
```

Response `200`: perfil atualizado.

### 7.3 Upload De Mídia

#### POST `/media/upload-url`

Request:

```json
{
  "fileName": "puppy-1.jpg",
  "contentType": "image/jpeg",
  "purpose": "PUPPY_PHOTO"
}
```

Response `201`:

```json
{
  "objectKey": "puppies/2026/06/uuid.jpg",
  "uploadUrl": "https://storage.example.com/presigned",
  "publicUrl": "https://cdn.example.com/puppies/2026/06/uuid.jpg",
  "expiresIn": 900
}
```

### 7.4 Anúncios

#### POST `/puppies`

Roles permitidas: `DONOR`, `SHELTER`.

Request:

```json
{
  "name": "Luna",
  "breed": "SRD",
  "species": "DOG",
  "ageMonths": 3,
  "size": "SMALL",
  "sex": "FEMALE",
  "description": "Filhote vacinada e brincalhona.",
  "address": {
    "city": "Sao Paulo",
    "state": "SP",
    "street": "Rua Exemplo",
    "number": "100"
  },
  "photos": [
    {
      "url": "https://cdn.example.com/puppies/2026/06/uuid.jpg",
      "sortOrder": 0,
      "isPrimary": true
    }
  ]
}
```

Response `201`:

```json
{
  "id": "7c0d20e4-d0bf-43ad-ac2b-548e43c8be99",
  "status": "AVAILABLE"
}
```

#### GET `/puppies/{id}`

Response `200`:

```json
{
  "id": "7c0d20e4-d0bf-43ad-ac2b-548e43c8be99",
  "name": "Luna",
  "breed": "SRD",
  "species": "DOG",
  "ageMonths": 3,
  "size": "SMALL",
  "sex": "FEMALE",
  "description": "Filhote vacinada e brincalhona.",
  "status": "AVAILABLE",
  "distanceKm": 2.4,
  "location": {
    "latitude": -23.55052,
    "longitude": -46.63331,
    "city": "Sao Paulo",
    "state": "SP"
  },
  "owner": {
    "id": "b3a2c3d4-80de-4efa-aabc-789d66545bd2",
    "name": "Abrigo Centro"
  },
  "photos": [
    {
      "url": "https://cdn.example.com/puppies/2026/06/uuid.jpg",
      "isPrimary": true,
      "sortOrder": 0
    }
  ]
}
```

#### PATCH `/puppies/{id}`

Permite edição parcial pelo dono do anúncio.

#### PATCH `/puppies/{id}/status`

Request:

```json
{
  "status": "ADOPTED"
}
```

Response `200`:

```json
{
  "id": "7c0d20e4-d0bf-43ad-ac2b-548e43c8be99",
  "status": "ADOPTED",
  "adoptedAt": "2026-06-08T12:30:00Z"
}
```

#### GET `/me/puppies`

Lista anúncios do usuário autenticado.

### 7.5 Busca

#### GET `/puppies/search`

Query params:

```text
lat=-23.55052
lng=-46.63331
radiusKm=25
species=DOG
breed=SRD
ageMinMonths=0
ageMaxMonths=12
size=SMALL
sex=FEMALE
page=0
size=20
sort=DISTANCE_ASC
```

Response `200`:

```json
{
  "items": [
    {
      "id": "7c0d20e4-d0bf-43ad-ac2b-548e43c8be99",
      "name": "Luna",
      "breed": "SRD",
      "species": "DOG",
      "ageMonths": 3,
      "size": "SMALL",
      "sex": "FEMALE",
      "status": "AVAILABLE",
      "distanceKm": 2.4,
      "city": "Sao Paulo",
      "state": "SP",
      "primaryPhotoUrl": "https://cdn.example.com/puppies/2026/06/uuid.jpg"
    }
  ],
  "page": 0,
  "size": 20,
  "totalElements": 1,
  "totalPages": 1
}
```

Cache:

- Chave: hash normalizado de filtros + grid aproximado de latitude/longitude.
- TTL: 60 segundos para buscas sem autenticação de estado do usuário.
- Invalidar por região quando anúncio é criado, removido ou alterado para `ADOPTED`.

### 7.6 Favoritos

```text
POST   /puppies/{id}/favorite
DELETE /puppies/{id}/favorite
GET    /me/favorites?page=0&size=20
```

### 7.7 Buscas Salvas

```text
POST   /me/saved-searches
GET    /me/saved-searches
DELETE /me/saved-searches/{id}
```

Request de criação:

```json
{
  "name": "Filhotes pequenos perto de mim",
  "filters": {
    "lat": -23.55052,
    "lng": -46.63331,
    "radiusKm": 25,
    "species": "DOG",
    "size": "SMALL"
  },
  "notifyOnMatch": true
}
```

### 7.8 Conversas

#### POST `/conversations`

Request:

```json
{
  "puppyId": "7c0d20e4-d0bf-43ad-ac2b-548e43c8be99"
}
```

Response `201`:

```json
{
  "conversationId": "c610f228-8655-448d-a8b7-2e131c7da8a4"
}
```

#### GET `/conversations`

Response `200`:

```json
{
  "items": [
    {
      "id": "c610f228-8655-448d-a8b7-2e131c7da8a4",
      "puppy": {
        "id": "7c0d20e4-d0bf-43ad-ac2b-548e43c8be99",
        "name": "Luna",
        "primaryPhotoUrl": "https://cdn.example.com/puppies/2026/06/uuid.jpg"
      },
      "otherParticipant": {
        "id": "b3a2c3d4-80de-4efa-aabc-789d66545bd2",
        "name": "Abrigo Centro"
      },
      "lastMessage": {
        "content": "Ela ainda esta disponivel?",
        "sentAt": "2026-06-08T12:30:00Z"
      },
      "unreadCount": 1
    }
  ]
}
```

#### GET `/conversations/{id}/messages`

Query params: `before`, `page`, `size`.

#### POST `/conversations/{id}/messages`

Request:

```json
{
  "content": "Ela ainda esta disponivel?"
}
```

Response `201`:

```json
{
  "id": "53bb4e30-7af7-4d29-902c-c9b29ce2cc83",
  "conversationId": "c610f228-8655-448d-a8b7-2e131c7da8a4",
  "senderId": "4a52c5cb-1cd6-49ed-9423-3c918bba6c13",
  "content": "Ela ainda esta disponivel?",
  "sentAt": "2026-06-08T12:30:00Z",
  "readAt": null
}
```

## 8. WebSocket

Endpoint:

```text
WSS /ws/chat?token=<jwt>
```

Eventos cliente -> servidor:

```json
{
  "type": "MESSAGE_SEND",
  "conversationId": "c610f228-8655-448d-a8b7-2e131c7da8a4",
  "clientMessageId": "local-uuid",
  "content": "Oi, tudo bem?"
}
```

```json
{
  "type": "MESSAGE_READ",
  "conversationId": "c610f228-8655-448d-a8b7-2e131c7da8a4",
  "messageId": "53bb4e30-7af7-4d29-902c-c9b29ce2cc83"
}
```

Eventos servidor -> cliente:

```json
{
  "type": "MESSAGE_CREATED",
  "conversationId": "c610f228-8655-448d-a8b7-2e131c7da8a4",
  "message": {
    "id": "53bb4e30-7af7-4d29-902c-c9b29ce2cc83",
    "senderId": "4a52c5cb-1cd6-49ed-9423-3c918bba6c13",
    "content": "Oi, tudo bem?",
    "sentAt": "2026-06-08T12:30:00Z"
  }
}
```

Redis:

- Guardar sessões WebSocket por usuário.
- Publicar eventos entre instâncias do backend quando houver múltiplas réplicas.

## 9. Banco De Dados

Extensões:

```sql
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
```

### 9.1 Tabelas

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(320) NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  role VARCHAR(20) NOT NULL CHECK (role IN ('ADOPTER', 'DONOR', 'SHELTER')),
  status VARCHAR(30) NOT NULL CHECK (status IN ('PENDING_VERIFICATION', 'ACTIVE', 'SUSPENDED', 'DELETED')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE profiles (
  user_id UUID PRIMARY KEY REFERENCES users(id),
  name VARCHAR(120) NOT NULL,
  phone VARCHAR(30),
  city VARCHAR(120) NOT NULL,
  state VARCHAR(2) NOT NULL,
  location GEOGRAPHY(POINT, 4326),
  avatar_url TEXT,
  bio TEXT
);

CREATE TABLE puppies (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  owner_id UUID NOT NULL REFERENCES users(id),
  name VARCHAR(80) NOT NULL,
  breed VARCHAR(120),
  species VARCHAR(20) NOT NULL CHECK (species IN ('DOG', 'CAT', 'OTHER')),
  age_months INT NOT NULL CHECK (age_months >= 0),
  size VARCHAR(20) NOT NULL CHECK (size IN ('SMALL', 'MEDIUM', 'LARGE')),
  sex VARCHAR(20) NOT NULL CHECK (sex IN ('MALE', 'FEMALE', 'UNKNOWN')),
  description TEXT NOT NULL,
  location GEOGRAPHY(POINT, 4326) NOT NULL,
  city VARCHAR(120) NOT NULL,
  state VARCHAR(2) NOT NULL,
  status VARCHAR(20) NOT NULL CHECK (status IN ('AVAILABLE', 'ADOPTED', 'PAUSED', 'REMOVED')),
  adopted_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE puppy_photos (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  puppy_id UUID NOT NULL REFERENCES puppies(id) ON DELETE CASCADE,
  url TEXT NOT NULL,
  sort_order INT NOT NULL,
  is_primary BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (puppy_id, sort_order)
);

CREATE TABLE favorites (
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  puppy_id UUID NOT NULL REFERENCES puppies(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, puppy_id)
);

CREATE TABLE conversations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  puppy_id UUID NOT NULL REFERENCES puppies(id),
  adopter_id UUID NOT NULL REFERENCES users(id),
  donor_id UUID NOT NULL REFERENCES users(id),
  status VARCHAR(20) NOT NULL CHECK (status IN ('OPEN', 'CLOSED', 'BLOCKED')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (puppy_id, adopter_id)
);

CREATE TABLE messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  sender_id UUID NOT NULL REFERENCES users(id),
  content TEXT NOT NULL,
  sent_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  read_at TIMESTAMPTZ
);

CREATE TABLE saved_searches (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR(120) NOT NULL,
  filters JSONB NOT NULL,
  notify_on_match BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### 9.2 Índices

```sql
CREATE INDEX idx_puppies_location ON puppies USING GIST (location);
CREATE INDEX idx_puppies_status_created_at ON puppies (status, created_at DESC);
CREATE INDEX idx_puppies_owner_id ON puppies (owner_id);
CREATE INDEX idx_puppies_filters ON puppies (species, size, sex, age_months);
CREATE INDEX idx_messages_conversation_sent_at ON messages (conversation_id, sent_at DESC);
CREATE INDEX idx_conversations_adopter_id ON conversations (adopter_id);
CREATE INDEX idx_conversations_donor_id ON conversations (donor_id);
CREATE INDEX idx_saved_searches_user_id ON saved_searches (user_id);
CREATE INDEX idx_saved_searches_filters ON saved_searches USING GIN (filters);
```

Restrição para foto primária:

```sql
CREATE UNIQUE INDEX idx_puppy_photos_one_primary
ON puppy_photos (puppy_id)
WHERE is_primary = true;
```

### 9.3 Query PostGIS

```sql
SELECT
  p.*,
  ST_Distance(
    p.location,
    ST_MakePoint(:lng, :lat)::geography
  ) / 1000 AS distance_km
FROM puppies p
WHERE p.status = 'AVAILABLE'
  AND ST_DWithin(
    p.location,
    ST_MakePoint(:lng, :lat)::geography,
    :radius_meters
  )
  AND (:species IS NULL OR p.species = :species)
  AND (:size IS NULL OR p.size = :size)
  AND (:sex IS NULL OR p.sex = :sex)
  AND (:age_min IS NULL OR p.age_months >= :age_min)
  AND (:age_max IS NULL OR p.age_months <= :age_max)
ORDER BY distance_km ASC, p.created_at DESC
LIMIT :limit OFFSET :offset;
```

## 10. Segurança

Autenticação:

- Access token JWT com expiração curta, recomendado 15 minutos.
- Refresh token com expiração maior, recomendado 30 dias.
- Refresh token deve ser revogável no logout e rotação.

Autorização:

- `ADOPTER`: buscar, favoritar, iniciar conversas, enviar mensagens.
- `DONOR`: tudo de `ADOPTER` mais criar e gerir próprios anúncios.
- `SHELTER`: igual `DONOR`, reservado para abrigos verificados.

Proteções:

- HTTPS obrigatório.
- WSS obrigatório.
- Rate limit em login, cadastro, envio de código e busca.
- Hash de senha com BCrypt ou Argon2id.
- Nunca retornar `password_hash`.
- Sanitização de texto livre para evitar XSS em clientes.
- Logs sem tokens, senhas ou dados sensíveis.
- LGPD: permitir exclusão/desativação de conta e remoção/anomização de dados pessoais.

## 11. Erros

Formato padrão:

```json
{
  "code": "PUPPY_NOT_FOUND",
  "message": "Anuncio nao encontrado.",
  "details": {
    "puppyId": "7c0d20e4-d0bf-43ad-ac2b-548e43c8be99"
  },
  "traceId": "01JY0000000000000000000000"
}
```

Status:

```text
400 VALIDATION_ERROR
401 UNAUTHORIZED
403 FORBIDDEN
404 NOT_FOUND
409 CONFLICT
422 BUSINESS_RULE_VIOLATION
429 RATE_LIMITED
500 INTERNAL_ERROR
```

Erros de domínio:

```text
EMAIL_ALREADY_REGISTERED
ACCOUNT_NOT_VERIFIED
INVALID_CREDENTIALS
USER_NOT_ALLOWED_TO_CREATE_LISTING
PUPPY_NOT_FOUND
PUPPY_NOT_AVAILABLE
ONLY_OWNER_CAN_EDIT_LISTING
CONVERSATION_NOT_FOUND
USER_NOT_IN_CONVERSATION
MESSAGE_TOO_LONG
UPLOAD_URL_EXPIRED
GEOLOCATION_FAILED
```

## 12. Observabilidade

Logs:

- `traceId`, `userId` quando autenticado, rota, status HTTP, duração.
- Eventos de domínio relevantes: criação de anúncio, adoção, início de conversa.

Métricas:

- Latência p50/p95/p99 por rota.
- Taxa de erro por rota.
- Quantidade de buscas por minuto.
- Cache hit ratio para buscas.
- Mensagens enviadas por minuto.
- Falhas de envio push/e-mail.

Tracing:

- Propagar `traceId` em adapters externos: S3, geocoding, push, e-mail/SMS.

## 13. Jobs

### 13.1 ExpireListingsJob

Objetivo:

- Pausar ou sinalizar anúncios muito antigos para revisão do dono.

Periodicidade:

- Diário.

Regra inicial:

- Anúncios `AVAILABLE` sem atualização há 90 dias geram notificação ao dono.
- Após 14 dias sem resposta, mover para `PAUSED`.

### 13.2 SavedSearchNotificationJob

Objetivo:

- Notificar usuários quando um novo anúncio corresponde a uma busca salva.

Periodicidade:

- A cada 15 minutos ou por evento de novo anúncio.

Regra:

- Não enviar mais de uma notificação por busca salva em janela de 24 horas.

## 14. Fluxos Detalhados

### 14.1 Busca De Filhotes Próximos

```text
1. App obtem GPS do dispositivo.
2. App chama GET /api/v1/puppies/search.
3. API valida token quando existir, filtros e paginacao.
4. `internal/app` monta `PuppySearchQuery`.
5. CachePort tenta recuperar resultado.
6. SearchRepository executa query PostGIS quando cache miss.
7. `internal/app` salva cache por TTL curto.
8. API retorna lista paginada com distancia e foto primaria.
```

### 14.2 Cadastro De Filhote

```text
1. Doador solicita URLs pre-assinadas em POST /api/v1/media/upload-url.
2. App envia fotos diretamente ao object storage.
3. App envia POST /api/v1/puppies com dados e URLs publicas.
4. API valida role DONOR/SHELTER.
5. `internal/app` chama `GeoLocationPort` quando endereco foi informado.
6. `internal/domain` valida regras de anuncio e fotos.
7. PuppyRepository persiste anuncio, fotos e localizacao.
8. Cache de buscas da regiao e invalidado.
9. Anuncio fica visivel para busca.
```

### 14.3 Chat In-App

```text
1. Adotante toca em contato no perfil do filhote.
2. App chama POST /api/v1/conversations.
3. `internal/app` cria ou retorna conversa existente.
4. App abre WSS /ws/chat?token=<jwt>.
5. Cliente envia MESSAGE_SEND.
6. SendMessageUseCase valida participante, persiste mensagem e publica evento realtime.
7. RealtimePort entrega MESSAGE_CREATED aos participantes conectados.
8. NotificationPort envia push ao participante offline.
```

## 15. Configuração

Variáveis:

```text
APP_ENV=local|dev|staging|prod
DATABASE_URL=postgres://...
DATABASE_USERNAME=...
DATABASE_PASSWORD=...
REDIS_URL=redis://...
JWT_ISSUER=adotapet
JWT_ACCESS_SECRET=...
JWT_REFRESH_SECRET=...
S3_ENDPOINT=...
S3_BUCKET=...
S3_ACCESS_KEY=...
S3_SECRET_KEY=...
CDN_BASE_URL=...
FCM_CREDENTIALS_JSON=...
APNS_KEY_ID=...
APNS_TEAM_ID=...
GEOCODING_API_KEY=...
EMAIL_PROVIDER_API_KEY=...
```

## 16. Testes

Unitários:

- Entidades e value objects de domínio.
- Regras de criação, edição, adoção e conversa.
- Use cases com ports mockados.

Integração:

- Repositórios Go PostgreSQL com PostgreSQL + PostGIS via Testcontainers.
- Query de busca geográfica com raio e ordenação por distância.
- Flyway migrations.
- WebSocket auth e entrega de eventos.

Contrato:

- OpenAPI para REST.
- Contratos de eventos WebSocket.

End-to-end:

- Cadastro/login.
- Criar anúncio com foto.
- Buscar por proximidade.
- Iniciar conversa e enviar mensagem.

## 17. Critérios De Aceite Do MVP

- Usuário consegue se cadastrar, verificar conta e autenticar.
- Doador consegue cadastrar filhote com foto e localização.
- Adotante consegue buscar filhotes próximos com filtros.
- Resultado de busca retorna distância aproximada e foto principal.
- Adotante consegue abrir perfil do filhote e iniciar conversa.
- Mensagens são persistidas e entregues em tempo real quando ambos estão conectados.
- Push é enviado quando o destinatário está offline.
- Dono consegue marcar anúncio como adotado.
- Busca p95 abaixo de 500 ms em dataset representativo do MVP.

## 18. Próximos Artefatos

- `openapi.yaml` com contratos REST.
- `V1__init_schema.sql` com migrations Flyway.
- Diagramas de sequência para busca, cadastro de anúncio e chat.
- Scaffolding Go do backend.
- Projeto Flutter com navegação base e clients REST/WebSocket.

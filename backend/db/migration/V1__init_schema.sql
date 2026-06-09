CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

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
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR(120) NOT NULL,
  phone VARCHAR(30),
  city VARCHAR(120) NOT NULL,
  state VARCHAR(2) NOT NULL,
  location GEOGRAPHY(POINT, 4326),
  avatar_url TEXT,
  bio TEXT
);

CREATE TABLE account_verification_codes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  code_hash TEXT NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  consumed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE refresh_tokens (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash TEXT NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
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
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK ((status = 'ADOPTED' AND adopted_at IS NOT NULL) OR (status <> 'ADOPTED' AND adopted_at IS NULL))
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

CREATE TABLE saved_searches (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR(120) NOT NULL,
  filters JSONB NOT NULL,
  notify_on_match BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE conversations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  puppy_id UUID NOT NULL REFERENCES puppies(id),
  adopter_id UUID NOT NULL REFERENCES users(id),
  donor_id UUID NOT NULL REFERENCES users(id),
  status VARCHAR(20) NOT NULL CHECK (status IN ('OPEN', 'CLOSED', 'BLOCKED')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (puppy_id, adopter_id),
  CHECK (adopter_id <> donor_id)
);

CREATE TABLE messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  sender_id UUID NOT NULL REFERENCES users(id),
  content TEXT NOT NULL CHECK (char_length(content) BETWEEN 1 AND 2000),
  sent_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  read_at TIMESTAMPTZ
);

CREATE INDEX idx_profiles_location ON profiles USING GIST (location);
CREATE INDEX idx_account_verification_codes_user_id ON account_verification_codes (user_id);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);
CREATE INDEX idx_puppies_location ON puppies USING GIST (location);
CREATE INDEX idx_puppies_status_created_at ON puppies (status, created_at DESC);
CREATE INDEX idx_puppies_owner_id ON puppies (owner_id);
CREATE INDEX idx_puppies_filters ON puppies (species, size, sex, age_months);
CREATE INDEX idx_messages_conversation_sent_at ON messages (conversation_id, sent_at DESC);
CREATE INDEX idx_conversations_adopter_id ON conversations (adopter_id);
CREATE INDEX idx_conversations_donor_id ON conversations (donor_id);
CREATE INDEX idx_saved_searches_user_id ON saved_searches (user_id);
CREATE INDEX idx_saved_searches_filters ON saved_searches USING GIN (filters);

CREATE UNIQUE INDEX idx_puppy_photos_one_primary
ON puppy_photos (puppy_id)
WHERE is_primary = true;

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_puppies_set_updated_at
BEFORE UPDATE ON puppies
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_conversations_set_updated_at
BEFORE UPDATE ON conversations
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

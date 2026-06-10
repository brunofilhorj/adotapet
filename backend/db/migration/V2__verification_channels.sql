ALTER TABLE account_verification_codes
  ADD COLUMN channel VARCHAR(20),
  ADD COLUMN destination VARCHAR(320);

UPDATE account_verification_codes c
SET
  channel = 'EMAIL',
  destination = u.email
FROM users u
WHERE u.id = c.user_id;

ALTER TABLE account_verification_codes
  ALTER COLUMN channel SET NOT NULL,
  ALTER COLUMN destination SET NOT NULL,
  ADD CONSTRAINT chk_account_verification_channel
    CHECK (channel IN ('EMAIL', 'SMS', 'WHATSAPP', 'PUSH'));

CREATE INDEX idx_account_verification_codes_lookup
ON account_verification_codes (user_id, channel, destination, code_hash, expires_at DESC)
WHERE consumed_at IS NULL;

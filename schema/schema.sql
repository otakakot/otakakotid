-- account

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

CREATE TABLE IF NOT EXISTS webauthn_credentials (
    raw_id_base64 TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    credential_base64 TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_webauthn_credentials_user_id ON webauthn_credentials(user_id);

-- oidc

CREATE TABLE IF NOT EXISTS clients (
    id TEXT PRIMARY KEY,
    der_public_key_base64 TEXT NOT NULL,
    name TEXT NOT NULL,
    redirect_uri TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS jwk_sets (
    id TEXT PRIMARY KEY,
    der_private_key_base64 TEXT NOT NULL
);

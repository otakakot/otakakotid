-- name: InsertUser :one
INSERT INTO
    users (id, email)
VALUES
    (?, ?) RETURNING *;

-- name: InsertWebAuthnCredential :one
INSERT INTO
    webauthn_credentials (raw_id_base64, user_id, credential_base64)
VALUES
    (?, ?, ?) RETURNING *;

-- name: FindWebAuthnCredentialByRawID :one
SELECT
    *
FROM
    webauthn_credentials
WHERE
    raw_id_base64 = ?;

package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/otakakot/otakakotid/internal/core"
	"github.com/otakakot/otakakotid/internal/database"
	"github.com/otakakot/otakakotid/pkg/schema"
	"github.com/syumai/workers/cloudflare"
)

func FinalizeAttestation(rw http.ResponseWriter, req *http.Request) {
	wa, err := core.NewWebAuthn(req)
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	sid, err := req.Cookie(core.CookeyAttestation)
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	code, err := req.Cookie(core.CookeyRegistration)
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	sessionKV, err := cloudflare.NewKVNamespace(database.KVNSWebAuthnSessionData)
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	sessionBase64, err := sessionKV.GetString(sid.Value, nil)
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	sessionByte, _ := base64.StdEncoding.DecodeString(sessionBase64)

	session := webauthn.SessionData{}

	if err := json.NewDecoder(bytes.NewReader(sessionByte)).Decode(&session); err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := sessionKV.Delete(sid.Value); err != nil {
		slog.InfoContext(req.Context(), "failed to delete session data error: "+err.Error())
	}

	emailKV, err := cloudflare.NewKVNamespace(database.KVNSRegistrationEmail)
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	email, err := emailKV.GetString(code.Value, nil)
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := emailKV.Delete(code.Value); err != nil {
		slog.WarnContext(req.Context(), err.Error())
	}

	user := &core.User{
		ID: uuid.MustParse(string(session.UserID)),
	}

	credential, err := wa.FinishRegistration(user, session, req)
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	buf := bytes.Buffer{}

	if err := json.NewEncoder(&buf).Encode(credential); err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	// NOTE: d1 Transaction Not Supported

	if _, err := schema.New(database.D1).InsertUser(req.Context(), schema.InsertUserParams{
		ID:    user.ID.String(),
		Email: email,
	}); err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	if _, err := schema.New(database.D1).InsertWebAuthnCredential(req.Context(), schema.InsertWebAuthnCredentialParams{
		RawIDBase64:      base64.StdEncoding.EncodeToString(credential.ID),
		UserID:           user.ID.String(),
		CredentialBase64: base64.StdEncoding.EncodeToString(buf.Bytes()),
	}); err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:   core.CookeyAttestation,
		Value:  "",
		MaxAge: -1,
	})

	http.SetCookie(rw, &http.Cookie{
		Name:   core.CookeyRegistration,
		Value:  "",
		MaxAge: -1,
	})

	rw.WriteHeader(http.StatusCreated)
}

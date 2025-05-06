package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/otakakot/otakakotid/internal/core"
	"github.com/otakakot/otakakotid/internal/database"
	"github.com/otakakot/otakakotid/pkg/schema"
	"github.com/syumai/workers/cloudflare"
)

func FinalizeAssertion(rw http.ResponseWriter, req *http.Request) {
	wa, err := core.NewWebAuthn(req)
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	sid, err := req.Cookie(core.CookeyAssertion)
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

	slog.InfoContext(req.Context(), fmt.Sprintf("session: %+v", session))

	if err := sessionKV.Delete(sid.Value); err != nil {
		slog.InfoContext(req.Context(), "failed to delete session data error: "+err.Error())
	}

	userHandler := func(rawID []byte, userHandle []byte) (webauthn.User, error) {
		uid, err := uuid.Parse(string(userHandle))
		if err != nil {
			return nil, err
		}

		rawIDBase64 := base64.StdEncoding.EncodeToString(rawID)

		cre, err := schema.New(database.D1).FindWebAuthnCredentialByRawID(req.Context(), rawIDBase64)
		if err != nil {
			return nil, err
		}

		creByte, err := base64.StdEncoding.DecodeString(cre.CredentialBase64)
		if err != nil {
			return nil, err
		}

		credential := webauthn.Credential{}

		if err := json.NewDecoder(bytes.NewBuffer(creByte)).Decode(&credential); err != nil {
			return nil, err
		}

		return &core.User{
			ID:          uid,
			Credentials: []webauthn.Credential{credential},
		}, nil
	}

	if _, err := wa.FinishDiscoverableLogin(
		func(rawID, userHandle []byte) (user webauthn.User, err error) {
			return userHandler(rawID, userHandle)
		},
		session,
		req,
	); err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:   core.CookeyAttestation,
		Value:  "",
		MaxAge: -1,
	})

	// TODO:
	// もし OIDC Authorization Code Flow であれば、ここでリダイレクトする

	rw.WriteHeader(http.StatusCreated)
}

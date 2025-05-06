package handler

import (
	"bytes"
	"net/http"
	"net/url"

	"github.com/otakakot/otakakotid/internal/core"
)

func FinalizeRegistration(rw http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")

	http.SetCookie(rw, &http.Cookie{
		Name:     core.CookeyRegistration,
		Value:    code,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	redirectBuf := bytes.Buffer{}

	// TODO: use env

	redirectBuf.WriteString("http://localhost:3000/link")

	redirectURL, _ := url.ParseRequestURI(redirectBuf.String())

	http.Redirect(rw, req, redirectURL.String(), http.StatusFound)
}

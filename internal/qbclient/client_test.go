package qbclient

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"
)

type stubRoundTripper struct {
	t        *testing.T
	handlers []func(*http.Request) *http.Response
	idx      int
}

func (s *stubRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if s.idx >= len(s.handlers) {
		s.t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
	}
	h := s.handlers[s.idx]
	s.idx++
	return h(r), nil
}

func TestLoginAndAddMagnet(t *testing.T) {
	jar, _ := cookiejar.New(nil)

	rt := &stubRoundTripper{
		t: t,
		handlers: []func(*http.Request) *http.Response{
			func(r *http.Request) *http.Response {
				body, _ := io.ReadAll(r.Body)
				if r.Method != http.MethodPost || r.URL.Path != "/api/v2/auth/login" {
					t.Fatalf("unexpected login request: %s %s", r.Method, r.URL.Path)
				}
				if !strings.Contains(string(body), "username=admin") || !strings.Contains(string(body), "password=password") {
					t.Fatalf("login body missing credentials: %s", string(body))
				}
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{},
					Body:       io.NopCloser(strings.NewReader("Ok.")),
				}
				resp.Header.Set("Set-Cookie", "SID=abc; Path=/")
				resp.Request = r
				return resp
			},
			func(r *http.Request) *http.Response {
				if r.Method != http.MethodPost || r.URL.Path != "/api/v2/torrents/add" {
					t.Fatalf("unexpected add request: %s %s", r.Method, r.URL.Path)
				}
				if cookie := r.Header.Get("Cookie"); !strings.Contains(cookie, "SID=abc") {
					t.Fatalf("expected SID cookie, got %s", cookie)
				}
				ct := r.Header.Get("Content-Type")
				mediaType, params, err := mime.ParseMediaType(ct)
				if err != nil {
					t.Fatalf("ParseMediaType: %v", err)
				}
				if mediaType != "multipart/form-data" {
					t.Fatalf("unexpected media type: %s", mediaType)
				}
				reader := multipart.NewReader(r.Body, params["boundary"])
				part, err := reader.NextPart()
				if err != nil {
					t.Fatalf("reading part: %v", err)
				}
				if part.FormName() != "urls" {
					t.Fatalf("unexpected form field: %s", part.FormName())
				}
				body, _ := io.ReadAll(part)
				if string(body) != "magnet:?xt=urn:btih:example" {
					t.Fatalf("unexpected magnet value: %s", string(body))
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{},
					Body:       io.NopCloser(strings.NewReader("Ok.")),
					Request:    r,
				}
			},
		},
	}

	client := &http.Client{Transport: rt, Jar: jar}
	qb := NewWithClient("http://example.test", "admin", "password", client)

	if err := qb.Login(); err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if err := qb.AddMagnet("magnet:?xt=urn:btih:example"); err != nil {
		t.Fatalf("AddMagnet() error = %v", err)
	}
}

func TestAddMagnetErrorStatus(t *testing.T) {
	rt := &stubRoundTripper{
		t: t,
		handlers: []func(*http.Request) *http.Response{
			func(r *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader("nope")),
					Header:     http.Header{},
					Request:    r,
				}
			},
		},
	}
	client := &http.Client{Transport: rt}
	qb := NewWithClient("http://example.test", "admin", "password", client)

	err := qb.AddMagnet("magnet:?xt=urn:btih:example")
	if err == nil || !strings.Contains(err.Error(), "nope") {
		t.Fatalf("expected error containing nope, got %v", err)
	}
}

func TestLoginFailureStatus(t *testing.T) {
	rt := &stubRoundTripper{
		t: t,
		handlers: []func(*http.Request) *http.Response{
			func(r *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(bytes.NewReader(nil)),
					Header:     http.Header{},
					Request:    r,
				}
			},
		},
	}
	client := &http.Client{Transport: rt}
	qb := NewWithClient("http://example.test", "admin", "password", client)

	err := qb.Login()
	if err == nil || !strings.Contains(err.Error(), "status 403") {
		t.Fatalf("expected login error containing status, got %v", err)
	}
}

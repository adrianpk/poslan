// Package mailer allows to maintain a redundant mail delivery service.
package mailer

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	c "github.com/adrianpk/poslan/internal/config"
	httptransport "github.com/go-kit/kit/transport/http"
)

const (
	loggingOn         = true
	tracingOn         = true
	instrumentationOn = true
	transactionOn     = true
)

var (
	authTokenCtxKey = contextKey("auth-token")
)

// Run the mailer service.
func Run() {

	// Context
	ctx, cancel := context.WithCancel(context.Background())
	go checkSigTerm(cancel)

	// Logger
	logger := makeLogger()

	// Config
	cfg, err := c.Load(logger)
	checkError(err)

	// Service
	svc, err := makeService(ctx, cfg, logger).Init()
	checkError(err)

	// Handlers
	http.Handle("/signin", SignInHandler(svc))
	http.Handle("/signout", SignOutHandler(svc))
	http.Handle("/send", SendHandler(svc))

	err = http.ListenAndServe(cfg.App.ServerPortFmt(), nil)

	logger.Log("level", c.LogLevel.Error, "msg", err.Error())
}

// SignInHandler manages signin up process.
func SignInHandler(svc Service) *httptransport.Server {
	return httptransport.NewServer(
		makeSignInEndpoint(svc),
		decodeSignInRequest,
		encodeResponse,
	)
}

// SignOutHandler manages signout up process.
func SignOutHandler(svc Service) *httptransport.Server {
	opts := httptransport.ServerBefore(tokenToContext)
	return httptransport.NewServer(
		makeSignOutEndpoint(svc),
		decodeSignOutRequest,
		encodeResponse,
		opts,
	)
}

// SendHandler manages email sending.
func SendHandler(svc Service) *httptransport.Server {
	opts := httptransport.ServerBefore(tokenToContext)
	return httptransport.NewServer(
		makeSendEndpoint(svc),
		decodeSendRequest,
		encodeResponse,
		opts,
	)
}

// Decoders
func decodeSignInRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request signInRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeSignOutRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request signOutRequest
	tokenToContext(ctx, r)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func decodeSendRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request sendRequest
	tokenToContext(ctx, r)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

// Encoders
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func readToken(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	splitToken := strings.Split(token, "Bearer")
	if len(splitToken) != 2 {
		return "", errors.New("invalid token format")
	}
	return strings.TrimSpace(splitToken[1]), nil
}

func tokenToContext(ctx context.Context, r *http.Request) context.Context {
	token := r.Header.Get("Authorization")
	token, err := readToken(r)
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, authTokenCtxKey, token)
}

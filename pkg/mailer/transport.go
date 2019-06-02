// Package mailer allows to maintain a redundant mail delivery service.
package mailer

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	c "github.com/adrianpk/poslan/internal/config"
	"github.com/adrianpk/poslan/pkg/auth"
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
	userDataCtxKey  = contextKey("user-data")
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
	opts := httptransport.ServerBefore(userDataToContext)
	return httptransport.NewServer(
		makeSignOutEndpoint(svc),
		decodeSignOutRequest,
		encodeResponse,
		opts,
	)
}

// SendHandler manages email sending.
func SendHandler(svc Service) *httptransport.Server {
	opts := httptransport.ServerBefore(userDataToContext)
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
	userDataToContext(ctx, r)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func decodeSendRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request sendRequest
	userDataToContext(ctx, r)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

// Encoders
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// userDataToContext read bearer token from request header.
func readToken(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	splitToken := strings.Split(token, "Bearer")
	if len(splitToken) != 2 {
		return "", errors.New("invalid token format")
	}
	return strings.TrimSpace(splitToken[1]), nil
}

// userDataToContext extracts bearer token from request header and stores it
// in context.
func userDataToContext(ctx context.Context, r *http.Request) context.Context {
	tk := r.Header.Get("Authorization")
	tk, err := readToken(r)
	if err != nil {
		return ctx
	}

	ud, err := auth.UserData(tk)
	if err != nil {
		return ctx
	}

	ctx = context.WithValue(ctx, userDataCtxKey, ud)
	ctx = context.WithValue(ctx, authTokenCtxKey, tk)

	return ctx
}

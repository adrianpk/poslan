// Package mailer allows to maintain a redundant mail delivery service.
package mailer

import (
	"context"
	"encoding/json"
	"net/http"

	c "github.com/adrianpk/poslan/internal/config"
	httptransport "github.com/go-kit/kit/transport/http"
)

const (
	loggingOn         = true
	tracingOn         = true
	instrumentationOn = true
	transactionOn     = true
)

// Run the mailer service.
func Run() {

	// Context
	ctx, cancel := context.WithCancel(context.Background())
	go checkSigTerm(cancel)

	// Config
	cfg, err := c.Load()
	checkError(err)

	// Logger
	logger := makeLogger()

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
	return httptransport.NewServer(
		makeSignOutEndpoint(svc),
		decodeSignOutRequest,
		encodeResponse,
	)
}

// SendHandler manages email sending.
func SendHandler(svc Service) *httptransport.Server {
	return httptransport.NewServer(
		makeSendEndpoint(svc),
		decodeSendRequest,
		encodeResponse,
	)
}

// Decoders
func decodeSignInRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request signInRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeSignOutRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request signOutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeSendRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request sendRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// Encoders
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

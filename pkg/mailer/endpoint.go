/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package mailer

import (
	"context"
	"fmt"

	c "github.com/adrianpk/poslan/internal/config"
	"github.com/go-kit/kit/endpoint"
)

func makeSignInEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(signInRequest)

		reqstr := fmt.Sprintf("Req: %+v", req)
		svc.Logger().Log("level", c.LogLevel.Info, "req", reqstr)

		user, err := svc.SignIn(req.Username, req.Password)
		if err != nil {
			return signInResponse{user, err.Error()}, nil
		}

		return signInResponse{user, ""}, nil
	}
}

func makeSignOutEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(signOutRequest)

		reqstr := fmt.Sprintf("Req: %+v", req)
		svc.Logger().Log("level", c.LogLevel.Info, "req", reqstr)

		err := svc.SignOut(req.ID)
		if err != nil {
			return signOutResponse{err.Error()}, nil
		}

		return signOutResponse{""}, nil
	}
}

func makeSendEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(sendRequest)

		reqstr := fmt.Sprintf("Req: %+v", req)
		svc.Logger().Log("level", c.LogLevel.Info, "req", reqstr)

		err := svc.Send(req.To, req.Cc, req.Bcc, req.Subject, req.Body)
		if err != nil {
			return signOutResponse{err.Error()}, nil
		}

		return sendResponse{Err: ""}, nil
	}
}

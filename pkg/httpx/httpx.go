// Package httpx provides you with extended http handler
package httpx

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/abdivasiyev/rester/pkg/encoder"
	"github.com/abdivasiyev/rester/pkg/errorsx"
	"github.com/abdivasiyev/rester/pkg/slogx"
)

// A Validatable interface to implement validation function individually for every request type using Validate function
type Validatable interface {
	Validate() error
}

// A Bindable interface to implement bindings between [http.Request] and your custom Request structure
type Bindable interface {
	Bind(*http.Request) error
}

// A Request is a type parameter to pass user provided request to Handle method
// In this type parameter following constraints are added: Bindable, Validatable and [fmt.Stringer]
type Request[Req any] interface {
	Bindable
	Validatable
	fmt.Stringer
	*Req
}

// A DefaultRequest is an implementation of empty request. If you don't want to implement all the methods of Request
// constraint, so you can embed DefaultRequest to your request structure
type DefaultRequest struct{}

func (*DefaultRequest) Bind(*http.Request) error {
	return nil
}

func (*DefaultRequest) Validate() error {
	return nil
}

type DefaultResponse struct {
	Message string `json:"message"`
}

// UseCaseFunc is a type to implement business logic functions
type UseCaseFunc[Req any, Resp any] func(context.Context, Req) (Resp, error)

type handlerOptions struct {
	successCode int
	encoder     encoder.Encoder
	logger      *slog.Logger
}

// An Option is a type to set optional parameters to handler
type Option func(h *handlerOptions)

// WithSuccessCode sets success code to handler. Default value is a [http.StatusOK]
func WithSuccessCode(code int) Option {
	return func(h *handlerOptions) {
		h.successCode = code
	}
}

// WithEncoder sets custom encoder to handler. Default value is a [encoder.JsonEncoder]
func WithEncoder(encoder encoder.Encoder) Option {
	return func(h *handlerOptions) {
		h.encoder = encoder
	}
}

// WithLogger sets custom slog instance to handler. Default value is generated from slogx.New()
func WithLogger(logger *slog.Logger) Option {
	return func(h *handlerOptions) {
		h.logger = logger
	}
}

func applyOptions(options ...Option) handlerOptions {
	var h handlerOptions

	for _, option := range options {
		option(&h)
	}

	if h.successCode <= 0 {
		h.successCode = http.StatusOK
	}

	if h.encoder == nil {
		h.encoder = encoder.JsonEncoder
	}

	if h.logger == nil {
		h.logger = slogx.New()
	}

	return h
}

// Handle receives request and response structs as type parameters to pass to use case function.
// Using options you can add your custom response codes and encoders to handler.
//
// Returns [http.HandlerFunc] to pass to your multiplexer
//
// Usage:
//
//	mux.HandleFunc("GET /", httpx.Handle[Request, Response](handleIndex, httpx.WithSuccessCode(http.StatusBadRequest), httpx.WithLogger(logger)))
func Handle[Req any, Resp any, _Req Request[Req]](useCase UseCaseFunc[Req, Resp], options ...Option) http.HandlerFunc {
	var h = applyOptions(options...)

	return func(w http.ResponseWriter, r *http.Request) {
		var (
			id   = uuid.New().String()
			req  Req
			_req = _Req(&req)
			err  error
		)

		err = _req.Bind(r)
		if err != nil {
			if errx, ok := errorsx.As(err); ok && !errx.Internal() {
				h.logger.WithGroup(id).Error("failed to bind request", slog.Any("err", errx))
				w.WriteHeader(errx.Code())
				err = h.encoder.New(w).Encode(DefaultResponse{Message: err.Error()})
				if err != nil {
					h.logger.WithGroup(id).Error("failed to write error response", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			err = h.encoder.New(w).Encode(DefaultResponse{Message: http.StatusText(http.StatusInternalServerError)})
			if err != nil {
				h.logger.WithGroup(id).Error("failed to write error response", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		h.logger.WithGroup(id).Info("request", _req)

		err = _req.Validate()
		if err != nil {
			if errx, ok := errorsx.As(err); ok && !errx.Internal() {
				h.logger.WithGroup(id).Error("failed to validate request", slog.Any("err", errx))
				w.WriteHeader(errx.Code())
				err = h.encoder.New(w).Encode(errx.Error())
				if err != nil {
					h.logger.WithGroup(id).Error("failed to write error response", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		response, err := useCase(r.Context(), req)
		if err != nil {
			if errx, ok := errorsx.As(err); ok && !errx.Internal() {
				w.WriteHeader(errx.Code())
				err = h.encoder.New(w).Encode(errx.Error())
				if err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		h.logger.WithGroup(id).Info("response", response)

		w.WriteHeader(h.successCode)
		err = h.encoder.New(w).Encode(response)
		if err != nil {
			if errx, ok := errorsx.As(err); ok && !errx.Internal() {
				http.Error(w, errx.Error(), errx.Code())
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

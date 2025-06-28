// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/alexliesenfeld/health"
	"go.brokedaear.com/internal/common/telemetry"
	"go.brokedaear.com/internal/common/utils/loggers"
	"go.brokedaear.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// HTTPServer represents an HTTP server that is capable of accepting routes
// and connections from clients via HTTP. Implementations for this interface
// include endpoints that serve GraphQL or simple barebones requests.
//
// This interface also implements the io.Closer interface, for use in global
// teardown operations.
type HTTPServer interface {
	ListenAndServe(context.Context) error
	// RegisterRoutes(...HTTPRoute)
	io.Closer
}

type httpServer struct {
	*Base
	srv *http.Server
}

const httpHealthTimeout = 10 * time.Second

// NewHTTPServer creates a new HTTP server using a logger and a config.
// The server comes with telemetry enabled by default.
func NewHTTPServer(logger Logger, configOpts ...ConfigOpts) (HTTPServer, error) {
	const (
		readTimeout  = 10 * time.Second
		writeTimeout = 30 * time.Second
	)

	b, err := NewBase(logger, configOpts...)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	checker := health.NewChecker(
		health.WithCacheDuration(1*time.Second),
		health.WithTimeout(httpHealthTimeout),
	)

	mux.Handle("/health", health.NewHandler(checker))

	return &httpServer{
		Base: b,
		srv: &http.Server{
			IdleTimeout:  time.Minute,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			Handler:      mux,
		},
	}, nil
}

// ListenAndServe listens to specified route endpoints given by route functions
// specified by the function signature. The server is terminated via error or
// the server's interface io.Closer method. An error is only returned when the
// closure results from an error.
func (s httpServer) ListenAndServe(ctx context.Context) error {
	var serverError error

	serverCtx, serverCancel := context.WithCancel(ctx)

	defer serverCancel()

	go func() {
		defer serverCancel()
		err := s.srv.Serve(s.listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error(err.Error())
			serverError = err
		}
	}()

	<-serverCtx.Done()

	return serverError
}

func (s httpServer) Close() error {
	const shutdownTimeout = 20 * time.Second
	shutdownCtx, shutdownCancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(shutdownTimeout),
	)

	defer shutdownCancel()

	err := s.srv.Shutdown(shutdownCtx)
	if err != nil {
		s.logger.Warn("failed to shutdown http server, killing", "err", err)
		err = s.srv.Close()
		if err != nil {
			return err
		}
	}

	err = s.listener.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close http listener")
	}

	s.logger.Info("http server closed")

	return nil
}

type HTTPRoute interface {
	String() string
	Route() http.HandlerFunc
}

func (s httpServer) RegisterRoutes(routes ...HTTPRoute) {
	s.srv.Handler = s.registerRoutes(routes...)
}

// registerRoutes needs to be REFACTORED.
// TODO: Refactor registerRoutes
func (s httpServer) registerRoutes(routes ...HTTPRoute) http.Handler {
	mux := http.NewServeMux()

	handleFunc := func(pattern string, handlerFunc http.HandlerFunc) {
		mux.HandleFunc(pattern, handlerFunc)
	}

	if s.config.telemetry != nil {
		handleFunc = func(pattern string, handlerFunc http.HandlerFunc) {
			h := otelhttp.WithRouteTag(pattern, handlerFunc)
			mux.Handle(pattern, h)
		}
	}

	for _, v := range routes {
		handleFunc(v.String(), v.Route())
	}

	if s.config.telemetry != nil {
		return otelhttp.NewHandler(mux, "/")
	}

	return mux
}

// BodyParser parses a body returned from an HTTP request and returns
// a specified type.
type BodyParser[T any] interface {
	// Parse parses a body returned from an HTTP request and translates it into
	// a type T.
	Parse(body []byte) (T, error)
}

// HTTPClient is an HTTPClient capable of making HTTP requests.
type HTTPClient[T any] interface {
	// Post()
	// Get makes a GET request to a URL and returns a parsed response of type T.
	Get(ctx context.Context, url string, headers map[string]string) (T, error)
	// Post makes a POST request to a given URL. It returns a a function that
	// executes the request.
}

type httpClient[T any] struct {
	client *http.Client
	logger loggers.Logger
	tel    telemetry.Telemetry
	parser BodyParser[T]
}

func NewHTTPRequestClient[T any](
	logger loggers.Logger,
	tel telemetry.Telemetry,
	parser BodyParser[T],
) HTTPClient[T] {
	const httpClientTimeout = time.Second * 30
	hc := &http.Client{
		Transport: nil,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			panic("TODO")
		},
		Jar:     nil,
		Timeout: httpClientTimeout,
	}
	return &httpClient[T]{
		client: hc,
		logger: logger,
		tel:    tel,
		parser: parser,
	}
}

func (h httpClient[T]) Get(ctx context.Context, url string, headers map[string]string) (T, error) {
	var (
		v   T
		b   []byte
		err error
		req *http.Request
		res *http.Response
	)

	req, err = http.NewRequestWithContext(
		ctx, http.MethodGet, url,
		bytes.NewReader(b),
	)
	if err != nil {
		return v, errors.Wrap(err, "failed to create request")
	}

	// Ensure we close the connection of the request

	req.Close = true

	// Add headers to the request.

	if headers != nil {
		h.mapHeaders(req, headers)
	}

	res, err = h.client.Do(req)
	if res != nil { // See: https://golang50shades.com/index.html#close_http_resp_body
		defer res.Body.Close()
	}
	if err != nil {
		return v, errors.Wrap(err, "failed to do request")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return v, errors.Wrap(err, "failed to read response body")
	}

	return h.parser.Parse(resBody)
}

// mapHeaders maps a KV pair of strings to a request.
func (h httpClient[T]) mapHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Add(k, v)
	}
}

func (h httpClient[T]) Post(url *url.URL, body any) (func(context.Context) (T, error), error) {
	var err error

	b, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare http request body")
	}

	return func(ctx context.Context) (T, error) {
		var (
			v       T
			req     *http.Request
			res     *http.Response
			resBody []byte
		)

		req, err = http.NewRequestWithContext(
			ctx, http.MethodPost, url.String(),
			bytes.NewReader(b),
		)
		if err != nil {
			return v, errors.Wrap(err, "failed to create request")
		}

		req.Close = true

		res, err = h.client.Do(req)
		if res != nil { // See: https://golang50shades.com/index.html#close_http_resp_body
			defer res.Body.Close()
		}
		if err != nil {
			return v, errors.Wrap(err, "failed to do request")
		}

		resBody, err = io.ReadAll(res.Body)
		if err != nil {
			return v, errors.Wrap(err, "failed to read response body")
		}

		err = json.Unmarshal(resBody, &v)
		if err != nil {
			return v, errors.Wrap(err, "failed to unmarshal response body")
		}

		return v, nil
	}, nil
}

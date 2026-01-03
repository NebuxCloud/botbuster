package api

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/NebuxCloud/botbuster/internal/captcha"
	"github.com/NebuxCloud/botbuster/internal/config"
)

type API struct {
	cfg     *config.Config
	log     *slog.Logger
	captcha *captcha.Manager

	server *http.Server
}

func New(cfg *config.Config, log *slog.Logger, captcha *captcha.Manager) *API {
	return &API{
		cfg:     cfg,
		log:     log,
		captcha: captcha,
	}
}

func (api *API) Serve(ctx context.Context) error {
	// Define HTTP request multiplexer
	mux := http.NewServeMux()

	mux.HandleFunc("/_health", api.HealthHandler)
	mux.HandleFunc("/v1/challenge", api.CaptchaChallengeHandler)
	mux.HandleFunc("/v1/verify", api.CaptchaVerifyHandler)

	// Initialize CORS middleware
	crs := cors.New(cors.Options{
		Debug:          api.cfg.Debug,
		AllowedOrigins: api.cfg.AllowedOrigins,
		MaxAge:         24 * 60 * 60, // 24 hours
	})

	// Chain middlewares
	handler := chain(mux,
		noSniffMiddleware,
		crs.Handler,
	)

	// Initialize HTTP server
	api.server = &http.Server{
		Addr: ":" + api.cfg.ListenPort,
		Handler: h2c.NewHandler(
			handler,
			&http2.Server{},
		),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	// Start HTTP server
	api.log.Info(
		"listening to HTTP requests",
		"address", api.server.Addr,
	)

	err := api.server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	api.log.Info("stopped")

	return nil
}

func (api *API) Shutdown(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return api.server.Shutdown(timeoutCtx)
}

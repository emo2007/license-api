package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/emo2007/block-accounting/examples/license-api/internal/pkg/config"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "blockd",
		Version: "0.1.0",
		Flags: []cli.Flag{
			// common
			&cli.StringFlag{
				Name:  "log-level",
				Value: "debug",
			},
			&cli.BoolFlag{
				Name: "log-local",
			},
			&cli.StringFlag{
				Name: "log-file",
			},
			&cli.BoolFlag{
				Name:  "log-add-source",
				Value: true,
			},
			&cli.StringFlag{
				Name: "jwt-secret",
			},

			// rest
			&cli.StringFlag{
				Name:  "rest-address",
				Value: "localhost:3312",
			},
			&cli.BoolFlag{
				Name: "rest-enable-tls",
			},
			&cli.StringFlag{
				Name: "rest-cert-path",
			},
			&cli.StringFlag{
				Name: "rest-key-path",
			},

			// database
			&cli.StringFlag{
				Name: "db-host",
			},
			&cli.StringFlag{
				Name: "db-database",
			},
			&cli.StringFlag{
				Name: "db-user",
			},
			&cli.StringFlag{
				Name: "db-secret",
			},
			&cli.BoolFlag{
				Name: "db-enable-tls",
			},
		},
		Action: func(c *cli.Context) error {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			config := config.Config{
				Common: config.CommonConfig{
					LogLevel:     c.String("log-level"),
					LogLocal:     c.Bool("log-local"),
					LogFile:      c.String("log-file"),
					LogAddSource: c.Bool("log-add-source"),
					JWTSecret:    []byte(c.String("jwt-secret")),
				},
				Rest: config.RestConfig{
					Address: c.String("rest-address"),
					TLS:     c.Bool("rest-enable-tls"),
				},
				DB: config.DBConfig{
					Host:      c.String("db-host"),
					EnableSSL: c.Bool("db-enable-ssl"),
					Database:  c.String("db-database"),
					User:      c.String("db-user"),
					Secret:    c.String("db-secret"),

					CacheHost:   c.String("cache-host"),
					CacheUser:   c.String("cache-user"),
					CacheSecret: c.String("cache-secret"),
				},
			}

			fmt.Println(config)

			service, cleanup, err := factory.ProvideService(config)
			if err != nil {
				panic(err)
			}

			defer func() {
				cleanup()
				service.Stop()
			}()

			if err = service.Run(ctx); err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

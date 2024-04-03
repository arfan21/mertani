package api

import (
	"github.com/arfan21/mertani/config"
	"github.com/arfan21/mertani/internal/server"
	dbpostgres "github.com/arfan21/mertani/pkg/db/postgres"
	"github.com/urfave/cli/v2"
)

func Serve() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Run the API server",
		Action: func(c *cli.Context) error {
			_, err := config.LoadConfig()
			if err != nil {
				return err
			}

			_, err = config.ParseConfig(config.GetViper())
			if err != nil {
				return err
			}

			db, err := dbpostgres.NewPgx()
			if err != nil {
				return err
			}

			server := server.New(
				db,
			)
			return server.Run()
		},
	}

}

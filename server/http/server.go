package http

import (
	"github.com/rrylee/go-tinyid/internal"
	"github.com/rrylee/go-tinyid/server/dbconnection/mysql"
	"github.com/rrylee/go-tinyid/server/service"
	"github.com/spf13/viper"
	"net/http"
)

func Run(configPath string) error {
	config := viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName("tinyid")
	config.AddConfigPath(configPath)

	if err := config.ReadInConfig(); err != nil {
		return err
	}

	if err := mysql.Init(config.GetStringSlice("mysql_dsnes")); err != nil {
		return err
	}

	service.AutoRefreshTokens()

	mux := http.NewServeMux()
	mux.Handle("/tinyid/id/nextId", nextIdHandler())
	mux.Handle("/tinyid/id/nextIdSimple", nextIdSimpleHandler())
	mux.Handle("/tinyid/id/nextSegmentIdSimple", nextSegmentIdSimpleHandler())
	mux.Handle("/tinyid/id/showSegmentId", segmentIdHandler())
	internal.Logf("tinyid server start on %s", config.GetString("server"))
	if err := http.ListenAndServe(config.GetString("server"), mux); err != nil {
		return err
	}

	return nil
}

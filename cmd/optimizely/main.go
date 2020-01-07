/****************************************************************************
 * Copyright 2019, Optimizely, Inc. and contributors                        *
 *                                                                          *
 * Licensed under the Apache License, Version 2.0 (the "License");          *
 * you may not use this file except in compliance with the License.         *
 * You may obtain a copy of the License at                                  *
 *                                                                          *
 *    http://www.apache.org/licenses/LICENSE-2.0                            *
 *                                                                          *
 * Unless required by applicable law or agreed to in writing, software      *
 * distributed under the License is distributed on an "AS IS" BASIS,        *
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. *
 * See the License for the specific language governing permissions and      *
 * limitations under the License.                                           *
 ***************************************************************************/
package main

import (
	"bytes"
	"context"
	"github.com/go-chi/chi"
	"github.com/optimizely/sidedoor/pkg/handler"
	"github.com/optimizely/sidedoor/pkg/middleware"
	"github.com/optimizely/sidedoor/pkg/router"
	"gopkg.in/yaml.v2"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/optimizely/sidedoor/config"
	"github.com/optimizely/sidedoor/pkg/api"
	"github.com/optimizely/sidedoor/pkg/optimizely"
	"github.com/optimizely/sidedoor/pkg/server"
	"github.com/optimizely/sidedoor/pkg/webhook"
)

// Version holds the admin version
var Version string // default set at compile time

func initConfig(v *viper.Viper) error {
	// Set explicit defaults
	v.SetDefault("config.filename", "config.yaml") // Configuration file name

	// Load defaults from the AgentConfig by loading the marshaled values as yaml
	// https://github.com/spf13/viper/issues/188
	defaultConf := config.NewDefaultConfig()
	defaultConf.Admin.Version = Version
	b, err := yaml.Marshal(defaultConf)
	if err != nil {
		return err
	}

	dc := bytes.NewReader(b)
	v.SetConfigType("yaml")
	return v.MergeConfig(dc)
}

func loadConfig(v *viper.Viper) *config.AgentConfig {
	// Configure environment variables
	v.SetEnvPrefix("optimizely")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read configuration from file
	configFile := v.GetString("config.filename")
	v.SetConfigFile(configFile)
	if err := v.MergeInConfig(); err != nil {
		log.Info().Err(err).Msg("Skip loading configuration from config file.")
	}

	conf := &config.AgentConfig{}
	if err := v.Unmarshal(conf); err != nil {
		log.Info().Err(err).Msg("Unable to marshal configuration.")
	}

	return conf
}

func initLogging(conf config.LogConfig) {
	if conf.Pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if lvl, err := zerolog.ParseLevel(conf.Level); err != nil {
		log.Warn().Err(err).Msg("Error parsing log level")
	} else {
		log.Logger = log.Logger.Level(lvl)
	}
}

func main() {
	v := viper.New()
	if err := initConfig(v); err != nil {
		log.Panic().Err(err).Msg("Unable to initialize config")
	}

	conf := loadConfig(v)
	initLogging(conf.Log)

	ctx, cancel := context.WithCancel(context.Background()) // Create default service context
	sg := server.NewGroup(ctx, conf.Server)                 // Create a new server group to manage the individual http listeners
	optlyCache := optimizely.NewCache(ctx, conf.Optly)

	// goroutine to check for signals to gracefully shutdown listeners
	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		// Wait for signal
		sig := <-signalChannel
		log.Info().Msgf("Received signal: %s\n", sig)
		cancel()
	}()

	mw := middleware.CachedOptlyMiddleware{optlyCache}

	r := chi.NewRouter()
	r.Use(mw.ClientCtx)
	r.Route("/features", router.WithFeatureRouter(&handler.FeatureHandler{}))
	r.Route("/experiments", router.WithExperimentRouter(&handler.ExperimentHandler{}))

	r.Route("/users/{userId}", func(r chi.Router) {
		r.Use(mw.ClientCtx)
		router.WithUserRouter(&handler.UserHandler{})(r)
	})

	log.Info().Str("version", conf.Admin.Version).Msg("Starting services.")
	sg.GoListenAndServe("api", conf.API.Port, api.NewDefaultRouter(optlyCache, conf.API))
	sg.GoListenAndServe("webhook", conf.Webhook.Port, webhook.NewRouter(optlyCache, conf.Webhook))
	sg.GoListenAndServe("admin", conf.Admin.Port, router.NewRouter(conf.Admin)) // Admin should be added last.

	// wait for server group to shutdown
	if err := sg.Wait(); err == nil {
		log.Info().Msg("Exiting.")
	} else {
		log.Fatal().Err(err).Msg("Exiting.")
	}

	optlyCache.Wait()
}

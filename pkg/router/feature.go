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

// Package api //
package router

import (
	"github.com/optimizely/sidedoor/pkg/handler"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/optimizely/sidedoor/pkg/middleware"
)

var listFeaturesTimer func(http.Handler) http.Handler
var getFeatureTimer func(http.Handler) http.Handler

func init() {
	listFeaturesTimer = middleware.Metricize("list-features")
	getFeatureTimer = middleware.Metricize("get-feature")
}

// NewRouter returns HTTP API router backed by an optimizely.Cache implementation
func WithFeatureRouter(api handler.FeatureAPI) func(chi.Router) {

	return func(r chi.Router) {
		r.With(listFeaturesTimer).Get("/", api.ListFeatures)
		r.With(getFeatureTimer).Get("/{featureKey}", api.GetFeature)
	}
}

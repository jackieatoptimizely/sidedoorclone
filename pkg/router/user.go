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

var trackEventTimer func(http.Handler) http.Handler
var listUserFeaturesTimer func(http.Handler) http.Handler
var trackUserFeaturesTimer func(http.Handler) http.Handler
var getUserFeatureTimer func(http.Handler) http.Handler
var trackUserFeatureTimer func(http.Handler) http.Handler
var getVariationTimer func(http.Handler) http.Handler
var activateExperimentTimer func(http.Handler) http.Handler
var setForcedVariationTimer func(http.Handler) http.Handler
var removeForcedVariationTimer func(http.Handler) http.Handler

func init() {
	trackEventTimer = middleware.Metricize("track-event")
	listUserFeaturesTimer = middleware.Metricize("list-user-features")
	trackUserFeaturesTimer = middleware.Metricize("track-user-features")
	getUserFeatureTimer = middleware.Metricize("get-user-feature")
	trackUserFeatureTimer = middleware.Metricize("track-user-feature")
	getVariationTimer = middleware.Metricize("get-variation")
	activateExperimentTimer = middleware.Metricize("activate-experiment")
	setForcedVariationTimer = middleware.Metricize("set-forced-variation")
	removeForcedVariationTimer = middleware.Metricize("remove-forced-variation")
}

// NewRouter returns HTTP API router backed by an optimizely.Cache implementation
func WithUserRouter(api handler.UserAPI) func(chi.Router) {

	return func(r chi.Router) {
		r.With(trackEventTimer).Post("/events/{eventKey}", api.TrackEvent)

		r.With(listUserFeaturesTimer).Get("/features", api.ListFeatures)
		r.With(trackUserFeaturesTimer).Post("/features", api.TrackFeatures)
		r.With(getUserFeatureTimer).Get("/features/{featureKey}", api.GetFeature)
		r.With(trackUserFeatureTimer).Post("/features/{featureKey}", api.TrackFeature)
		r.With(getVariationTimer).Get("/experiments/{experimentKey}", api.GetVariation)
		r.With(activateExperimentTimer).Post("/experiments/{experimentKey}", api.ActivateExperiment)
		r.With(setForcedVariationTimer).Put("/experiments/{experimentKey}/variations/{variationKey}", api.SetForcedVariation)
		r.With(removeForcedVariationTimer).Delete("/experiments/{experimentKey}/variations", api.RemoveForcedVariation)
	}
}

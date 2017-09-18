// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"code.cloudfoundry.org/garden"

	"github.com/nwright-nz/openfaas-guardian-backend/metrics"
	"github.com/nwright-nz/openfaas-guardian-backend/requests"
)

// MakeFunctionReader gives a summary of Function structs with Docker service stats overlaid with Prometheus counters.
func MakeFunctionReader(metricsOptions metrics.MetricOptions, c garden.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var m = make(map[string]string)
                services, err := c.Containers(m)
		if err != nil {
			fmt.Println(err)
		}

		// TODO: Filter only "faas" functions (via metadata?)
		var functions []requests.Function

		for _, service := range services {
			functionProp, err := service.Property("function")
			if err != nil {
                           fmt.Printf("error trying to get function property")
                        }
                        if functionProp == "true" {
                        containerName, err := service.Property("name")
				imageName, err := service.Property("image")
				var envProcess string

				
				f := requests.Function{
					Name: containerName,
					Image:           imageName,
					InvocationCount: 0,
					//Replicas:        *service.Spec.Mode.Replicated.Replicas,
					Replicas:   1, //note: doesnt look like garden has replicas, this is handled by CF?
					EnvProcess: envProcess,
				}

				functions = append(functions, f)

				if err != nil {
					print("There was an error retrieving info about the service: ", service)
				}
                           }
			}
		

		functionBytes, _ := json.Marshal(functions)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(functionBytes)
	}
}

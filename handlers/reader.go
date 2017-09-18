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

		//serviceFilter := filters.NewArgs()

		// options := types.ServiceListOptions{
		// 	Filters: serviceFilter,
		// }

		//services, err := c.ServiceList(context.Background(), options)

		var m = make(map[string]string)

		services, err := c.Containers(m)
		if err != nil {
			fmt.Println(err)
		}

		// TODO: Filter only "faas" functions (via metadata?)
		var functions []requests.Function

		for _, service := range services {
			functionProp, err := service.Property("function")
			print(functionProp)
			//if len(service.Spec.TaskTemplate.ContainerSpec.Labels["function"]) > 0 {
			if err != nil {
				print("*********************There are no functions", err)
			} else {
				print(functionProp)
				containerName, err := service.Property("name")
				fmt.Println(err)
				imageName, err := service.Property("image")
				//if len(service.Property("function")) > 0 {
				var envProcess string

				//envs, err := service.Property("env")
				// for _, envs := range service.Spec.TaskTemplate.ContainerSpec.Env {
				// 	if strings.Contains(env, "fprocess=") {
				// 		envProcess = env[len("fprocess="):]
				// 	}
				// }

				f := requests.Function{
					Name: containerName,
					//Image:           service.Spec.TaskTemplate.ContainerSpec.Image,
					Image:           imageName,
					InvocationCount: 0,
					//Replicas:        *service.Spec.Mode.Replicated.Replicas,
					Replicas:   1, //note: doesnt look like garden has replicas, this is handled by CF?
					EnvProcess: envProcess,
				}

				functions = append(functions, f)

				if err != nil {
					print("There was an error doing something, error is: ", service)
				}
			}
		}

		functionBytes, _ := json.Marshal(functions)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(functionBytes)
	}
}

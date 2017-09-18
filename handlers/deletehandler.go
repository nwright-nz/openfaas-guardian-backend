// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"code.cloudfoundry.org/garden"

	"io/ioutil"

	"github.com/nwright-nz/openfaas-guardian-backend/metrics"
	"github.com/nwright-nz/openfaas-guardian-backend/requests"
)

func MakeDeleteFunctionHandler(metricsOptions metrics.MetricOptions, c garden.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		req := requests.DeleteFunctionRequest{}
		defer r.Body.Close()
		reqData, _ := ioutil.ReadAll(r.Body)
		unmarshalErr := json.Unmarshal(reqData, &req)

		if (len(req.FunctionName) == 0) || unmarshalErr != nil {
			log.Printf("Error parsing request to remove service: %s\n", unmarshalErr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("Attempting to remove service %s\n", req.FunctionName)

		serviceName := map[string]string{"name": req.FunctionName}

		services, err := c.Containers(serviceName)
		if err != nil {
			fmt.Println("something something : ", err)
		} else {
			fmt.Println("i can get th service list ok")
		}

		//fmt.Println("********** - Service Name: " + req.FunctionName)

		// TODO: Filter only "faas" functions (via metadata?)
		var serviceIDs []string

		for _, service := range services {

			fmt.Println("i;ve made it to the iteration of serivces")
			serviceProperties, err := service.Property("function")
			if err != nil {
				fmt.Println("Service is not a function", err)
			} else {
				fmt.Println("got the handle: ", serviceProperties)
			}
			nameProperty, err := service.Property("name")
			if err != nil {
				fmt.Println("There is an error getting the container handle..", err)
			}
			isFunction := len(serviceProperties) > 0
			fmt.Println("*************----", isFunction, nameProperty)

			if isFunction && req.FunctionName == nameProperty {
				fmt.Println("Its a function, and the name matches!")
				serviceIDs = append(serviceIDs, nameProperty)
			} else {
				fmt.Println("There is a problem matching the name")
			}
		}

		log.Println(len(serviceIDs))
		if len(serviceIDs) == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("No such service found: %s.", req.FunctionName)))
			return
		}

		var serviceRemoveErrors []error
		for _, serviceID := range serviceIDs {
			fmt.Println(serviceID)
			err := c.Destroy(serviceID)
			if err != nil {
				serviceRemoveErrors = append(serviceRemoveErrors, err)
			}
		}

		if len(serviceRemoveErrors) > 0 {
			log.Printf("Error(s) removing service: %s\n", req.FunctionName)
			log.Println(serviceRemoveErrors)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}

	}
}

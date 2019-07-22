// Copyright 2019 b<>com
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"gopkg.in/Graylog2/go-gelf.v1/gelf"
)

func main() {

	// Define config file
	viper.SetDefault("listenOn", "localhost:5001")
	viper.SetDefault("graylogAddr", "localhost:12201")
	viper.SetConfigName("alertmanager2gelf")       // name of config file (without extension)
	viper.AddConfigPath("/etc/alertmanager2gelf/") // path to look for the config file in

	// Read config file
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	// Checks config parameters
	listenOnConfig := viper.GetString("listenOn")
	if listenOnConfig == "" {
		panic(fmt.Errorf("Fatal error in config file: missing 'listenOn' parameter"))
	}

	graylogAddrConfig := viper.GetString("graylogAddr")
	if graylogAddrConfig == "" {
		panic(fmt.Errorf("Fatal error in config file: missing 'graylogAddrConfig' parameter"))
	}

	// Log service informations on startup
	log.Printf("Service is listening on: '%s'", listenOnConfig)
	log.Printf("Graylog server defined is: '%s'", graylogAddrConfig)

	if graylogAddrConfig != "" {
		gelfWriter, err := gelf.NewWriter(graylogAddrConfig)
		if err != nil {
			log.Fatalf("gelf.NewWriter: %s", err)
		}
		// Log to graylog2
		log.SetFlags(0) // Remove date and time
		log.SetOutput(gelfWriter)
	}

	log.Fatal(http.ListenAndServe(listenOnConfig, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get payload
		promJSON, err := ioutil.ReadAll(r.Body)
		// set payload as string
		spromJSON := string(promJSON)

		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		// Get alerts subset only
		result := gjson.Get(spromJSON, "alerts")
		// iterate over alerts
		result.ForEach(func(key, value gjson.Result) bool {
			log.Printf(value.String())
			return true // keep iterating
		})

	})))
}

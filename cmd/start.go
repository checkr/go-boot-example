// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/checkr/go-boot-example/core"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/checkr/go-boot-example/controllers"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		api := rest.NewApi()

		if viper.GetString("environment") == string(core.Development) {
			api.Use(rest.DefaultDevStack...)
		} else {
			api.Use(&rest.AccessLogJsonMiddleware{})
		}

		api.Use(&rest.CorsMiddleware{
			RejectNonCorsRequests: false,
			OriginValidator: func(origin string, request *rest.Request) bool {
				allowedOrigins := viper.GetStringSlice("allowed_origins")
				if viper.GetString("environment") == string(core.Development) {
					return true
				} else if core.SliceContains(origin, allowedOrigins) {
					return true
				}
				return false
			},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders: []string{
				"Authorization", "Accept", "Content-Type", "X-Custom-Header", "Origin"},
			AccessControlAllowCredentials: true,
			AccessControlMaxAge:           3600,
		})

		api.SetApp(rest.AppSimple(func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson(map[string]string{"Body": "Hello World!"})
		}))

		var routes []*rest.Route
		var controllerRoutes []*rest.Route

		for _, controller := range core.AvailableControllers() {
			controllerRoutes = controller.Routes(api)
			routes = append(routes, controllerRoutes...)
		}
		spew.Dump(routes)
		router, err := rest.MakeRouter(routes...)
		if err != nil {
			log.Fatal(err)
		}

		api.SetApp(router)

		log.Println(fmt.Sprintf("Running HTTP server on 0.0.0.0:%v", viper.GetInt("port")))
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", viper.GetInt("port")), api.MakeHandler()))
	},
}

func init() {
	RootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().IntP("port", "p", 3000, "application port")
	viper.BindPFlag("port", startCmd.PersistentFlags().Lookup("port"))
}

/*
Copyright 2017 Samsung CNCT. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/samsung-cnct/windward/apkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Verbose bool
var ExitCode int

var windwardConfigFilename string

// !? example-app does not explicitly allocate a viper object. It is even
// necessary? Can we just call the library helpers directly?
var windwardConfig = viper.New()

var (
	etcdConfigFilename string
	etcdEndpoints      string
	etcdCaCert         string
	etcdClientCert     string
	etcdClientKey      string
	etcdUuid           string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "windward",
	//Short: "TODO: Short descriptoion of root command",
	Short: "Determine etcd environment for new or existing cluster",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		apkg.Verbose = Verbose
	},

	Run: func(cmd *cobra.Command, args []string) {
		GenerateEtcdConfig()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initWindwardConfig)

	RootCmd.SetHelpCommand(helpCmd)

	RootCmd.PersistentFlags().StringVarP(
		&windwardConfigFilename,
		"in",
		"i",
		"",
		"windward config filename")
	RootCmd.PersistentFlags().StringVarP(
		&etcdConfigFilename,
		"out",
		"o",
		"",
		"etcd config filename")
	// !? Pass in a SRV record instead
	RootCmd.PersistentFlags().StringVarP(
		&etcdEndpoints,
		"endpoints",
		"e",
		"",
		"endpoints")
	RootCmd.PersistentFlags().StringVarP(
		&etcdCaCert,
		"ca-file",
		"a",
		"",
		"Certificate authority filename")
	RootCmd.PersistentFlags().StringVarP(
		&etcdClientCert,
		"cert-file",
		"c",
		"",
		"Client certificate filename")
	RootCmd.PersistentFlags().StringVarP(
		&etcdClientKey,
		"key-file",
		"k",
		"",
		"Client key filename")
	RootCmd.PersistentFlags().StringVarP(
		&etcdUuid,
		"uuid",
		"u",
		"",
		"etcd cluster uuid")
}

// Initializes windwardConfig to use flags, ENV variables and finally configuration files (in that order).
func initWindwardConfig() {
	windwardConfig.BindPFlag("output", RootCmd.Flags().Lookup("output"))
	windwardConfig.BindPFlag("endpoints", RootCmd.Flags().Lookup("endpoints"))
	windwardConfig.BindPFlag("uuid", RootCmd.Flags().Lookup("uuid"))

	windwardConfig.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	windwardConfig.SetEnvPrefix("WINDWARD") // prefix for env vars to configure cluster
	windwardConfig.AutomaticEnv()           // read in environment variables that match

	windwardConfig.SetConfigName(".windward") // name of config file (without extension)
	windwardConfig.AddConfigPath("$HOME")     // path to look for the config file in
	windwardConfig.AddConfigPath(".")         // optionally look for config in the working directory

	configFilename := windwardConfig.GetString("config")
	if configFilename != "" { // enable ability to specify config file via flag
		windwardConfig.SetConfigFile(configFilename)
	}

	// If a config file is found, read it in.
	if err := windwardConfig.ReadInConfig(); err == nil {
		fmt.Println("INFO: Using windward config file:", windwardConfig.ConfigFileUsed())
	}

	// No default for config
	windwardConfig.SetDefault("endpoints", "127.0.0.1:2380")
	windwardConfig.SetDefault("uuid", "default")
	windwardConfig.SetDefault("ca-cert", "/etcd/etcd/ssl/client-ca.pem")
	windwardConfig.SetDefault("cert-file", "/etcd/etcd/ssl/client.pem")
	windwardConfig.SetDefault("key-file", "/etcd/etcd/ssl/client-key.pem")
	windwardConfig.SetDefault("out", "/etc/")
}

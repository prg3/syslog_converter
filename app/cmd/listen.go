/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"log"
	"log/syslog"

	"github.com/spf13/cobra"
	syslogv2 "gopkg.in/mcuadros/go-syslog.v2"
)

//
var target_ip string
var target_port string
var listen_port string

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for syslog and forward",
	Long: `This application listens for syslog messages on UDP and forwards
	them onto the designated server via TCP`,
	Run: func(cmd *cobra.Command, args []string) {
		listen()
	},
}

func listen() {
	fmt.Printf("Forwarding to %v:%v\n", target_ip, target_port)
	channel := make(syslogv2.LogPartsChannel)
	handler := syslogv2.NewChannelHandler(channel)

	// setup the server
	server := syslogv2.NewServer()
	server.SetFormat(syslogv2.RFC5424)
	server.SetHandler(handler)
	server.ListenUDP(fmt.Sprintf("0.0.0.0:%v", listen_port))
	server.Boot()

	// dial the client
	syslog_target := fmt.Sprintf("%v:%v", target_ip, target_port)
	logwriter, e := syslog.Dial("tcp", syslog_target, syslog.LOG_DEBUG, "syslog_forwarder")
	if e != nil {
		log.Fatal(e)
	}

	logwriter.Info("Syslog Forwarder starting up")

	go func(channel syslogv2.LogPartsChannel) {
		for logParts := range channel {
			fmt.Println(logParts)
		}
	}(channel)

	server.Wait()
}

func init() {
	rootCmd.AddCommand(listenCmd)
	listenCmd.Flags().StringVarP(&target_ip, "target_ip", "t", "", "Target IP")
	listenCmd.Flags().StringVarP(&target_port, "target_port", "p", "", "Target Port")
	listenCmd.Flags().StringVarP(&listen_port, "listen_port", "l", "", "Listener Port")
	listenCmd.MarkFlagRequired("target_ip")
	listenCmd.MarkFlagRequired("target_port")
	listenCmd.MarkFlagRequired("listen_port")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

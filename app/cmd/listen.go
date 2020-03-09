/*
Copyright © 2020 NAME HERE paul@majestik.org

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
	"bytes"
	"fmt"
	"log"

	// "log/syslog"
	"net"
	"time"

	"github.com/crewjam/rfc5424"
	"github.com/jeromer/syslogparser/rfc3164"
	"github.com/spf13/cobra"
	"github.com/tidwall/evio"
)

//
var target_ip string
var target_port string
var listen_port string

const maxBufferSize = 1024

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

	var events evio.Events

	// events call
	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {

		// Parse the incoming UDP message as a RFC3164 BSD Syslog message
		p := rfc3164.NewParser(in)
		err := p.Parse()
		if err != nil {
			panic(err)
		}

		// Dump the parsed message into a list so we can rebuild it into a RFC5424 message
		arr := p.Dump()

		// Build a RFC5424 message
		m := rfc5424.Message{
			Priority:  rfc5424.Daemon | rfc5424.Info,
			Timestamp: arr["timestamp"].(time.Time),
			Hostname:  arr["hostname"].(string),
			AppName:   arr["tag"].(string),
			Message:   []byte(arr["content"].(string)),
		}

		// Setup an IO buffer and fill it with the RFC5424 message
		buf := new(bytes.Buffer)
		m.WriteTo(buf)

		// Connect to the TCP Syslog server
		conn, e := net.Dial("tcp", fmt.Sprintf("%v:%v", target_ip, target_port))
		if e != nil {
			log.Fatal(e)
		}

		// Send the message to the Syslog server
		fmt.Fprintf(conn, buf.String())

		// Close the connection and return
		_ = conn.Close()
		return
	}

	// Serve Loop
	if err := evio.Serve(events, fmt.Sprintf("udp://0.0.0.0:%v", listen_port)); err != nil {
		panic(err.Error())
	}
}

func init() {
	rootCmd.AddCommand(listenCmd)
	listenCmd.Flags().StringVarP(&target_ip, "target_ip", "t", "", "Target IP")
	listenCmd.Flags().StringVarP(&target_port, "target_port", "p", "", "Target Port")
	listenCmd.Flags().StringVarP(&listen_port, "listen_port", "l", "", "Listener Port")
	listenCmd.MarkFlagRequired("target_ip")
	listenCmd.MarkFlagRequired("target_port")
	listenCmd.MarkFlagRequired("listen_port")
}

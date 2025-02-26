/*
 * Copyright (c) 2008-2021, Hazelcast, Inc. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License")
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/chroma/quick"
	"github.com/hazelcast/hazelcast-go-client/serialization"
	"github.com/spf13/cobra"

	"github.com/hazelcast/hazelcast-commandline-client/internal"
)

var mapGetCmd = &cobra.Command{
	Use:   "get [--name mapname | --key keyname]",
	Short: "get from map",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(cmd.Context(), time.Second*3)
		defer cancel()
		config, err := internal.MakeConfig()
		// TODO error look like unhandled although it is handled in MakeConfig. Find a better approach
		if err != nil {
			return
		}
		m, err := getMap(ctx, config, mapName)
		if err != nil {
			return
		}
		value, err := m.Get(ctx, mapKey)
		if err != nil {
			fmt.Printf("Cannot get value for key %s from map %s\n", mapKey, mapName)
			isCloudCluster := config.Cluster.Cloud.Enabled
			if networkErrMsg, handled := internal.TranslateNetworkError(err, isCloudCluster); handled {
				fmt.Println("Error: ", networkErrMsg)
				return
			}
			fmt.Printf("Unknown cause: %s\n", err)
			return
		}
		if value != nil {
			switch v := value.(type) {
			case serialization.JSON:
				if err := quick.Highlight(os.Stdout, v.String(),
					"json", "terminal", "tango"); err != nil {
					fmt.Println(v.String())
				}
			default:
				fmt.Println(value)
			}
			return
		}
		fmt.Println("There is no value corresponding to the provided key")
		return
	},
}

func init() {
	decorateCommandWithMapNameFlags(mapGetCmd)
	decorateCommandWithKeyFlags(mapGetCmd)
}

// Copyright © 2019 The Knative Authors
//
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

package service

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	api_errors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/knative/client/pkg/kn/commands"
)

func NewServiceUpdateCommand(p *commands.KnParams) *cobra.Command {
	var editFlags ConfigurationEditFlags
	var waitFlags commands.WaitFlags

	serviceUpdateCommand := &cobra.Command{
		Use:   "update NAME [flags]",
		Short: "Update a service.",
		Example: `
  # Updates a service 'mysvc' with new environment variables
  kn service update mysvc --env KEY1=VALUE1 --env KEY2=VALUE2

  # Update a service 'mysvc' with new port
  kn service update mysvc --port 80

  # Updates a service 'mysvc' with new requests and limits parameters
  kn service update mysvc --requests-cpu 500m --limits-memory 1024Mi`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) != 1 {
				return errors.New("requires the service name.")
			}

			namespace, err := p.GetNamespace(cmd)
			if err != nil {
				return err
			}

			client, err := p.NewClient(namespace)
			if err != nil {
				return err
			}

			var retries = 0
			for {
				name := args[0]
				service, err := client.GetService(name)
				if err != nil {
					return err
				}
				service = service.DeepCopy()

				err = editFlags.Apply(service, cmd)
				if err != nil {
					return err
				}

				err = client.UpdateService(service)
				if err != nil {
					// Retry to update when a resource version conflict exists
					if api_errors.IsConflict(err) && retries < MaxUpdateRetries {
						retries++
						continue
					}
					return err
				}

				if !waitFlags.Async {
					out := cmd.OutOrStdout()
					err := waitForService(client, name, out, waitFlags.TimeoutInSeconds)
					if err != nil {
						return err
					}
				}

				fmt.Fprintf(cmd.OutOrStdout(), "Service '%s' updated in namespace '%s'.\n", args[0], namespace)
				return nil
			}
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preCheck(cmd, args)
		},
	}

	commands.AddNamespaceFlags(serviceUpdateCommand.Flags(), false)
	editFlags.AddUpdateFlags(serviceUpdateCommand)
	waitFlags.AddConditionWaitFlags(serviceUpdateCommand, 60, "Update", "service")
	return serviceUpdateCommand
}

func preCheck(cmd *cobra.Command, args []string) error {
	if cmd.Flags().NFlag() == 0 {
		return errors.New(fmt.Sprintf("flag(s) not set\nUsage: %s", cmd.Use))
	}

	return nil
}

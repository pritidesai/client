// Copyright © 2018 The Knative Authors
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

package plugin

import (
	"testing"

	"github.com/spf13/cobra"
	"gotest.tools/assert"
)

func TestAddPluginFlags(t *testing.T) {
	var (
		pluginFlags *PluginFlags
		cmd         *cobra.Command
	)

	setup := func() {
		pluginFlags = &PluginFlags{}

		cmd = &cobra.Command{}
	}

	t.Run("adds plugin flag", func(t *testing.T) {
		setup()
		pluginFlags.AddPluginFlags(cmd)

		assert.Assert(t, pluginFlags != nil)
		assert.Assert(t, cmd.Flags() != nil)

		nameOnly, err := cmd.Flags().GetBool("name-only")
		assert.Assert(t, err == nil)
		assert.Assert(t, nameOnly == false)
	})
}

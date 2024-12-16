// -----------------------------------------------------------------------------
// SPDX-FileCopyrightText: Copyright © 2024 bomctl a Series of LF Projects, LLC
// SPDX-FileName: cmd/completions.go
// SPDX-FileType: SOURCE
// SPDX-License-Identifier: Apache-2.0
// -----------------------------------------------------------------------------
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
// -----------------------------------------------------------------------------

package cmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bomctl/bomctl/internal/pkg/db"
)

func completions(
	cmd *cobra.Command,
	_ []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	backend := backendFromContext(cmd)

	comps := []string{}

	defer backend.CloseClient()

	documents, err := backend.GetDocumentsByIDOrAlias()
	if err != nil {
		backend.Logger.Fatal(err)
	}

	for _, document := range documents {
		documentID := document.GetMetadata().GetId()
		if slices.Contains(comps, documentID) {
			continue
		}

		if strings.HasPrefix(documentID, toComplete) {
			comps = cobra.AppendActiveHelp(comps, documentID)

			continue
		}

		annotations, err := backend.GetDocumentAnnotations(documentID)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		for _, annotation := range annotations {
			if strings.HasPrefix(annotation.Name, db.AliasAnnotation) &&
				strings.HasPrefix(annotation.Value, toComplete) {
				comps = cobra.AppendActiveHelp(comps, fmt.Sprintf("%s (%s)", documentID, annotation.Value))

				break
			}
		}
	}

	return comps, cobra.ShellCompDirectiveNoFileComp
}

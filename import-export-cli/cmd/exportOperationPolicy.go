/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package cmd

import (
	"fmt"
	"net/http"

	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"path/filepath"
)

var exportAPIPolicyName string
var exportAPIPolicyVersion string
var runningExportAPIPolicyCommand bool

// ExportAPIPolicy command related usage info
const ExportAPIPolicyCmdLiteral = "api"
const exportAPIPolicyCmdShortDesc = "Export API Policies"
const exportAPIPolicyCmdLongDesc = "Export API Policies from an environment"

const exportAPIPolicyCmdExamples = utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportPolicyCmdLiteral + ` ` + ExportAPIPolicyCmdLiteral + ` -n AddHeader -e dev
 NOTE: All the 2 flags (--name (-n) and --environment (-e)) are mandatory.`

// ExportAPIPolicyCmd represents the export policy api command
var ExportAPIPolicyCmd = &cobra.Command{
	Use: ExportAPIPolicyCmdLiteral + " (--name <name-of-the-api-policy> --environment " +
		"<environment-from-which-the-api-policies-should-be-exported>)",
	Short:   exportAPIPolicyCmdShortDesc,
	Long:    exportAPIPolicyCmdLongDesc,
	Example: exportAPIPolicyCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ExportAPIPolicyCmdLiteral + " called")
		var apiPoliciesExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedPoliciesDirName, utils.ExportedAPIPoliciesDirName)

		cred, err := GetCredentials(CmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}

		exportAPIPolicyVersion = utils.APIPolicyVersion
		executeExportAPIPolicyCmd(cred, apiPoliciesExportDirectory, exportAPIPolicyVersion, exportAPIPolicyName)
	},
}

func executeExportAPIPolicyCmd(credential credentials.Credential, exportDirectory string, exportAPIPolicyVersion string, exportAPIPolicyName string) {
	runningExportAPIPolicyCommand = true
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, CmdExportEnvironment)
	if preCommandErr == nil {
		resp, err := impl.ExportAPIPolicyFromEnv(accessToken, CmdExportEnvironment, exportAPIPolicyName, exportAPIPolicyVersion)
		if err != nil {
			utils.HandleErrorAndExit("Error while exporting", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusOK {
			impl.WriteAPIPolicyToFile(exportDirectory, resp, exportAPIPolicyVersion, exportAPIPolicyName, runningExportAPIPolicyCommand)
		} else if resp.StatusCode() == http.StatusInternalServerError {
			// 500 Internal Server Error
			fmt.Println(string(resp.Body()))
		} else {
			// neither 200 nor 500
			fmt.Println("Error exporting API Policies:", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		// error exporting Operarion Policy
		fmt.Println("Error getting OAuth tokens while exporting API Policies:" + preCommandErr.Error())
	}
}

// init using Cobra
func init() {
	ExportPolicyCmd.AddCommand(ExportAPIPolicyCmd)
	ExportAPIPolicyCmd.Flags().StringVarP(&exportAPIPolicyName, "name", "n",
		"", "Name of the API Policy to be exported")
	ExportAPIPolicyCmd.Flags().StringVarP(&CmdExportEnvironment, "environment", "e",
		"", "Environment to which the API Policies should be exported")
	_ = ExportAPIPolicyCmd.MarkFlagRequired("name")
	_ = ExportAPIPolicyCmd.MarkFlagRequired("environment")

}

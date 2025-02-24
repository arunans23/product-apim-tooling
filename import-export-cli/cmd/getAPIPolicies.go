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
	"strconv"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var getAPIPoliciesCmdEnvironment string
var getAPIPoliciesCmdFormat string
var getAPIPolicyListCmdLimit string
var getAllAPIPoliciesAvailable bool

// GetAPIPoliciesCmdLiteral related info
const GetAPIPoliciesCmdLiteral = "api"
const getAPIPoliciesCmdShortDesc = "Display a list of API Policies"

const getAPIPoliciesCmdLongDesc = `Display a list of API Policies in the environment`

var getAPIPoliciesCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetPoliciesCmdLiteral + ` ` + GetAPIPoliciesCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetPoliciesCmdLiteral + ` ` + GetAPIPoliciesCmdLiteral + ` -e dev --all
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetPoliciesCmdLiteral + ` ` + GetAPIPoliciesCmdLiteral + ` -e dev -l 30
 NOTE: The flag (--environment (-e)) is mandatory
 NOTE: Flags (--all) and (--limit (-l)) cannot be used at the same time`

// getAPIPoliciesCmd represents the get policies api command
var getAPIPoliciesCmd = &cobra.Command{
	Use:     GetAPIPoliciesCmdLiteral,
	Short:   getAPIPoliciesCmdShortDesc,
	Long:    getAPIPoliciesCmdLongDesc,
	Example: getAPIPoliciesCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetAPIPoliciesCmdLiteral + " called")
		cred, err := GetCredentials(getAPIPoliciesCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}

		if getAllAPIPoliciesAvailable {
			getAPIPolicyListCmdLimit = ""
		} else {
			limit, err := utils.ValidateFlagWithIntegerValues(getAPIPolicyListCmdLimit)
			if limit < 0 {
				fmt.Println("Limit value should be greater than 0")
			} else if err != nil {
				utils.HandleErrorAndExit("Error converting limit value", err)
			}
		}

		executeGetAPIPoliciesCmd(cred)
	},
}

func executeGetAPIPoliciesCmd(credential credentials.Credential) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, getAPIPoliciesCmdEnvironment)

	if preCommandErr == nil {
		resp, err := impl.GetAPIPolicyListFromEnv(accessToken, getAPIPoliciesCmdEnvironment, getAPIPolicyListCmdLimit)
		if err != nil {
			utils.HandleErrorAndExit("Error while getting API Policies", err)
		}
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())

		if resp.StatusCode() == http.StatusOK {
			impl.PrintAPIPolicies(resp, getAPIPoliciesCmdFormat)
		} else {
			// neither 200 nor 500
			fmt.Println("Error getting API Policies:", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		fmt.Println("Error getting OAuth tokens while getting API Policies:" + preCommandErr.Error())
	}
}

func init() {
	GetPoliciesCmd.AddCommand(getAPIPoliciesCmd)
	getAPIPoliciesCmd.Flags().StringVarP(&getAPIPoliciesCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	getAPIPoliciesCmd.Flags().StringVarP(&getAPIPoliciesCmdFormat, "format", "", "", "Pretty-print API Policies "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	getAPIPoliciesCmd.Flags().StringVarP(&getAPIPolicyListCmdLimit, "limit", "l",
		strconv.Itoa(utils.DefaultPoliciesDisplayLimit), "Maximum number of API Policies to return")
	getAPIPoliciesCmd.Flags().BoolVarP(&getAllAPIPoliciesAvailable, "all", "", false, "Get all API Policies")
	_ = getAPIPoliciesCmd.MarkFlagRequired("environment")
	getAPIPoliciesCmd.MarkFlagsMutuallyExclusive("limit", "all")
}

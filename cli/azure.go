package cli

import (
	"fmt"
	"path/filepath"

	"github.com/BishopFox/cloudfox/azure"
	"github.com/BishopFox/cloudfox/constants"
	"github.com/BishopFox/cloudfox/utils"
	"github.com/aws/smithy-go/ptr"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	AzSubFilter       string
	AzRGFilter        string
	AzOutputFormat    string
	AzOutputDirectory string
	AzVerbosity       int
	AzCommands        = &cobra.Command{
		Use:     "azure",
		Aliases: []string{"az"},
		Long: `
See \"Available Commands\" for Azure Modules`,
		Short: "See \"Available Commands\" for Azure Modules",

		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	AzInstancesCommand = &cobra.Command{
		Use:     "instances",
		Aliases: []string{"instances-map"},
		Short:   "Enumerates compute instances for specified Resource Group",
		Long:    `Enumerates compute instances for specified Resource Group`,
		Run: func(cmd *cobra.Command, args []string) {
			AzRunInstancesCommand(AzSubFilter, AzRGFilter, AzOutputFormat, AzOutputDirectory, AzVerbosity)
		},
	}
	AzUserFilter     string
	AzRBACMapCommand = &cobra.Command{
		Use:     "rbac-map",
		Aliases: []string{"rbac"},
		Short:   "Display all role assignemts for all principals",
		Long:    `Display all role assignemts for all principals`,
		Run: func(cmd *cobra.Command, args []string) {
			color.Red("This command is under development! Use at your own risk!")
			m := azure.RBACMapModule{Scope: utils.AzGetScopeInformation()}
			m.RBACMap(AzVerbosity, AzOutputFormat, AzOutputDirectory, AzUserFilter)
		},
	}
)

func init() {
	// Global flags for the Azure modules
	AzCommands.PersistentFlags().StringVarP(&AzOutputFormat, "output", "o", "all", "[\"table\" | \"csv\" | \"all\" ]")
	AzCommands.PersistentFlags().IntVarP(&AzVerbosity, "verbosity", "v", 1, "1 = Print control messages only\n2 = Print control messages, module output\n3 = Print control messages, module output, and loot file output\n")
	AzCommands.PersistentFlags().StringVar(&AzOutputDirectory, constants.CLOUDFOX_BASE_OUTPUT_DIRECTORY, "cloudfox-output", "Output Directory ")

	// Instance Command Flags
	AzInstancesCommand.Flags().StringVarP(&AzSubFilter, "subscription", "s", "interactive", "Subscription ID")
	AzInstancesCommand.Flags().StringVarP(&AzRGFilter, "resource-group", "g", "interactive", "Resource Group's Name")

	// RBAC Command Flags
	AzRBACMapCommand.Flags().StringVarP(&AzUserFilter, "user", "u", "all", "Display name of user to query")

	AzCommands.AddCommand(AzInstancesCommand, AzRBACMapCommand)
}

func AzRunInstancesCommand(AzSubFilter string, AzRGFilter string, AzOutputFormat string, AzOutputDirectory string, AzVerbosity int) {
	if AzRGFilter == "interactive" && AzSubFilter == "interactive" {
		for _, scopeItem := range azure.ScopeSelection(nil) {

			head, body := azure.GetComputeRelevantData(
				ptr.ToString(scopeItem.Sub.ID),
				ptr.ToString(scopeItem.Rg.Name))

			utils.OutputSelector(
				AzVerbosity,
				AzOutputFormat,
				head,
				body,
				filepath.Join(constants.CLOUDFOX_BASE_OUTPUT_DIRECTORY, fmt.Sprintf("%s_%s", constants.AZ_OUTPUT_DIRECTORY, AzRGFilter)),
				constants.AZ_INTANCES_MODULE_NAME,
				constants.AZ_INTANCES_MODULE_NAME,
				ptr.ToString(scopeItem.Rg.Name))
		}
	} else if AzRGFilter == "interactive" && AzSubFilter != "interactive" {
		// Get all instances from all resource groups in AzSubFilter
		return
	} else if AzRGFilter != "interactive" && AzSubFilter == "interactive" {
		// Get instances for all resource groups with name AzRGFilter
		return
	}
}

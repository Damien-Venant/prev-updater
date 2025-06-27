package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	token        string = ""
	baseUrl      string = ""
	organisation string = ""
	version_tool string = "debug_X.X.X"
)

var rootCommand = &cobra.Command{
	Short: "Update previsional version",
	Long:  "Update previsional version of an user storie on ADO",
	Run:   funcRun,
}

var versionCommand = &cobra.Command{
	Use:  "version",
	Long: "Get the version",
	Run:  funcVersion,
}

func init() {
	rootCommand.Flags().StringVarP(&token, "token", "t", "", "set ADO token (required)")
	rootCommand.Flags().StringVarP(&baseUrl, "base-url", "b", "https://dev.azure.com/", "set base url")
	rootCommand.Flags().StringVarP(&organisation, "organisation", "o", "", "set organisation (required)")
	rootCommand.MarkFlagRequired("token")
	rootCommand.MarkFlagRequired("organisation")
	rootCommand.AddCommand(versionCommand)
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func funcVersion(cmd *cobra.Command, args []string) {
	fmt.Println(version_tool)
}

func funcRun(cmd *cobra.Command, args []string) {

}

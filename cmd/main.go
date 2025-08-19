package main

import (
	"fmt"
	"os"

	"github.com/prev-updater/internal/infra"
	"github.com/prev-updater/internal/repository"
	"github.com/spf13/cobra"
)

var (
	token        string = ""
	baseUrl      string = ""
	organisation string = ""
	version_tool string = "debug_X.X.X"
	pipelineId   int32
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
	token = os.Getenv("PREV_UPDATER_TOKEN")
	baseUrl = os.Getenv("PREV_UPDATER_BASEURL")

	rootCommand.Flags().StringVarP(&token, "token", "t", "", "set ADO token (required)")
	rootCommand.Flags().StringVarP(&baseUrl, "base-url", "b", "https://dev.azure.com/", "set base url")
	rootCommand.Flags().StringVarP(&organisation, "organisation", "o", "", "set organisation (required)")
	rootCommand.Flags().Int32VarP(&pipelineId, "pipeline-id", "p", 0, "set pipeline id")
	rootCommand.MarkFlagRequired("token")
	rootCommand.MarkFlagRequired("organisation")
	rootCommand.MarkFlagRequired("pipeline-id")
	rootCommand.AddCommand(versionCommand)
}

func main() {

	infra.ConfigureHttpClient(&infra.HttpClientConfiguration{
		BaseUrl: baseUrl,
		Token:   token,
	})

	client := infra.GetHttpClient()

	repo := repository.New(client)
	repo.GetPipelineRuns(862)
	//if err := rootCommand.Execute(); err != nil {
	//	fmt.Println(err)
	//	os.Exit(-1)
	//}
}

func funcVersion(cmd *cobra.Command, args []string) {
	fmt.Println(version_tool)
}

func funcRun(cmd *cobra.Command, args []string) {

}

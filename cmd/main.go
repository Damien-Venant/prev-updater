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
	project      string = ""
	versionTool  string = "debug_X.X.X"
	pipelineId   int32
)

var rootCommand = &cobra.Command{
	Short: "Update previsional version",
	Long:  "Update previsional version of an user storie on ADO",
	Run:   funcRun,
}

var versionCommand = &cobra.Command{
	Use:   "version",
	Long:  "Get the version",
	Short: "Get the version",
	Run:   funcVersion,
}

var launchCommand = &cobra.Command{
	Use:   "start",
	Short: "Start the prev-updater",
	Long:  "Start the prev-updater command to changes fields in ADO cards",
	Run:   funcStart,
}

func init() {
	token = os.Getenv("PREV_UPDATER_TOKEN")
	baseUrl = os.Getenv("PREV_UPDATER_BASEURL")

	launchCommand.Flags().StringVarP(&token, "token", "t", "", "set ADO token (required)")
	launchCommand.Flags().StringVarP(&baseUrl, "base-url", "b", "https://dev.azure.com/", "set base url")
	launchCommand.Flags().StringVarP(&organisation, "organisation", "o", "", "set organisation (required)")
	launchCommand.Flags().Int32VarP(&pipelineId, "pipeline-id", "i", 0, "set pipeline id")
	launchCommand.Flags().StringVarP(&project, "project", "p", "", "project name")

	launchCommand.MarkFlagRequired("token")
	launchCommand.MarkFlagRequired("organisation")
	launchCommand.MarkFlagRequired("project")
	launchCommand.MarkFlagRequired("pipeline-id")

	rootCommand.AddCommand(versionCommand)
	rootCommand.AddCommand(launchCommand)
}

func main() {

	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}

func funcVersion(cmd *cobra.Command, args []string) {
	fmt.Println(versionTool)
}

func funcRun(cmd *cobra.Command, args []string) {

}

func funcStart(cmd *cobra.Command, args []string) {
	url := fmt.Sprintf("%s/%s/%s/", baseUrl, organisation, project)
	infra.ConfigureHttpClient(&infra.HttpClientConfiguration{
		BaseUrl: url,
		Token:   token,
	})
	client := infra.GetHttpClient()

	repo := repository.New(client)
	lastRun, err := repo.GetPipelineRuns(862)
	if err != nil {
		panic(err)
	}

	fmt.Println(lastRun[0].Name)
}

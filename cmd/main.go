package cmd

import (
	"fmt"
	"os"

	"github.com/Damien-Venant/prev-updater/internal/infra"
	"github.com/Damien-Venant/prev-updater/internal/repository"
	"github.com/Damien-Venant/prev-updater/internal/usescases"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const (
	EXIT_FAILURE = -1
	EXIT_SUCCESS = 0
)

var (
	token        string = ""
	baseUrl      string = ""
	organisation string = ""
	project      string = ""
	versionTool  string = "debug_X.X.X"
	pipelineId   int32
	repositoryId string
	fieldName    string

	logger *zerolog.Logger = nil
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
	Run:   funcStartBatching,
}

func init() {
	launchCommand.Flags().StringVarP(&token, "token", "t", "", "set ADO token (required)")
	launchCommand.Flags().StringVarP(&baseUrl, "base-url", "b", "https://dev.azure.com/", "set base url")
	launchCommand.Flags().StringVarP(&organisation, "organisation", "o", "", "set organisation")
	launchCommand.Flags().Int32VarP(&pipelineId, "pipeline-id", "i", 0, "set pipeline id")
	launchCommand.Flags().StringVarP(&project, "project", "p", "", "project name")
	launchCommand.Flags().StringVarP(&repositoryId, "repository", "r", "", "set repository id")
	launchCommand.Flags().StringVarP(&fieldName, "field", "f", "", "set field name")

	launchCommand.MarkFlagRequired("token")
	launchCommand.MarkFlagRequired("organisation")
	launchCommand.MarkFlagRequired("project")
	launchCommand.MarkFlagRequired("pipeline-id")
	launchCommand.MarkFlagRequired("repository")
	launchCommand.MarkFlagRequired("field")

	rootCommand.AddCommand(versionCommand)
	rootCommand.AddCommand(launchCommand)

	_, err := infra.ConfigDirectory()
	if err != nil {
		panic(err)
	}
	loggerWriter, err := infra.OpenLogFile()
	if err != nil {
		panic(err)
	}
	logger = infra.NewLogger(loggerWriter)
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(exitWithError())
	}
	os.Exit(EXIT_SUCCESS)
}

func funcVersion(cmd *cobra.Command, args []string) {
	fmt.Println(versionTool)
}

func funcRun(cmd *cobra.Command, args []string) {

}

func funcStartBatching(cmd *cobra.Command, args []string) {
	url := fmt.Sprintf("%s/%s/%s/", baseUrl, organisation, project)
	infra.ConfigureHttpClient(&infra.HttpClientConfiguration{
		BaseUrl: url,
		Token:   token,
	}, logger)
	client := infra.GetHttpClient()
	repo := repository.New(client)

	use := usescases.NewAdoUsesCases(repo, logger)

	if err := use.UpdateFieldsByLastRuns(int(pipelineId), repositoryId, fieldName); err != nil {
		logger.
			Err(err).
			Send()
		os.Exit(exitWithError())
	}
}

func exitWithError() int {
	infra.CloseLogFile()

	return EXIT_FAILURE
}

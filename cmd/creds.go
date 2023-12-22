package cmd

import (
	"cli/api"
	"fmt"

	"github.com/spf13/cobra"
)

var getCredsCmd = &cobra.Command{
	Use:   "creds",
	Short: "get image registry creds",
	Long:  "get image registry creds",
	Run: func(cmd *cobra.Command, args []string) {
		creds, err := api.GetCreds()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(creds)

	},
}

func init() {

}

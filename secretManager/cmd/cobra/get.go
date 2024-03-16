package cobra

import (
	"fmt"

	"github.com/okaberintaroubeta/secretManager/secret"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a secret in your local strage path",

	Run: func(cmd *cobra.Command, args []string) {
		v := secret.File(encodingKey, secretsPath())
		key, _ := cmd.Flags().GetString("key")
		val, err := v.Get(key)
		if err != nil {
			fmt.Println("No value set for this key")
			return
		} else {
			fmt.Printf("%s= %s\n", key, val)
		}
	},
}

func init() {

	getCmd.Flags().String("key", "", "Key to retrieve its value")
	RootCmd.AddCommand(getCmd)
}

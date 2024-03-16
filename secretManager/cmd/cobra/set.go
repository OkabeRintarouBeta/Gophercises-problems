package cobra

import (
	"fmt"

	"github.com/okaberintaroubeta/secretManager/secret"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a secret in your local strage path",

	Run: func(cmd *cobra.Command, args []string) {
		v := secret.File(encodingKey, secretsPath())

		key, _ := cmd.Flags().GetString("key")
		val, _ := cmd.Flags().GetString("value")

		if key == "" || val == "" {
			fmt.Println("Both key and value must be provided")
			return
		}

		err := v.Set(key, val)
		if err != nil {
			fmt.Println("Error when setting key value pair")
			return
		}
		fmt.Println("Set successfully!")
	},
}

func init() {

	setCmd.Flags().String("key", "", "Key to set value for")
	setCmd.Flags().String("value", "", "Value for the key to be set")
	RootCmd.AddCommand(setCmd)
}

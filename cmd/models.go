package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "Manage models",
}

var modelsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available models",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		resp, err := client.ListModels(cmd.Context())
		if err != nil {
			return fmt.Errorf("list models: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tLABEL\tVENDOR\tTYPE\tBYOK")
		for _, m := range resp.Models {
			mType := ""
			if m.Type != nil {
				mType = *m.Type
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\n", m.ModelName, m.Label, m.Vendor, mType, m.IsByok)
		}
		w.Flush()
		return nil
	},
}

func init() {
	modelsCmd.AddCommand(modelsListCmd)
	rootCmd.AddCommand(modelsCmd)
}

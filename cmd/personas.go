package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var personasCmd = &cobra.Command{
	Use:   "personas",
	Short: "Manage personas",
}

var personasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available personas",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		resp, err := client.ListPersonas(cmd.Context())
		if err != nil {
			return fmt.Errorf("list personas: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tMODEL")
		for _, p := range resp.Personas {
			model := ""
			if p.Model != nil {
				model = *p.Model
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", p.ID, p.Name, model)
		}
		w.Flush()
		return nil
	},
}

func init() {
	personasCmd.AddCommand(personasListCmd)
	rootCmd.AddCommand(personasCmd)
}

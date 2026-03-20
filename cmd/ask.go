package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/irfansofyana/zo-cli/api"
	"github.com/spf13/cobra"
)

var (
	askModel          string
	askConversationID string
	askPersona        string
	askOutputFormat   string
)

var askCmd = &cobra.Command{
	Use:   "ask [message]",
	Short: "Send a message to Zo",
	Long:  "Send a single message to Zo and get a response. Use --conversation-id to continue a conversation.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAPIKey(); err != nil {
			return err
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		req := api.AskRequest{
			Input:          strings.Join(args, " "),
			ConversationID: askConversationID,
			ModelName:      askModel,
			PersonaID:      askPersona,
		}

		if askOutputFormat != "" {
			data, err := os.ReadFile(askOutputFormat)
			if err != nil {
				return fmt.Errorf("read output format file: %w", err)
			}
			raw := json.RawMessage(data)
			req.OutputFormat = &raw
		}

		resp, err := client.Ask(cmd.Context(), req)
		if err != nil {
			return err
		}

		// Print output
		var strOutput string
		if json.Unmarshal(resp.Output, &strOutput) == nil {
			fmt.Fprintln(cmd.OutOrStdout(), strOutput)
		} else {
			// Structured output - print as pretty JSON
			var pretty json.RawMessage
			if json.Unmarshal(resp.Output, &pretty) == nil {
				formatted, _ := json.MarshalIndent(pretty, "", "  ")
				fmt.Fprintln(cmd.OutOrStdout(), string(formatted))
			}
		}

		// Print conversation ID to stderr for scripting
		if resp.ConversationID != "" {
			fmt.Fprintf(os.Stderr, "conversation_id: %s\n", resp.ConversationID)
		}

		return nil
	},
}

func init() {
	askCmd.Flags().StringVar(&askModel, "model", "", "Override the default model")
	askCmd.Flags().StringVar(&askConversationID, "conversation-id", "", "Continue an existing conversation")
	askCmd.Flags().StringVar(&askPersona, "persona", "", "Override the active persona")
	askCmd.Flags().StringVar(&askOutputFormat, "output-format", "", "Path to JSON schema file for structured output")
	rootCmd.AddCommand(askCmd)
}

package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/irfansofyana/zo-cli/api"
	"github.com/spf13/cobra"
)

var (
	chatModel   string
	chatPersona string
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start an interactive conversation with Zo",
	Long:  "Enter a REPL loop where you can chat continuously with Zo. Conversation context is maintained automatically. Type 'exit' or 'quit' to end.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAPIKey(); err != nil {
			return err
		}

		client, err := getClient()
		if err != nil {
			return err
		}
		return chatLoop(cmd.Context(), client, chatModel, chatPersona, os.Stdin, cmd.OutOrStdout())
	},
}

func chatLoop(ctx context.Context, client api.ZoClient, model, persona string, in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	var conversationID string

	fmt.Fprintln(out, "Zo Chat — type 'exit' or 'quit' to end")
	fmt.Fprintln(out, "")

	for {
		fmt.Fprint(out, "you> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "exit" || input == "quit" {
			fmt.Fprintln(out, "Goodbye!")
			break
		}

		req := api.AskRequest{
			Input:          input,
			ConversationID: conversationID,
			ModelName:      model,
			PersonaID:      persona,
		}

		resp, err := client.Ask(ctx, req)
		if err != nil {
			fmt.Fprintf(out, "Error: %v\n\n", err)
			continue
		}

		// Print response
		var strOutput string
		if json.Unmarshal(resp.Output, &strOutput) == nil {
			fmt.Fprintf(out, "zo> %s\n\n", strOutput)
		} else {
			formatted, _ := json.MarshalIndent(resp.Output, "", "  ")
			fmt.Fprintf(out, "zo> %s\n\n", string(formatted))
		}

		// Track conversation
		if resp.ConversationID != "" {
			conversationID = resp.ConversationID
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read input: %w", err)
	}
	return nil
}

func init() {
	chatCmd.Flags().StringVar(&chatModel, "model", "", "Override the default model")
	chatCmd.Flags().StringVar(&chatPersona, "persona", "", "Override the active persona")
	rootCmd.AddCommand(chatCmd)
}

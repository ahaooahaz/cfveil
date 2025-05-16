package python

import (
	"fmt"
	"os"

	"github.com/ahaooahaz/cfveil/internal/python"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "python",
	Short:   "python project.",
	Long:    `python project.`,
	Aliases: []string{"py"},
	Run: func(cmd *cobra.Command, args []string) {
		if *arg_INPUT == "" || *arg_OUTPUT == "" {
			fmt.Println("invalid")
			return
		}

		err := python.Process(*arg_INPUT, *arg_OUTPUT, *arg_EXCLUDE)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		return
	},
}

var (
	arg_INPUT, arg_OUTPUT *string
	arg_EXCLUDE           *[]string
)

func init() {
	arg_INPUT = Cmd.Flags().StringP("input", "i", "", "project dir or file path")
	arg_OUTPUT = Cmd.Flags().StringP("output", "o", "", "output dir")
	arg_EXCLUDE = Cmd.Flags().StringSliceP("exclude", "e", []string{}, "exclude files")
}

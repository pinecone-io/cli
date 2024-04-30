package pcio

import (
	"github.com/spf13/cobra"
)

var Messages = (&cobra.Command{}).OutOrStdout()

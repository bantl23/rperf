package cmd

import (
	"github.com/bantl23/yabba/run"
	"github.com/spf13/cobra"
)

var (
	listenAddrs []string
	listenSize  uint64
)

func init() {
	listenCmd.Flags().StringSliceVarP(&listenAddrs, "addrs", "a", []string{":5201"}, "bind address(es)")
	listenCmd.Flags().Uint64VarP(&listenSize, "size", "s", 128*1024, "buffer size")
	rootCmd.AddCommand(listenCmd)
}

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listens for clients",
	Run: func(cmd *cobra.Command, args []string) {
		server := run.Server{
			Addresses: listenAddrs,
			Size:      listenSize,
		}
		server.Run()
	},
}

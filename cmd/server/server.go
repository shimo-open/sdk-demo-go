package server

import (
	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egovernor"
	"github.com/spf13/cobra"

	"sdk-demo-go/cmd"
	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/server/http"
)

// CmdRun is the cobra command for running the HTTP server
var CmdRun = &cobra.Command{
	Use:   "server",
	Short: "Start sdk-demo-go HTTP API server",
	Long:  `Start sdk-demo-go HTTP API server`,
	Run:   CmdFunc,
}

func init() {
	CmdRun.InheritedFlags()
	cmd.RootCommand.AddCommand(CmdRun)
}

// CmdFunc is the entry point for the server command
func CmdFunc(cmd *cobra.Command, args []string) {
	e := ego.New()
	e.Invoker(invoker.Init)
	if err := e.Serve(
		egovernor.Load("server.governor").Build(),
		http.ServeHTTP(),
	).Run(); err != nil {
		elog.Panic("Startup failed", l.E(err))
	}
}

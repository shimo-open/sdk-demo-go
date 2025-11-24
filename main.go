package main

import (
	"fmt"
	"os"

	_ "github.com/gotomicro/gorm-to-dm"

	"sdk-demo-go/cmd"
	_ "sdk-demo-go/cmd"
	_ "sdk-demo-go/cmd/sdkctl"
	_ "sdk-demo-go/cmd/server"
)

func main() {
	if err := cmd.RootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return
}

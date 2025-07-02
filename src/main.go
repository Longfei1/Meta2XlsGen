package main

import (
	"Meta2XlsGen/src/cmd"
	"Meta2XlsGen/src/logic"
	"fmt"
)

func main() {
	cmdArgs, err := cmd.ParseCmdArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	if !cmdArgs.SuccessRun {
		return
	}

	l := logic.NewGenCodeLogic(cmdArgs)
	err = l.Run()
	if err != nil {
		fmt.Println(err)
	}
}

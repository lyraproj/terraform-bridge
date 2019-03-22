package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/lyraproj/pcore/pcore"
	"github.com/lyraproj/pcore/px"
	"github.com/lyraproj/servicesdk/grpc"
)

func main() {
	pcore.Do(func(c px.Context) {
		if len(os.Args) < 2 {
			fmt.Println("You must specify the name of a goplugin binary as a command line argument e.g. build/goplugins/terraformazurerm")
			return
		}
		service, _ := grpc.Load(exec.Command(os.Args[1]), hclog.Default())
		fmt.Println("----------------------------------------------------------------")
		fmt.Println("Identifier")
		fmt.Println("----------------------------------------------------------------")
		fmt.Println(service.Identifier(c))
		fmt.Println("----------------------------------------------------------------")
		fmt.Println("Metadata")
		fmt.Println("----------------------------------------------------------------")
		fmt.Println(service.Metadata(c))
	})
}

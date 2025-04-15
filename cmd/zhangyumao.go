package main

import (
	"flag"

	"zhangyumao/apps/app"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/ergo/lib"
	"github.com/spf13/viper"
)

var (
	OptionNodeName   string
	OptionNodeCookie string
)

func init() {
	flag.StringVar(&OptionNodeName, "name", "ZhangyuMao@localhost", "node name")
	flag.StringVar(&OptionNodeCookie, "cookie", lib.RandomString(16), "a secret cookie for the network messaging")
}

func main() {
	var options gen.NodeOptions

	flag.Parse()

	viper.AutomaticEnv()

	// create applications that must be started
	apps := []gen.ApplicationBehavior{
		app.CreateApp(),
	}
	options.Applications = apps

	// set network options
	options.Network.Cookie = OptionNodeCookie

	// starting node
	node, err := ergo.StartNode(gen.Atom(OptionNodeName), options)
	if err != nil {
		panic(err)
	}

	node.Wait()
}

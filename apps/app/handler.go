package app

import (
	"fmt"
	"net/http"
	"zhangyumao/internal/github"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"ergo.services/ergo/meta"
)

var alias_Handler = gen.Atom("handler")

func factory_Handler() gen.ProcessBehavior {
	return &Handler{}
}

type HandlerInitArgs struct {
	ClientCreator github.ClientCreator
}

type Handler struct {
	act.Pool
	clientCreator github.ClientCreator
}

// Init invoked on a start this process.
func (w *Handler) Init(args ...any) (act.PoolOptions, error) {
	var webOptions meta.WebServerOptions
	var poolOptions act.PoolOptions

	initError := fmt.Errorf("unable to start Handler. Expected Args to be of type []HandlerWorkerInitArgs, got: %v", args)
	if len(args) < 1 {
		return poolOptions, initError
	}

	switch initArgs := args[0].(type) {
	case HandlerInitArgs:
		w.clientCreator = initArgs.ClientCreator
	default:
		return poolOptions, initError
	}

	mux := http.NewServeMux()

	// create and spawn root handler meta-process.
	root := meta.CreateWebHandler(meta.WebHandlerOptions{})
	rootid, err := w.SpawnMeta(root, gen.MetaOptions{})
	if err != nil {
		w.Log().Error("unable to spawn WebHandler meta-process: %s", err)
		return poolOptions, err
	}

	// add it to the mux. you can also use middleware functions:
	// mux.Handle("/", middleware(root))
	mux.Handle("/", root)
	w.Log().Info("started WebHandler to serve '/' (meta-process: %s)", rootid)

	webOptions.Port = 3000
	webOptions.Host = "0.0.0.0"

	webOptions.Handler = mux

	webserver, err := meta.CreateWebServer(webOptions)
	if err != nil {
		w.Log().Error("unable to create Web server meta-process: %s", err)
		return poolOptions, err
	}
	webserverid, err := w.SpawnMeta(webserver, gen.MetaOptions{})
	if err != nil {
		// invoke Terminate to close listening socket
		webserver.Terminate(err)
		return poolOptions, err
	}

	https := "http"
	if webOptions.CertManager != nil {
		https = "https"
	}
	w.Log().Info("started Web server %s: use %s://%s:%d/", webserverid, https, webOptions.Host, webOptions.Port)
	w.Log().Info("you may check it with command below:")
	w.Log().Info("   $ curl -k %s://%s:%d", https, webOptions.Host, webOptions.Port)

	poolOptions.WorkerFactory = factory_HandlerWebWorker
	poolOptions.WorkerArgs = []any{
		HandlerWorkerInitArgs{ClientCreator: w.clientCreator},
	}
	return poolOptions, nil
}

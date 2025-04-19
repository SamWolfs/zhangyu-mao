package app

import (
	"fmt"
	"net/http"
	"zhangyumao/config"
	"zhangyumao/internal/github"
	"zhangyumao/internal/handlers"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

func factory_HandlerWebWorker() gen.ProcessBehavior {
	return &HandlerWebWorker{}
}

type HandlerWorkerInitArgs struct {
	ClientCreator github.ClientCreator
}

type HandlerWebWorker struct {
	act.WebWorker
	clientCreator github.ClientCreator
}

// Init invoked on a start this process.
func (w *HandlerWebWorker) Init(args ...any) error {
	initError := fmt.Errorf("unable to start HandlerWorker. Expected Args to be of type []HandlerWorkerInitArgs, got: %v", args)
	if len(args) < 1 {
		return initError
	}

	switch initArgs := args[0].(type) {
	case HandlerWorkerInitArgs:
		w.clientCreator = initArgs.ClientCreator
	default:
		return initError
	}
	w.Log().Info("started web worker process with args %v", args)
	return nil
}

func (w *HandlerWebWorker) HandleGet(from gen.PID, writer http.ResponseWriter, request *http.Request) error {
	w.Log().Info("got HTTP GET request %q", request.URL.Path)
	writer.Header().Set("Content-Type", "application/json")
	return nil
}

func (w *HandlerWebWorker) HandlePost(from gen.PID, writer http.ResponseWriter, request *http.Request) error {
	cfg := config.GetGitHubAppConfig()

	payload, err := github.ValidatePayload(request, []byte(cfg.WebhookSecret))
	if err != nil {
		return err
	}

	event, err := github.ParseWebHook(github.WebHookType(request), payload)
	if err != nil {
		return err
	}

	switch event := event.(type) {
	case *github.PushEvent:
		client, err := w.clientCreator.NewInstallationClient(*event.Installation.ID)
		if err != nil {
			return err
		}

		err = handlers.HandlePushEvent(client, event)
		if err != nil {
			return err
		}
	}

	writer.Header().Set("Content-Type", "application/json")

	return nil
}

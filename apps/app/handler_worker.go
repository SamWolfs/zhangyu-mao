package app

import (
	"bytes"
	"encoding/json"
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
	if len(args) < 1 {
		return fmt.Errorf("unable to start HandlerWorker. Expected Args to be of type []HandlerWorkerInitArgs, got: %v", args)
	}

	switch initArgs := args[0].(type) {
	case HandlerWorkerInitArgs:
		w.clientCreator = initArgs.ClientCreator
	default:
		return fmt.Errorf("unable to start HandlerWorker. Expected Args to be of type []HandlerWorkerInitArgs, got: %v", args)
	}
	w.Log().Info("started web worker process with args %v", args)
	return nil
}

// Handle GET requests. For the other HTTP methods (POST, PATCH, etc)
// you need to add the accoring callback-method implementation. See act.WebWorkerBehavior.

func (w *HandlerWebWorker) HandleGet(from gen.PID, writer http.ResponseWriter, request *http.Request) error {
	var buf bytes.Buffer

	w.Log().Info("got HTTP request %q", request.URL.Path)
	writer.Header().Set("Content-Type", "application/json")
	// response JSON message with information about this process
	info, _ := w.Info()
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(info)
	checkError(err)
	_, err = writer.Write(buf.Bytes())
	checkError(err)
	return nil
}

func (w *HandlerWebWorker) HandlePost(from gen.PID, writer http.ResponseWriter, request *http.Request) error {
	var buf bytes.Buffer

	cfg := config.GetGitHubAppConfig()

	payload, err := github.ValidatePayload(request, []byte(cfg.WebhookSecret))
	checkError(err)

	event, err := github.ParseWebHook(github.WebHookType(request), payload)
	checkError(err)
	switch event := event.(type) {
	case *github.PushEvent:
		client, err := w.clientCreator.NewInstallationClient(*event.Installation.ID)
		checkError(err)
		err = handlers.HandlePushEvent(client, event)
		checkError(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	// response JSON message with information about this process
	info, _ := w.Info()
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err = enc.Encode(info)
	checkError(err)
	_, err = writer.Write(buf.Bytes())
	checkError(err)
	return nil
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

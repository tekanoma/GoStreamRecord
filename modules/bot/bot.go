package bot

import (
	"GoRecordurbate/modules/web/provider"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"sync"
)

// Bot encapsulates the recording bot’s state.
type controller struct {
	mux        sync.Mutex
	status     []Recorder
	isFirstRun bool
	logger     *log.Logger
	// ctx is used to signal shutdown.
	ctx    context.Context
	cancel context.CancelFunc
}
type Recorder struct {
	Web         *Provider
	WebType     string    `json:"web_type"` // name of interface were working with
	Enabled     bool      `json:"enabled"`
	StopStatus  bool      `json:"-"`
	WasRestart  bool      `json:"restarting"`
	Name        string    `json:"name"`
	Cmd         *exec.Cmd `json:"-"`
	IsRecording bool      `json:"isRecording"`
}
type Provider struct {
	Type     string             `json:"type"`     // Which website we are using from the supported ones in web/provider
	Url      string             `json:"url"`      // Which website we are using from the supported ones in web/provider
	Username string             `json:"username"` // Which website we are using from the supported ones in web/provider
	Site     provider.IProvider `json:"-"`
}

var Bot *controller

func Init() *controller {
	Bot = NewBot(log.New(os.Stdout, "lpg.log", log.LstdFlags))
	return Bot
}

// NewBot creates a new Bot, sets up its cancellation context.
func NewBot(logger *log.Logger) *controller {
	ctx, cancel := context.WithCancel(context.Background())
	b := &controller{
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
		status:     []Recorder{},
		isFirstRun: true,
	}
	return b
}

func NewProvider(webType, username string) *Provider {
	// Create a new instance of the provider (assuming "bc" is the type name)
	prov := provider.Init(webType)

	// Initialize the provider with the given webType and username.
	// The Init method returns a value (likely *bc) that we can marshal.
	initResult := prov.Init(webType, username)

	// Marshal the initialized provider into JSON.
	data, err := json.Marshal(initResult)
	if err != nil {
		return nil
	}

	// Unmarshal the JSON into your Provider struct.
	var p Provider
	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil
	}

	// Keep a reference to the underlying provider instance.
	p.Site = prov
	return &p
}

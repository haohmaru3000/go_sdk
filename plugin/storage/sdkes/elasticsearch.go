// github: https://github.com/olivere/elastic

package sdkes

import (
	"context"
	"flag"
	"os"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/haohmaru3000/go_sdk/logger"
)

const (
	DefaultElasticSearchName = "DefaultES"
)

type ESConfig struct {
	Prefix         string
	URL            string
	HasSniff       bool
	HasHealthCheck bool
	Index          string
	Username       string
	Password       string
}

type es struct {
	name   string
	client *elastic.Client
	logger logger.Logger
	*ESConfig
}

func NewES(name, flagPrefix string) *es {
	return &es{
		name: name,
		ESConfig: &ESConfig{
			Prefix: flagPrefix,
		},
	}
}

func (es *es) GetDefaultES() *es {
	return NewES(DefaultElasticSearchName, "")
}

func (es *es) GetPrefix() string {
	return es.Prefix
}

func (es *es) GetIndex() string {
	return es.Index
}

func (es *es) isDisabled() bool {
	return es.URL == ""
}

func (es *es) InitFlags() {
	prefix := es.Prefix
	if es.Prefix != "" {
		prefix += "-"
	}
	flag.StringVar(&es.Index, prefix+"es-index", "", "elastic search index")
	flag.StringVar(&es.URL, prefix+"es-url", "", "elastic search connection-string. ex: http://localhost:9200")
	flag.BoolVar(&es.HasSniff, prefix+"es-has-sniff", false, "elastic search sniffing mode. default: false")
	flag.BoolVar(&es.HasHealthCheck, prefix+"es-has-health-check", false, "elastic search health-check mode. default: false")
	flag.StringVar(&es.Username, prefix+"es-username", "", "elasticsearch username")
	flag.StringVar(&es.Password, prefix+"es-password", "", "elasticsearch password")
	flag.Parse()
}

func (es *es) Configure() error {
	if es.isDisabled() {
		return nil
	}
	es.logger = logger.GetCurrent().GetLogger(es.name)
	es.logger.Info("connecting to elastic search at ", es.URL, "...")

	var client *elastic.Client
	var err error

	if es.Username != "" && es.Password != "" {
		client, err = elastic.NewClient(
			elastic.Config{
				Addresses: []string{es.URL},
				Username:  es.Username,
				Password:  es.Password,
				Logger: &elastictransport.ColorLogger{
					Output:             os.Stdout,
					EnableRequestBody:  true,
					EnableResponseBody: true,
				},
			},
		)
	} else {
		client, err = elastic.NewClient(
			elastic.Config{
				Addresses: []string{es.URL},
				Logger: &elastictransport.ColorLogger{
					Output:             os.Stdout,
					EnableRequestBody:  true,
					EnableResponseBody: true,
				},
			},
		)
	}

	if err != nil {
		es.logger.Error("cannot connect to elastic search. ", err.Error())
		return err
	}

	// Connect successfully, assign client
	es.client = client
	return nil
}

func (es *es) Name() string {
	return es.name
}

func (es *es) Get() interface{} {
	return es.client
}

func (es *es) Run() error {
	return es.Configure()
}

func (es *es) Stop(ctx context.Context) <-chan bool {
	if ctx == nil {
		ctx = context.TODO()
	}

	c := make(chan bool)
	go func() {
		if es.client != nil {
			es.client.InstrumentationEnabled().Close(ctx)
		}
		c <- true
	}()
	return c
}

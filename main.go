package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gertjaap/dlcoracle/crypto"
	"github.com/gertjaap/dlcoracle/logging"
	"github.com/gertjaap/dlcoracle/publisher"
	"github.com/gertjaap/dlcoracle/routes"
	"github.com/gertjaap/dlcoracle/store"

	"github.com/awnumar/memguard"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"path/filepath"
	flags "github.com/jessevdk/go-flags"
	"github.com/gertjaap/dlcoracle/gcfg"
)


type dlcConfig struct { 
	DataDir      string `long:"DataDir" description:"Connect to bitcoin testnet3."`
	HttpPort     string `long:"HttpPort" description:"bc2 full node."`
	// For testing purposes
	Interval 	 uint64  `long:"Interval" description:"Interval in seconds."`
	ValueToPublish	uint64  `long:"ValueToPublish" description:"Value to publish"`
}


// newConfigParser returns a new command line flags parser.
func newConfigParser(conf *dlcConfig, options flags.Options) *flags.Parser {
	parser := flags.NewParser(conf, options)
	return parser
}

var (
	defaultDataDir  = "oracle"
	defaultHttpPort = "3000"
	// For testing purposes
	defaultInterval = uint64(300)
	defaultValueToPublish = uint64(11)
)

func main() {

	logging.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	conf := dlcConfig{
		DataDir:		defaultDataDir,
		HttpPort:		defaultHttpPort,
		// For testing purposes
		Interval:		defaultInterval,
		ValueToPublish:		defaultValueToPublish,
	}

	var err error

	preParser := newConfigParser(&conf, flags.HelpFlag)
	_, err = preParser.ParseArgs(os.Args) // parse the cli
	if err != nil {
		logging.Error.Fatal(err)
	}

	// create home directory
	_, err = os.Stat(conf.DataDir)
	if err != nil {
		logging.Info.Println("dlcoracle home directory does not exist.")
	}
	if os.IsNotExist(err) {
		// first time the guy is running dlcoracle
		os.Mkdir(conf.DataDir, 0700)
		logging.Info.Println("dlcoracle home directory have been created")

	}


	gcfg.DataDir = conf.DataDir
	gcfg.Interval = conf.Interval
	gcfg.ValueToPublish = conf.ValueToPublish


	logging.Info.Println("MIT-DCI Discreet Log Oracle starting...")

	key, err := crypto.ReadKeyFile(filepath.Join(conf.DataDir, "privkey.hex"))
	if err != nil {
		logging.Error.Fatal("Could not open or create keyfile:", err)
		os.Exit(1)
	}
	crypto.StoreKeys(key)
	// Tell memguard to listen out for interrupts, and cleanup in case of one.
	memguard.CatchInterrupt(func() {
		fmt.Println("Interrupt signal received. Exiting...")
	})
	// Make sure to destroy all LockedBuffers when returning.
	defer memguard.DestroyAll()

	store.Init()
	logging.Info.Println("Initialized store...")

	publisher.Init()
	logging.Info.Println("Started publisher...")

	r := mux.NewRouter()
	r.HandleFunc("/api/datasources", routes.ListDataSourcesHandler)
	r.HandleFunc("/api/datasource/{id}/value", routes.DataSourceValueHandler)
	r.HandleFunc("/api/pubkey", routes.PubKeyHandler)
	r.HandleFunc("/api/rpoint/{datasource}/{timestamp}", routes.RPointHandler)
	r.HandleFunc("/api/publication/{R}", routes.PublicationHandler)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static"))))

	// CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	logging.Info.Printf("Listening on port: %s", conf.HttpPort)

	logging.Error.Fatal(http.ListenAndServe(":" + conf.HttpPort, handlers.CORS(originsOk, headersOk, methodsOk)(logging.WebLoggingMiddleware(r))))
}

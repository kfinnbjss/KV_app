package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"store/logging"
	"store/server"
	"store/users"
)

const (
	ConnHost = "localhost"
)

func main() {
	//initialise loggers
	logging.SetupLoggers("../info.log", "../htaccess.log") //pass in the log files so they can be closed at the end of the main function
	defer logging.Shutdown()

	//read command line flags
	portPtr := flag.Int("port", 0, "Port number to run on")
	depthPtr := flag.Int("depth", 1000, "Maximum size of the KV store")
	bufferPtr := flag.Int("buffer", 100, "KV store buffer")

	flag.Parse()
	if *portPtr <= 0 { //Todo, distinguish between no port received and port set to 0
		logging.WarningLogger.Println("no port received")
		fmt.Println("Port is required")
		os.Exit(-1)
	}

	if *depthPtr <= 0 {
		logging.WarningLogger.Println("invalid store depth received", *depthPtr)
		fmt.Println("Invalid store depth")
		os.Exit(-1)
	}

	if *bufferPtr < 0 {
		logging.WarningLogger.Println("invalid store buffer received", *bufferPtr)
		fmt.Println("Invalid store buffer")
		os.Exit(-1)
	}

	//Fill the users database
	errUser := users.FillUserDB()
	if errUser != nil {
		logging.ErrorLogger.Println("Unable to fill the users database", errUser)
		os.Exit(-1)
	}

	errSetupServer := server.Setup(*portPtr, ConnHost, *bufferPtr, *depthPtr)
	if errSetupServer != nil {
		logging.ErrorLogger.Println("problem setting up server", errSetupServer)
	}

	err := server.Start()
	if err != nil {
		if err == http.ErrServerClosed {
			logging.InfoLogger.Println("server shut down")
		} else {
			logging.ErrorLogger.Println("problem starting server", err)
			fmt.Println("Problem starting server")
			os.Exit(-2)
		}
	}

	fmt.Println("Done")
}

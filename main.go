package main // import "hello-bbolt"

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blevesearch/bleve/v2"
)

func setupDB() (bleve.Index, *EncryptedDB, error) {
	var err error

	// Bleve
	index, err = bleve.Open(indexName)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(indexName, mapping)
		if err != nil {
			return nil, nil, err
		}
	} else if err != nil {
		return nil, nil, err
	}

	// BBolt
	db, err := OpenEncryptedDB(databaseName, 0600, nil)
	if err != nil {
		return nil, nil, err
	}

	return index, db, err
}

func startServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /init-data", dataInitHandler)

	mux.HandleFunc("POST /person", addPersonHandler)
	mux.HandleFunc("DELETE /person", deletePersonHandler)
	mux.HandleFunc("PUT /person", updatePersonHandler)
	mux.HandleFunc("GET /person", getPersonListHandler)

	mux.HandleFunc("POST /txt-file", setTextFileHandler)
	mux.HandleFunc("GET /txt-file", getTextFileHandler)

	server := &http.Server{Addr: listenADDR, Handler: mux}
	go func() {
		fmt.Println("Server starting on " + listenADDR)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("Server error: " + err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v", err)
	}

	if index != nil {
		index.Close()
	}
	if db != nil {
		db.Close()
	}

	fmt.Println("Server exited")
}

func main() {
	var err error

	index, db, err = setupDB()
	if err != nil {
		panic(err)
	}

	startServer()
}

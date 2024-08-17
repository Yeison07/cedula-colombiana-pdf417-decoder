package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/yeison07/cedula-colombiana-pdf417-decoder/server"
	"github.com/yeison07/cedula-colombiana-pdf417-decoder/utils"
)

var (
	mu         sync.Mutex
	cancelFunc context.CancelFunc
)

/*
Inicio del programa, configura los logger y inicia el servidor http
*/
func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		router.Get("/getdata", getCedulaDataHandler)
		router.Get("/cancel", cancelDataHandler)
		return router
	})

}

/*Handler para el endpoint getdata*/
func getCedulaDataHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	if cancelFunc != nil {
		cancelFunc()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	cancelFunc = cancel
	mu.Unlock()

	inputScanner, err := utils.StartSerialPort(ctx)
	if err != nil && err.Error() == "lectura del puerto serie cancelada" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cancelación exitosa"))
		return
	}
	if err != nil {
		http.Error(w, "Error al escanear el documento: \ninternal-error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	persona := utils.FormatterScannerInput(inputScanner)
	if persona == nil {
		http.Error(w, "Hubo un error al momento de procesar la cedula: \ninternal-error: ", http.StatusInternalServerError)
		return
	}
	personaJSON, err := json.Marshal(persona)
	if err != nil {
		log.Error().Err(err).Msg("Error al convertir los datos de la persona a JSON")
		http.Error(w, "Error al convertir los datos de la persona a JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(personaJSON)
}

/*Cancela la lectura*/
func cancelDataHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	if cancelFunc != nil {
		cancelFunc()
		cancelFunc = nil
	}
	mu.Unlock()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cancelación exitosa"))
}

package utils

import (
	"bytes"
	"context"
	"errors"
	"runtime"
	"sync"

	"github.com/rs/zerolog/log"
	"go.bug.st/serial"
)

var (
	portMu sync.Mutex
)

func StartSerialPort(ctx context.Context) (string, error) {
	portMu.Lock()
	defer portMu.Unlock()

	var devicePortDetected string

	if runtime.GOOS == "windows" {
		devicePortDetected = "COM3"
		log.Info().Msgf("Puerto serie detectado: %s", devicePortDetected)
	} else {
		ports, err := detectSerialPorts()
		if err != nil {
			log.Error().Err(err).Msg("Error al detectar los puertos serie")
			return "", err
		}
		devicePortDetected = ports[0]
	}

	port, err := openSerialPort(devicePortDetected)
	if err != nil {
		log.Error().Err(err).Msg("Error al abrir el puerto serie")
		return "", err
	}
	defer func() {
		log.Info().Msg("Cerrando el puerto serie")
		port.Close()
	}()

	var resultBuffer bytes.Buffer
	buff := make([]byte, 100)
	scanSize := 150
	bytesRead := 0

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("Lectura del puerto serie cancelada o timeout")
				return
			default:
				n, err := port.Read(buff)
				if err != nil {
					log.Error().Err(err).Msg("Error al leer datos del puerto serie")
					return
				}
				if n == 0 {
					log.Info().Msg("Fin de la lectura")
					return
				}
				resultBuffer.Write(buff[:n])
				bytesRead += n
				if bytesRead >= scanSize {
					log.Info().Msg("Escaneo completado")
					return
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			port.Close()
			log.Error().Msg("Timeout de lectura del puerto serie")
			return "", errors.New("timeout de lectura del puerto serie")
		}
		return "", errors.New("lectura del puerto serie cancelada")
	case <-done:
		log.Info().Msg("Lectura del puerto serie completada")
		if resultBuffer.Len() <= 0 {
			log.Error().Err(err).Msg("La lectura se completó pero no hay datos, revisa la conexión USB")
			return "", errors.New("la lectura se completó pero no hay datos")
		}
		return resultBuffer.String(), nil
	}

}

func detectSerialPorts() ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, err
	}
	if len(ports) == 0 {
		return nil, errors.New("no se encontraron puertos serie")
	}
	for _, port := range ports {
		log.Info().Msgf("Puerto serie detectado: %s", port)
		return ports, nil
	}
	return nil, errors.New("no se encontraron puertos serie (Hubo un problema inesperado)")
}

func openSerialPort(devicePort string) (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open(devicePort, mode)
	if err != nil {
		return nil, err
	}
	return port, nil
}

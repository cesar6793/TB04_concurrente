package main

import (
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// Estructura para los puntos y centroides
type DataPoint struct {
	Dimensions []float64
}

// Función principal del cliente
func client(points []DataPoint) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error al conectar con el servidor:", err)
		return
	}
	defer conn.Close()

	// Codificar y enviar los puntos al servidor
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(points)
	if err != nil {
		fmt.Println("Error al codificar:", err)
		return
	}

	// Recibir y decodificar los centroides resultantes del servidor
	decoder := gob.NewDecoder(conn)
	var centroids []DataPoint
	err = decoder.Decode(&centroids)
	if err != nil {
		fmt.Println("Error al decodificar:", err)
		return
	}

	// Mostrar los centroides resultantes
	fmt.Println("Centroides recibidos del servidor:")
	for _, centroid := range centroids {
		fmt.Println(centroid.Dimensions)
	}
}

func main() {
	url := "https://raw.githubusercontent.com/cesar6793/concurrente/main/dataset.csv"

	// Realizar la solicitud GET a la URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error al realizar la solicitud GET:", err)
		return
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error al leer el CSV:", err)
		return
	}

	// Convertir los registros CSV en puntos de datos
	var points []DataPoint
	for _, record := range records {
		var dimensions []float64
		for _, value := range record {
			number, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
			if err != nil {
				fmt.Println("Error al convertir el valor a float:", err)
				return
			}
			dimensions = append(dimensions, number)
		}
		points = append(points, DataPoint{Dimensions: dimensions})
	}

	// Llamar a la función cliente con los puntos cargados
	client(points)
}

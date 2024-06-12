package main

import (
	"encoding/gob"
	"fmt"
	"math"
	"net"
	"sync"
)

// Estructura para los puntos y centroides
type DataPoint struct {
	Dimensions []float64
}

// Función para calcular la distancia euclidiana
func distance(a, b DataPoint) float64 {
	var sum float64
	for i := range a.Dimensions {
		sum += math.Pow(a.Dimensions[i]-b.Dimensions[i], 2)
	}
	return math.Sqrt(sum)
}

// Función para encontrar el centroide más cercano
func closestCentroid(point DataPoint, centroids []DataPoint) int {
	minDist := math.MaxFloat64
	closestIdx := -1
	for i, centroid := range centroids {
		dist := distance(point, centroid)
		if dist < minDist {
			minDist = dist
			closestIdx = i
		}
	}
	return closestIdx
}

// Función para actualizar los centroides
func updateCentroids(points []DataPoint, assignments []int, k int) []DataPoint {
	centroids := make([]DataPoint, k)
	counts := make([]int, k)

	for i := range centroids {
		centroids[i].Dimensions = make([]float64, len(points[0].Dimensions))
	}

	for i, point := range points {
		cluster := assignments[i]
		if cluster >= k {
			fmt.Println("Error: Asignación de cluster inválida.")
			return nil
		}
		for j := range point.Dimensions {
			centroids[cluster].Dimensions[j] += point.Dimensions[j]
		}
		counts[cluster]++
	}

	for i := range centroids {
		for j := range centroids[i].Dimensions {
			if counts[i] > 0 {
				centroids[i].Dimensions[j] /= float64(counts[i])
			}
		}
	}

	return centroids
}

// Función para manejar cada conexión de cliente
func handleClient(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()

	decoder := gob.NewDecoder(conn)
	var points []DataPoint
	err := decoder.Decode(&points)
	if err != nil {
		fmt.Println("Error al decodificar:", err)
		return
	}

	k := 4 // Número de clusters

	if len(points) < k {
		fmt.Println("Error: No hay suficientes puntos para formar los clusters solicitados.")
		return
	}
	centroids := make([]DataPoint, k)
	for i := range centroids {
		centroids[i] = points[i] // Inicialización simple para el ejemplo
	}

	assignments := make([]int, len(points))
	for i := range points {
		assignments[i] = closestCentroid(points[i], centroids)
	}

	centroids = updateCentroids(points, assignments, k)

	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(centroids)
	if err != nil {
		fmt.Println("Error al codificar:", err)
		return
	}
}

// Función principal del servidor
func server() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
		return
	}
	defer ln.Close()

	var wg sync.WaitGroup

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error al aceptar conexión:", err)
			continue
		}

		wg.Add(1)
		go handleClient(conn, &wg)
	}

	wg.Wait()
}

func main() {
	server()
}

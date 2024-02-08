package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
)

// STLCalc represents a calculator for STL files

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/calculate", calculateHandler)
	fmt.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}

}
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("stl")
	if err != nil {
		http.Error(w, "Error retrieving the STL file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a temporary file to save the uploaded STL
	tempFile, err := ioutil.TempFile("temp", "*.stl")
	if err != nil {
		http.Error(w, "Error creating a temporary file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	// Copy the uploaded STL content to the temporary file
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading the STL file", http.StatusInternalServerError)
		return
	}
	tempFile.Write(fileBytes)
	basePrice := r.URL.Query().Get("basePrice")
	if basePrice == "" {
		http.Error(w, "basePrice query parameter is required", http.StatusBadRequest)
		return
	}

	if convertprice, err := strconv.ParseFloat(basePrice, 64); err == nil {

		// Process the STL file
		response, err := processSTL(tempFile.Name(), convertprice)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the volume and weight
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is accepted", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve the filepath from the query string
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "Path query parameter is required", http.StatusBadRequest)
		return
	}
	basePrice := r.URL.Query().Get("basePrice")
	if basePrice == "" {
		http.Error(w, "basePrice query parameter is required", http.StatusBadRequest)
		return
	}

	if convertprice, err := strconv.ParseFloat(basePrice, 64); err == nil {

		// Process the STL file
		response, err := processSTL(filePath, convertprice)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the volume and weight
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func processSTL(filePath string, basePrice float64) (map[string]interface{}, error) {
	stlCalc, err := NewSTLCalc(filePath)
	if err != nil {
		return nil, err
	}
	defer stlCalc.Close()
	//defer os.Remove(filePath) // Remove the temp file after processing

	volume, err := stlCalc.GetVolume("cm")
	if err != nil {
		return nil, err
	}

	weight, err := stlCalc.GetWeight()
	if err != nil {
		return nil, err
	}
	var price = volume * basePrice
	roundedValue := math.Round(price)

	// Convert to integer
	completePrice := int(roundedValue)
	return map[string]interface{}{
		"volume": volume,
		"weight": weight,
		"price":  completePrice,
	}, nil
}

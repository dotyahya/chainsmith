package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"net/http"
	"os"
	"os/exec"
)

// fetchFromIPFS fetches files from Kubo/IPFS using their CID and saves them locally with the correct extension
func fetchFromIPFS(cid, extension string) (string, error) {
	// Using the "ipfs cat" command to fetch the file locally
	cmd := exec.Command("ipfs", "cat", cid)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to fetch file from IPFS with CID %s: %v", cid, err)
	}

	// Save the output to a file with the given extension
	filePath := fmt.Sprintf("./execution/%s%s", cid, extension) // Save the file locally with the CID and extension
	err = os.WriteFile(filePath, output, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to save file locally: %v", err)
	}
	return filePath, nil
}

type Selection struct {
	DatasetCid string `json:"datasetCid"`
	Algorithm  string `json:"algorithm"`
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var selection Selection
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&selection)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Fetch the dataset file from IPFS with the ".csv" extension
		datasetFilePath, err := fetchFromIPFS(selection.DatasetCid, ".csv")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching dataset: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Printf("Dataset fetched from IPFS: %s\n", datasetFilePath)

		// Fetch the algorithm file from IPFS with the ".go" extension
		algoFilePath, err := fetchFromIPFS(selection.Algorithm, ".go")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching algorithm: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Printf("Algorithm fetched from IPFS: %s\n", algoFilePath)

		// Return a response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Dataset and Algorithm selected and fetched successfully",
		})
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/process", processHandler)

	// CORS handling
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),                                       // Allow all origins (use specific URLs in production)
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), // Allowed methods
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),           // Allowed headers
	)(http.DefaultServeMux)

	// Start server with CORS enabled
	fmt.Println("Starting server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

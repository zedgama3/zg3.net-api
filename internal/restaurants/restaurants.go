package restaurants

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// RestaurantSearch is a handler for the /restaurant/search endpoint.
func AutoComplete(w http.ResponseWriter, r *http.Request) {
	// Check that needed components are present
	var db *sql.DB
	if newDb, exists := r.Context().Value("db").(*sql.DB); !exists {
		http.Error(w, "DB Handle not found", http.StatusInternalServerError)
		return
	} else {
		db = newDb
	}

	// Get the name parameter from the URL
	params := mux.Vars(r)
	name := params["name"]
	//jsonOutput(w, name, http.StatusOK)
	//return

	restaurants, err := searchDbForRestaurant(db, name)
	if err != nil {
		jsonError(w, "Error searching for restaurant", http.StatusInternalServerError)
		return
	}
	if restaurants == nil {
		jsonError(w, "No restaurants found", http.StatusNotFound)
		return
	}
	jsonOutput(w, restaurants, http.StatusOK)
	return

}

func searchDbForRestaurant(db *sql.DB, name string) ([][2]string, error) {
	rows, err := db.Query("SELECT name, restaurant_id FROM restaurants WHERE name ILIKE '%' || $1 || '%'", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants [][2]string
	for rows.Next() {
		var restaurant [2]string
		err := rows.Scan(&restaurant[0], &restaurant[1])
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}

/************************************************************/

type ErrorResponse struct {
	Error string `json:"error"`
}

func jsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func jsonOutput(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

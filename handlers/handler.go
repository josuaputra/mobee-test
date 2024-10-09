package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// Get all items
func GetItems(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var items []Item

		rows, err := db.Query(context.Background(), "SELECT id, name, price FROM items")
		if err != nil {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var item Item
			if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
				http.Error(w, "Error reading row", http.StatusInternalServerError)
				return
			}
			items = append(items, item)
		}

		json.NewEncoder(w).Encode(items)
	}
}

// Create an item
func CreateItem(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item Item
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(context.Background(), "INSERT INTO items (name, price) VALUES ($1, $2)", item.Name, item.Price)
		if err != nil {
			http.Error(w, "Error inserting item", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(item)
	}
}

// Update an item
func UpdateItem(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		var item Item
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(context.Background(), "UPDATE items SET name = $1, price = $2 WHERE id = $3", item.Name, item.Price, id)
		if err != nil {
			http.Error(w, "Error updating item", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(item)
	}
}

// Delete an item
func DeleteItem(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		_, err := db.Exec(context.Background(), "DELETE FROM items WHERE id = $1", id)
		if err != nil {
			http.Error(w, "Error deleting item", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Item deleted successfully")
	}
}

func Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	}
}

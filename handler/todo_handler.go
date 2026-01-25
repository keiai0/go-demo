package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/keiai0/go-demo/model"
	"gorm.io/gorm"
)

func ListTodos(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todos []model.Todo

		if err := db.Find(&todos).Error; err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to retrieve todos", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(todos); err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to encode todos", http.StatusInternalServerError)
			return
		}
	}
}

func UpdateTodo(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 特定の ID のtodoを探す
		var existingTodo model.Todo
		if err := db.First(&existingTodo, 1).Error; err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to retrieve todo", http.StatusInternalServerError)
			return
		}

		// 更新するフィールドを設定
		existingTodo.Title = "更新後のタイトル"
		existingTodo.Note = "更新後の備考"
		existingTodo.UpdatedAt = time.Now()

		// 更新処理
		if err := db.Updates(&existingTodo).Error; err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to update todo", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(existingTodo); err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to encode todo", http.StatusInternalServerError)
			return
		}
	}
}

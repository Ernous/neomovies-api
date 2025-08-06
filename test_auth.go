package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	fmt.Println("=== AUTH TEST ===")
	
	// Тестируем аутентификацию
	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "test123",
	}
	
	jsonData, _ := json.Marshal(loginData)
	
	resp, err := http.Post("http://localhost:3000/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("📊 Status: %d\n", resp.StatusCode)
	fmt.Printf("📄 Response: %s\n", string(body))
	
	fmt.Println("=== END AUTH TEST ===")
}
package main

type jsonResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
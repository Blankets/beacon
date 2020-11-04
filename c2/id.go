package main

import "github.com/google/uuid"

func getNewID() string {
	id := uuid.New()
	return id.String()
}

package main

import "github.com/google/uuid"

func getNewId() string {
	id := uuid.New()
	return id.String()
}

package id

import "github.com/google/uuid"

func GenerateID() string {
	return uuid.NewString()
}

// func generateRoomID() string {
// 	return uuid.NewString()[:8]
// }

// func generateProverbID() string {
// 	return uuid.NewString()[:4]
// }

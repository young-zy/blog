package services

import (
	"log"
	"testing"
)

func TestRegister(t *testing.T) {
	err := Register("admin", "admin", "test@young-zy.com")
	if err != nil {
		log.Println(err)
	}
}

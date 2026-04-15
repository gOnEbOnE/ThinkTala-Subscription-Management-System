package auth

import (
	"log"

	"github.com/master-abror/zaframework/core/utils"
)

// InitRedis initializes Redis used by gateway auth/token flow.
func InitRedis() {
	if err := utils.InitRedis(); err != nil {
		log.Printf("[GW][AUTH] Redis init failed: %v", err)
	}
}

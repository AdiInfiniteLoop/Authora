package handlers

import "github.com/AdiInfiniteLoop/Authora/internal/config"

type LocalApiConfig struct {
	//Composition relationship ("has a" relationship)
	*config.ApiConfig
}

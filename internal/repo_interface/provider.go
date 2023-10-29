package repo_interface

import "github.com/Monstergogo/beauty-share/internal/repo"

func MongoRepoProvider() MongoRepo {
	return &repo.MongoRepoImpl{}
}

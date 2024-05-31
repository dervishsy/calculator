package healthz

import (
	"calculator/internal/shared/entities"
	"calculator/pkg/logger"
	"calculator/pkg/utils"
	"net/http"
)

type response struct {
	Name         string `json:"name"`
	BuildVersion string `json:"build_version"`
	BuildTime    string `json:"build_time"`
	GitTag       string `json:"git_tag"`
	GitHash      string `json:"git_hash"`
}

func MakeHandler(info *entities.AppInfo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		response := &response{
			Name:         info.Name,
			BuildVersion: info.BuildVersion,
			BuildTime:    info.BuildTime,
			GitTag:       info.GitTag,
			GitHash:      info.GitHash,
		}
		err := utils.SuccessRespondWith200(w, response)
		if err != nil {
			logger.Error("failed to decode response", response, err)
		}
	}
}

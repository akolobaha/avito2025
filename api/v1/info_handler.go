package httphandler

import (
	"avito2015/internal/info"
	"avito2015/internal/user"
	"avito2015/pkg/jsonresponse"
	"errors"
	"net/http"
)

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	usr, ok := r.Context().Value("user").(*user.User)
	if !ok {
		jsonresponse.Error(w, errors.New("unauthorized"), 500)
		return
	}

	repo := info.NewInfoRepository()
	s := info.NewInfoService(repo)

	infoR, err := s.Get(usr)
	if err != nil {
		jsonresponse.Error(w, err, 500)
		return
	}

	jsonresponse.StatusOK(w, infoR)

}

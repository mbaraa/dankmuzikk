package main

import (
	"dankmuzikk/config"
	"dankmuzikk/handlers/middlewares/logger"
	"dankmuzikk/log"
	"net/http"
)

func main() {
	err := StartServer()
	if err != nil {
		log.Fatalln(err)
	}
}

func StartServer() error {
	// mariadbRepo, err := mariadb.New()
	// if err != nil {
	// return err
	// }
	// app := app.New(mariadbRepo)
	// jwtUtil := jwt.New[actions.TokenPayload]()
	//	usecases := actions.New(
	//		app,
	//		nil,
	//		jwtUtil,
	//		nil,
	//		nil,
	//	)
	// authMw := auth.New(usecases)
	// applicationHandler.Handle("/muzikkx/", authMw.AuthHandler(http.StripPrefix("/muzikkx", http.FileServer(http.Dir(config.Env().YouTube.MuzikkDir)))))
	applicationHandler := http.NewServeMux()

	muzikkxDir := config.Env().YouTube.MuzikkDir + "/muzikkx/"
	pixDir := config.Env().YouTube.MuzikkDir + "/pix/"

	applicationHandler.Handle("/muzikkx/", http.StripPrefix("/muzikkx", http.FileServer(http.Dir(muzikkxDir))))
	applicationHandler.Handle("/muzikkx-raw/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment")
		http.
			StripPrefix("/muzikkx-raw", http.FileServer(http.Dir(muzikkxDir))).
			ServeHTTP(w, r)
	}))
	applicationHandler.Handle("/pix/", http.StripPrefix("/pix", http.FileServer(http.Dir(pixDir))))

	log.Info("Starting http cdn server at port " + config.Env().CdnPort)
	if config.Env().GoEnv == "dev" || config.Env().GoEnv == "beta" {
		return http.ListenAndServe(":"+config.Env().CdnPort, logger.Handler(applicationHandler))
	}
	return http.ListenAndServe(":"+config.Env().CdnPort, applicationHandler)
}

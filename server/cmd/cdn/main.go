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
	applicationHandler.Handle("/muzikkx/", http.StripPrefix("/muzikkx", http.FileServer(http.Dir(config.Env().YouTube.MuzikkDir))))
	applicationHandler.Handle("/muzikkx-raw/", http.StripPrefix("/muzikkx-raw", http.FileServer(http.Dir(config.Env().YouTube.MuzikkDir))))

	log.Info("Starting http cdn server at port " + config.Env().CdnPort)
	if config.Env().GoEnv == "dev" || config.Env().GoEnv == "beta" {
		return http.ListenAndServe(":"+config.Env().CdnPort, logger.Handler(applicationHandler))
	}
	return http.ListenAndServe(":"+config.Env().CdnPort, applicationHandler)
}

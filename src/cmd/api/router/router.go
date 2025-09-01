package router

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rishirishhh/vought/src/cmd/api/config"
	"github.com/rishirishhh/vought/src/cmd/api/controllers"
	"github.com/rishirishhh/vought/src/cmd/api/db/dao"
	"github.com/rishirishhh/vought/src/pkg/clients"
)

type Clients struct {
	S3Client              clients.IS3Client
	AmqpClient            clients.AmqpClient
	AmqpVideoStatusUpdate clients.AmqpClient
	UUIDGen               clients.IUUIDGenerator
}

type DAOs struct {
	Db         *sql.DB
	VideosDAO  dao.VideosDAO
	UploadsDAO dao.UploadsDAO
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewRouter(config config.Config, clients *Clients, DAOs *DAOs) http.Handler {
	r := mux.NewRouter()
	r.PathPrefix("/ws").Handler(controllers.WSHandler{Config: config, AmqpVideoStatusUpdate: clients.AmqpVideoStatusUpdate}).Methods("GET")

	v1 := r.PathPrefix("/api/v1").Subrouter()
	v1.PathPrefix("/videos/upload").Handler(controllers.VideoUploadHandler{S3Client: clients.S3Client, AmqpClient: clients.AmqpClient, AmqpVideoStatusUpdate: clients.AmqpVideoStatusUpdate, VideosDAO: &DAOs.VideosDAO, UploadsDAO: &DAOs.UploadsDAO, UUIDGen: clients.UUIDGen}).Methods("POST")
	return handlers.CORS(getCORS())(r)
}

func getCORS() (handlers.CORSOption, handlers.CORSOption, handlers.CORSOption, handlers.CORSOption) {
	corsObj := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS", "DELETE"})
	headers := handlers.AllowedHeaders([]string{"Authorization"})
	credentials := handlers.AllowCredentials()

	return corsObj, methods, headers, credentials
}

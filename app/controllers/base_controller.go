package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gorilla/sessions"

	"github.com/gieart87/gotoko/app/models"
	"github.com/gieart87/gotoko/database/seeders"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB        *gorm.DB
	Router    *mux.Router
	AppConfig *AppConfig
}

type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string
	AppURL  string
}

type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBDriver   string
}

type PageLink struct {
	Page          int32
	Url           string
	IsCurrentPage bool
}

type PaginationLinks struct {
	CurrentPage string
	NextPage    string
	PrevPage    string
	TotalRows   int32
	TotalPages  int32
	Links       []PageLink
}

type PaginationParams struct {
	Path        string
	TotalRows   int32
	PerPage     int32
	CurrentPage int32
}

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var sessionShoppingCart = "shopping-cart-session"
var sessionFlash = "flash-session"

func (server *Server) Initialize(appConfig AppConfig, dbConfig DBConfig) {
	fmt.Println("Welcome to " + appConfig.AppName)

	server.initializeDB(dbConfig)
	server.initializeAppConfig(appConfig)
	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func (server *Server) initializeDB(dbConfig DBConfig) {
	var err error
	if dbConfig.DBDriver == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)
		server.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
		server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		panic("Failed on connecting to the database server")
	}
}

func (server *Server) initializeAppConfig(appConfig AppConfig) {
	server.AppConfig = &appConfig
}

func (server *Server) dbMigrate() {
	for _, model := range models.RegisterModels() {
		err := server.DB.Debug().AutoMigrate(model.Model)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Database migrated successfully.")
}

func (server *Server) InitCommands(config AppConfig, dbConfig DBConfig) {
	server.initializeDB(dbConfig)

	cmdApp := cli.NewApp()
	cmdApp.Commands = []cli.Command{
		{
			Name: "db:migrate",
			Action: func(c *cli.Context) error {
				server.dbMigrate()
				return nil
			},
		},
		{
			Name: "db:seed",
			Action: func(c *cli.Context) error {
				err := seeders.DBSeed(server.DB)
				if err != nil {
					log.Fatal(err)
				}

				return nil
			},
		},
	}

	err := cmdApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func GetPaginationLinks(config *AppConfig, params PaginationParams) (PaginationLinks, error) {
	var links []PageLink

	totalPages := int32(math.Ceil(float64(params.TotalRows) / float64(params.PerPage)))

	for i := 1; int32(i) <= totalPages; i++ {
		links = append(links, PageLink{
			Page:          int32(i),
			Url:           fmt.Sprintf("%s/%s?page=%s", config.AppURL, params.Path, fmt.Sprint(i)),
			IsCurrentPage: int32(i) == params.CurrentPage,
		})
	}

	var nextPage int32
	var prevPage int32

	prevPage = 1
	nextPage = totalPages

	if params.CurrentPage > 2 {
		prevPage = params.CurrentPage - 1
	}

	if params.CurrentPage < totalPages {
		nextPage = params.CurrentPage + 1
	}

	return PaginationLinks{
		CurrentPage: fmt.Sprintf("%s/%s?page=%s", config.AppURL, params.Path, fmt.Sprint(params.CurrentPage)),
		NextPage:    fmt.Sprintf("%s/%s?page=%s", config.AppURL, params.Path, fmt.Sprint(nextPage)),
		PrevPage:    fmt.Sprintf("%s/%s?page=%s", config.AppURL, params.Path, fmt.Sprint(prevPage)),
		TotalRows:   params.TotalRows,
		TotalPages:  totalPages,
		Links:       links,
	}, nil
}

func (server *Server) GetProvinces() ([]models.Province, error) {
	response, err := http.Get(os.Getenv("API_ONGKIR_BASE_URL") + "province?key=" + os.Getenv("API_ONGKIR_KEY"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	provinceResponse := models.ProvinceResponse{}
	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return nil, readErr
	}

	jsonErr := json.Unmarshal(body, &provinceResponse)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return provinceResponse.ProvinceData.Results, nil

}

func (server *Server) GetCitiesByProvinceID(provinceID string) ([]models.City, error) {
	response, err := http.Get(os.Getenv("API_ONGKIR_BASE_URL") + "city?key=" + os.Getenv("API_ONGKIR_KEY") + "&province=" + provinceID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	cityResponse := models.CityResponse{}

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return nil, readErr
	}

	jsonErr := json.Unmarshal(body, &cityResponse)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return cityResponse.CityData.Results, nil
}

func (server *Server) CalculateShippingFee(shippingParams models.ShippingFeeParams) ([]models.ShippingFeeOption, error) {
	if shippingParams.Origin == "" || shippingParams.Destination == "" || shippingParams.Weight <= 0 || shippingParams.Courier == "" {
		return nil, errors.New("invalid params")
	}

	params := url.Values{}
	params.Add("key", os.Getenv("API_ONGKIR_KEY"))
	params.Add("origin", shippingParams.Origin)
	params.Add("destination", shippingParams.Destination)
	params.Add("weight", strconv.Itoa(shippingParams.Weight))
	params.Add("courier", shippingParams.Courier)

	response, err := http.PostForm(os.Getenv("API_ONGKIR_BASE_URL")+"cost", params)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	ongkirResponse := models.OngkirResponse{}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	jsonErr := json.Unmarshal(body, &ongkirResponse)
	if jsonErr != nil {
		return nil, jsonErr
	}

	var shippingFeeOptions []models.ShippingFeeOption
	for _, result := range ongkirResponse.OngkirData.Results {
		for _, cost := range result.Costs {
			shippingFeeOptions = append(shippingFeeOptions, models.ShippingFeeOption{
				Service: cost.Service,
				Fee:     cost.Cost[0].Value,
			})
		}
	}

	return shippingFeeOptions, nil
}

func SetFlash(w http.ResponseWriter, r *http.Request, name string, value string) {
	session, err := store.Get(r, sessionFlash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.AddFlash(value, name)
	session.Save(r, w)
}

func GetFlash(w http.ResponseWriter, r *http.Request, name string) []string {
	session, err := store.Get(r, sessionFlash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	fm := session.Flashes(name)
	if len(fm) < 0 {
		return nil
	}

	session.Save(r, w)
	var flashes []string
	for _, fl := range fm {
		flashes = append(flashes, fl.(string))
	}

	return flashes
}

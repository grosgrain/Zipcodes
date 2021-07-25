package freightMatchingService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	Redis "github.com/go-redis/redis"
	ZipCodeService "github.com/grosgrain/zipcodes/src/zipCodesService"
	_ "github.com/lib/pq"
)

type RequestService struct {
	Error error
}

type Carrier struct {
	Id string `json:"id"`
	Name string `json:"carriername"`
	Zip string `json:"zip"`
	Active int `json:"active"`
	NoLoad int `json:"noload"`
	PrimaryBrokerId *string `json:"primarybrokerid"`
	PowerUnit int `json:"powerunits"`
	Distance float32`json:"distance"`
}

var (
	redisClient *Redis.Client
	DB *sql.DB
	zipCodeService = ZipCodeService.NewRequestService()
)

const googleMaximumZipsPerRequest  = 20

func NewRequestService() *RequestService {
	s := new(RequestService)
	return s
}

func (s *RequestService)GetNearbyCarriers(zip string, radius float64, page int) ([]Carrier, error){
	connectRedis()
	val, err := redisClient.Get(zip + fmt.Sprintf("%f", radius)).Result()
	if err != nil {
		fmt.Println(err)
	}
	if len(val) > 0 {
		mapJSON := make([]Carrier, len(val))
		err = json.Unmarshal([]byte(val), &mapJSON)
		if err != nil {
			fmt.Println(err)
		}
		return mapJSON, nil
	} else {
		//Get zips within radius
		c := make(chan []string, 20)
		go func() {
			data, _ := zipCodeService.GetZipCodesWithinRadius(zip, radius, true)
			c <- data
		}()
		zipLists := <-c

		//Get carriers info from DB
		results := make(chan []Carrier)

		go func() {
			zipString := strings.Join(zipLists, "','")
			DB = connectDB()
			sql :=  fmt.Sprintf(`SELECT carriers.id, carriers.zip, carriers.carriername, carriers.active, carriers.noload, carriers.primarybrokerid, carriersaferinfo.powerunits FROM carriers JOIN carriersaferinfo ON carriersaferinfo.carrierid=carriers.id WHERE carriers.zip IN ('%s')`, zipString)
			rows, err := DB.Query(sql)
			if err != nil {
				panic(err)
			}
			defer rows.Close()
			defer DB.Close()

			var tempArray []Carrier
			for rows.Next() {
				var temp Carrier
				err = rows.Scan(&temp.Id, &temp.Zip, &temp.Name, &temp.Active, &temp.NoLoad, &temp.PrimaryBrokerId, &temp.PowerUnit)
				if err != nil {
					panic(err)
				}
				tempArray = append(tempArray, temp)
			}
			results <- tempArray
		}()

		carrierResults := <-results

		//Calculate driving distances between zips

		if len(carrierResults) > 0 {
			var lists []ZipCodeService.ZipToDistanceMapping
			var wg sync.WaitGroup
			for i := 0; i < len(zipLists)/googleMaximumZipsPerRequest+1; i++ {
				var end int
				if i*googleMaximumZipsPerRequest+googleMaximumZipsPerRequest < len(zipLists) {
					end = i*googleMaximumZipsPerRequest + googleMaximumZipsPerRequest
				} else {
					end = len(zipLists)
				}
				temp := zipLists[i*googleMaximumZipsPerRequest : end]
				wg.Add(1)
				go func() {
					data, _ := zipCodeService.GetDistanceFromOnePointToMultiplePoints(zip, temp, &wg)
					lists = append(lists, data...)
				}()
			}
			wg.Wait()
			zipMap := make(map[string]float32)
			for _, v := range lists {
				zipMap[v.Zip] = v.Distance
			}
			for i := len(carrierResults) - 1; i >= 0; i-- {
				if val, ok := zipMap[carrierResults[i].Zip]; ok && float64(val) <= radius {
					carrierResults[i].Distance = val
				} else {
					//Wipe out carriers that driving distance larger than radius
					carrierResults = append(carrierResults[:i], carrierResults[i + 1 :]...)
				}
			}

			//Add results to Redis cache
			json, _ := json.Marshal(carrierResults)
			err = redisClient.Set(zip + fmt.Sprintf("%f", radius), json, 86400).Err()
			if err != nil {
				fmt.Println(err)
			}
		}
		return carrierResults, nil
	}
}

func connectRedis() {
	redisClient = Redis.NewClient(&Redis.Options{
		Addr: "localhost:6379",
		Password: "dev",
		DB: 0,
	})
}

func connectDB() *sql.DB {
	host, _ := os.LookupEnv("DB_HOST")
	port, _ := os.LookupEnv("DB_PORT")
	name,_ := os.LookupEnv("DB_NAME")
	user,_ := os.LookupEnv("DB_USER")
	password,_ := os.LookupEnv("DB_PASSWORD")
	dbDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)

	DB, err := sql.Open("postgres", dbDSN)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	return DB
}

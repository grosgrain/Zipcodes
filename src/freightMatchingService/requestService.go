package freightMatchingService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	PrimaryBrokerId string `json:"primarybrokerid"`
	PowerUnit int `json:"powerunits"`
	Distance float64 `json:"distance"`
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
		/*results := make([]Carrier, 2)
		results[0] = Carrier{
			Id:              "1",
			Name:            "haha",
			Zip:             "45140",
			Active:          1,
			NoLoad:          1,
			PrimaryBrokerId: "",
			PowerUnit:       42,
			Distance:        0,
		}
		results[1] = Carrier{
			Id:              "1",
			Name:            "jaja",
			Zip:             "93065",
			Active:          1,
			NoLoad:          1,
			PrimaryBrokerId: "",
			PowerUnit:       14,
			Distance:        0,
		}
		json, _ := json.Marshal(results)
		err = redisClient.Set(zip + fmt.Sprintf("%f", radius), json, 0).Err()
		if err != nil {
			fmt.Println(err)
		}
		return results, nil*/
		c := make(chan []string, 20)
		go func() {
			data, _ := zipCodeService.GetZipCodesWithinRadius(zip, radius, false)
			c <- data
		}()
		zipLists := <-c

		DB = connectDB()
		sql := "SELECT id, zip, carriername, active, noload, primarybrokerid, powerunit FROM carriers JOIN carriersaferinfo ON carriersaferinfo.carrierid=carriers.id WHERE zip IN ($1) DISTINCT"
		rows, err := DB.Query(sql, zipLists)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		var results []Carrier
		for rows.Next() {
			var temp Carrier
			err = rows.Scan(&temp.Id, &temp.Zip, &temp.Name, &temp.Active, &temp.NoLoad, &temp.PrimaryBrokerId, &temp.PowerUnit)
			if err != nil {
				panic(err)
			}
			results = append(results, temp)
		}
		var lists []ZipCodeService.ZipToDistanceMapping
		var wg sync.WaitGroup
		for i := 0; i < len(zipLists) / googleMaximumZipsPerRequest + 1; i++ {
			var end int
			if i * googleMaximumZipsPerRequest + googleMaximumZipsPerRequest < len(zipLists) {
				end = i * googleMaximumZipsPerRequest + googleMaximumZipsPerRequest
			} else {
				end = len(zipLists)
			}
			temp := zipLists[i * googleMaximumZipsPerRequest : end]
			wg.Add(1)
			go func() {
				data,_ := zipCodeService.GetDistanceFromOnePointToMultiplePoints(zip, temp, &wg)
				lists = append(lists, data...)
			}()
		}
		wg.Wait()
		res := make(map[string]float32)
		for _, v:= range lists {
			if float64(v.Distance) <= radius {
				res[v.Zip] = v.Distance
			}
		}
		return results, nil
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

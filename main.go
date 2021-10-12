package main

import (
	"WebCrawler/webCrawler"
	//"WebCrawler/webCrawler"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"log"
	"strconv"
)

/*
	TODO - SAVE ALL REQUIRED DATA TO DATABASE(IN PROGRESS)
	TODO -  MOVIE DATA FROM API AND SAVING TO DATABASE
	MovieData Response Format :
	adult	false -> Bool
	backdrop_path	null -> URL
	genre_ids	[] -> a list of genre ids  -> use another table to save those information
	id	883725 -> primary key of movie id
	original_language	"en" -> string
	original_title	"Motions" -> string
	overview	""
	popularity	0
	poster_path	null
	release_date	"2021-04-06"
	title	"Motions"
	video	false
	vote_average	0
	vote_count	0

	TODO -  GENRE DATA FROM API -/genre/movie/list (Genre table -> one to many movie info)
	Genre Data Format :
	id - int
	name - string

	TODO -  ACTOR DATA FROM API AND SAVING TO DATABASE
	-PEOPLE RESPONSE DATA
	adult	- bool
	gender	- int
	id	- int
	known_for - [KnownFor]
	known_for_department- string
	name - string
	popularity	- float
	profile_path	- string

	-PEOPLE KnownFor RESPONSE
	backdrop_path - string
	first_air_date	- string
	genre_ids - []int
	id	- int
	media_type	- string
	name - string
	origin_country - []string
	original_language	- string
	original_name	- string
	overview	- string
	poster_path	- string
	vote_average	- float
	vote_count	- int

*/

const (
	host string = "https://api.themoviedb.org/3"
	apiKey string = "29570e7acc52b3e085ab46f6a60f0a55"
	upcomingURI = "/movie/upcoming"
	allMovieURI string ="/discover/movie"
	genreAllURI string = "/genre/movie/list"

	sqlHOST string = "127.0.0.1"
	userName string = "postgres"
	password string = "admin"
	port int = 5432
	db string = "tmdb"
)

func dbConfigure() string{

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",userName,password,sqlHOST,port,db)
	//return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d ",sqlHOST,userName,password,db,port)
}

func main(){
	config := dbConfigure()
	fmt.Println(config)
	db, err := gorm.Open(postgres.Open(config),&gorm.Config{
		//some config here....
	})
	if err != nil {
		log.Println(err)
		return
	}

	////////create table
	db.AutoMigrate(&webCrawler.GenreInfo{}) //generate genreInfo Table
	db.AutoMigrate(&webCrawler.MovieInfo{}) // generate movieInfo table
	//

	//TODO - Create all needed genre data with id : with uri:
	genreAPIUri := host + genreAllURI +"?api_key=" + apiKey + "&language=zh-TW"


	//TODO - Insert Data to Database
	err = webCrawler.GenreTableCreate(genreAPIUri,db)
	if err != nil{
		log.Fatalln(err)
		return
	}

	var uris []string
	//var res []*webCrawler.APIResult
	//https://api.themoviedb.org/3/discover/movie?api_key=29570e7acc52b3e085ab46f6a60f0a55&language=zh-TW&sort_by=popularity.desc&include_adult=false&include_video=false&page=1&with_watch_monetization_types=flatrate
	uri := host + allMovieURI +"?api_key=" + apiKey + "&language=zh-TW&sort_by=popularity.desc&include_adult=false&include_video=false&page=1&with_watch_monetization_types=flatrate"
	page := webCrawler.FetchPageInfo(uri)
	uris = uriGenerator(uri,page)

	//store data to db
	webCrawler.FetchMovieInfos(uris,db)

}

func uriGenerator(uri string,page int) []string{
	var uris []string
	for i := 0;i<page;i++{
		newURI := uri + "&page=" + strconv.Itoa(i+1)
		uris = append(uris, newURI)
	}

	return uris
}

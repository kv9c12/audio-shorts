package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"errors"
	"fmt"
	"os"

	// models package where Short schema is defined
	"github.com/kv9c12/audio-shorts/db/models"
	"github.com/kv9c12/audio-shorts/graph/model"
	_ "github.com/lib/pq" // postgres golang driver
)

// create connection with postgres db
func createConnection() (*sql.DB, error) {
    
    // Open the connection
    db, err := sql.Open("postgres", string(os.Getenv("POSTGRES_CONN")))

    if err != nil {
        return nil, err
    }

    // check the connection
    err = db.Ping()

    if err != nil {
        return nil, err
    }

    fmt.Println("Successfully connected!")
    
    // return the connection
    return db, nil
}

// Create a short in the postgres db
func CreateShort(input model.NewShort) (int, error) {
    var err error

    // get the file upload url from s3
    
    fileUrl, err := HandleFileUpload(input.File);

    if err != nil {
        fmt.Println(err)
        return 0, err
    }

    //populate the data to be stored in DB

    var creator models.Creator
    var short models.Short

    short.Title = input.Title
    short.Description = input.Description
    short.Category = input.Category
    creator.Name = input.Name
    creator.Email = input.Email
    short.Creator, err = json.Marshal(creator)
    short.FileUrl = fileUrl

    // call insert short in Postgres
    insertID, err := insertShort(short)

    if err != nil {
        return 0, err
    }

    return insertID, nil
}

// GetShort will return a single short by its id
func GetShort(id int) (models.Short, error){

    var short models.Short

    // call the getShort function with short id to retrieve a single short
    short, err := getShort(id)

    if err != nil {
        fmt.Println("Unable to get short. ", err)
        return short, err 
    }

    return short, nil
}

// GetAllShorts will return all the shorts
func GetAllShorts(page int) ([]models.Short, error) {

    shorts, err := getAllShorts(page)

    if err != nil {
        fmt.Println("Unable to get all short. ", err)
        return nil, err
    }

    return shorts, nil
}

//------------------------- handler functions ----------------
// insert one short in the DB
func insertShort(short models.Short) (int, error) {

    // create the postgres db connection
    db, err := createConnection()

    if err != nil {
        fmt.Println("Create Connection to DB Error: ", err)
        return 0, errors.New("Error Occurred while Creating Connection to DB")
    }

    // close the db connection
    defer db.Close()

    // create the insert sql query
    // returning id will return the id of the inserted short
    sqlStatement := `INSERT INTO shorts (title, description, category, fileUrl, creator) VALUES ($1, $2, $3, $4, $5) RETURNING id`

    // the inserted id will store in this id
    var id int

    // execute the sql statement
    // Scan function will save the insert id in the id
    err = db.QueryRow(sqlStatement, short.Title, short.Description, short.Category, short.FileUrl, short.Creator).Scan(&id)

    if err != nil {
        fmt.Println("Unable to execute the query: ", err)
        return 0, errors.New("Unable to execute the query")
    }

    fmt.Println("Inserted a single record: ", id)

    // return the inserted id
    return id, nil
}

// get one short from the DB by its userid
func getShort(id int) (models.Short, error) {
    // create a short of models.Short type
    var short models.Short

    // create the postgres db connection
    db, err := createConnection()

    if err != nil {
        fmt.Println("Create Connection to DB Error: ", err)
        return short, errors.New("Error Connecting to DB")
    }

    // close the db connection
    defer db.Close()

    // create the select sql query
    sqlStatement := `SELECT * FROM shorts WHERE id=$1`

    // execute the sql statement
    row := db.QueryRow(sqlStatement, id)

    // unmarshal the row object to short
    err = row.Scan(&short.ID, &short.Title, &short.Description, &short.Category, &short.FileUrl, &short.Creator)

    if err != nil {

        if err == sql.ErrNoRows {
            fmt.Println("No rows were returned!")
            return short, errors.New("No rows were returned!")
        }

        fmt.Println("Unable to scan the row. ", err)
        return short, errors.New("Unable to read from DB") 
    }

    return short, nil
}

// get one short from the DB by its userid
func getAllShorts(page int) ([]models.Short, error) {
    // create the postgres db connection
    db,err := createConnection()

    if err != nil {
        fmt.Println("Unable to connect to DB: ",err)
        return nil, errors.New("Unable to establish a connection to DB")
    }

    // close the db connection
    defer db.Close()

    var shorts []models.Short
    offset := 0

    if page > 1 {
        offset = (page-1)*10
    }

    // create the select sql query
    sqlStatement := `SELECT * FROM shorts limit 10 offset $1`

    // execute the sql statement
    rows, err := db.Query(sqlStatement, offset)

    if err != nil {
        fmt.Println("Unable to execute the query. %v", err)
        return nil, errors.New("Unable to execute the query")
    }

    // close the statement
    defer rows.Close()
    
    // iterate over the rows
    for rows.Next() {
        var short models.Short

        // unmarshal the row object to short
        err = rows.Scan(&short.ID, &short.Title, &short.Description, &short.Category, &short.FileUrl, &short.Creator)

        if err != nil {
            fmt.Println("Unable to scan the row. %v", err)
            return nil, errors.New("Unable to read from DB")
        }

        // append the short in the shorts slice
        shorts = append(shorts, short)

    }

    if len(shorts) == 0 {
        shorts = make([]models.Short, 0)
    }

    // return empty short on error
    return shorts, nil
}

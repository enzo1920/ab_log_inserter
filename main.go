package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "database/sql"
   _ "github.com/lib/pq"
     "time"
)

const (
     DB_USER     = "postgres"
     DB_PASSWORD = ""
     DB_NAME     = "ab_log_db"
)


func main() {
        light := getLight("http://192.168.71.74/sec/?pt=31&scl=30&i2c_dev=max44009")
        fmt.Println(light)
        inserter(light)
}

func inserter(light_val string) {
        dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
            DB_USER, DB_PASSWORD, DB_NAME)
        db, err := sql.Open("postgres", dbinfo)
        checkErr(err)
        defer db.Close()

        fmt.Println("# Inserting values")
        dt := time.Now()
        var lastInsertId int
        err = db.QueryRow("INSERT INTO light (light_val,light_date) VALUES($1,$2) returning light_id;", light_val, dt).Scan(&lastInsertId)
        checkErr(err)
        fmt.Println("last inserted id =", lastInsertId)

}

func getLight(host string )(light_val string){
    resp, err := http.Get(host) 
    if err != nil { 
        fmt.Println(err) 
        return
    } 
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
          fmt.Println(err)
          return
    }
    //fmt.Println(string(body))
    light_val = string(body)
    return light_val
}

func checkErr(err error) {
        if err != nil {
            panic(err)
        }
    }


package main

import (
    "fmt"
    "strconv"
    "net/http"
    "io/ioutil"
    "database/sql"
   _ "github.com/lib/pq"
     "time"
)

const (
     DB_USER     = ""
     DB_PASSWORD = ""
     DB_NAME     = "ab_log_db"
)


func main() {
        //counter чтобы получить точно значение с датчика
        cntChecker := 0
        var light string
        for i := 1; i < 10; i++ {
           light = getLight("http://192.168.71.74/sec/?pt=31&scl=30&i2c_dev=max44009")
           if i == 10 || light !="NA"{ 
                break
           }
        cntChecker += i
        }
        //fmt.Println(light)
        inserter(light)
}

func inserter(light_val string) {
        dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
            DB_USER, DB_PASSWORD, DB_NAME)
        db, err := sql.Open("postgres", dbinfo)
        checkErr(err)
        defer db.Close()
        flight_val, err := strconv.ParseFloat(light_val, 64)
        //иногда датчик выдает занчение больше 35000.0 или N/A
        if err != nil || flight_val>=4000.0{
              fmt.Println("error or light_val >4000:")
              fmt.Println(err)
              err = db.QueryRow("select light_val from light order by light_id desc limit 1;").Scan(&flight_val)
              checkErr(err) 
             // return
        }
        //иногда датчик выдает занчение больше 35000.0
        /*if flight_val>=4000.0 {
           err = db.QueryRow("select light_val from light where light order by light_id desc limit 1;").Scan(&flight_val)
           checkErr(err)
        }*/
        insertLightval := fmt.Sprintf("%f",flight_val)
        fmt.Println("# Inserting values")
        fmt.Println(flight_val)
        dt := time.Now()
        var lastInsertId int
        err = db.QueryRow("INSERT INTO light (light_val,light_date) VALUES($1,$2) returning light_id;", insertLightval, dt).Scan(&lastInsertId)
        checkErr(err)
        //fmt.Println("last inserted id =", lastInsertId)

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


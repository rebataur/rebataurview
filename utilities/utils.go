package utilities

import(
    "io/ioutil"
    "encoding/json"
    "log"
)
type DB struct{
  DBType string
  Path string
}
type Config struct{
  Database DB
}

func ReadConfig() Config{

    var config Config
    content,_ := ioutil.ReadFile("../config.json")

    if err := json.Unmarshal(content,&config); err != nil{
      log.Println(err)
    }
    log.Println(config.Database.DBType)
    return config
}

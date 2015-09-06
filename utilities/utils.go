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

func ReadConfig() (Config){

    var config Config
    // Read the config.json file
    if content,err := ioutil.ReadFile("../config.json"); err == nil{
      // Unmarshal the config.json file
      if err := json.Unmarshal(content,&config); err != nil{
        log.Fatal(err, "Error parsing config.json file")
      }
    }else{
      log.Fatal(err,"Error in reading config.json file, check whether present")
    }


    return config
}

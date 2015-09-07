package utilities

import(
    "io/ioutil"
    "encoding/json"
    "log"
    "os"
    "strings"
)
type DB struct{
  DBType string
  DBPath string
}
type Config struct{
  Database DB
}

func ReadConfig() (Config,error){

    var config Config
    // Read the config.json file
    pwd,err := os.Getwd()
    if err != nil{
      log.Println("Error getting Working Dir", err)
      return config,err
    }
    var fullPath string = strings.Join([]string{pwd,"\\config.json"},"")
    log.Println(fullPath)
    if content,err := ioutil.ReadFile(fullPath); err == nil{
      // Unmarshal the config.json file
      if err := json.Unmarshal(content,&config); err != nil{
        log.Println(err, "Error parsing config.json file")
        return config,err
      }
    }else{
      log.Println(err,"Error in reading config.json file, check whether present")
      return config,err
    }


    return config,nil
}

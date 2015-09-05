package utilities

import(
  "log"
  "testing"
)


func TestReadConfig(t *testing.T){
  config := ReadConfig()

  if config.Database.DBType == "pg"{
    log.Println("Passed **")
  }else {
    t.Fail()
  }
}

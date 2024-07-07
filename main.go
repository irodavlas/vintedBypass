package main

import (
  "github.com/vintedMonitor/utils"
  "log"
)
func main(){
   
  client, err := utils.NewClient("https://www.vinted.com")
  if err != nil {
    log.Fatal(err)
  }
  //start mointor for starting page
}



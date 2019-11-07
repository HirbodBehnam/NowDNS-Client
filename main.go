package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const VERSION = "1.0.0 / Build 1"

func main() {
	Hostname := flag.String("h", "", "The host of now dns")
	Username := flag.String("u", "", "Username of now dns")
	Password := flag.String("p", "", "Password of now dns")
	Interval := flag.Int("interval",600,"The interval between refreshing the IP address. In seconds")
	LogFile := flag.String("log","","The log file that the app writes log into it. If specified, it will not log to terminal")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		fmt.Println("Created by Hirbod Behnam")
		fmt.Println("Source at https://github.com/HirbodBehnam/NowDNS-Client")
		fmt.Println("Version", VERSION)
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *Hostname == "" || *Username == "" || *Password == ""{
		fmt.Println("Please fill hostname, username and password")
		flag.PrintDefaults()
		os.Exit(0)
	}

	//Set log file
	if *LogFile != ""{
		f, err := os.OpenFile(*LogFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error opening log file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)
	}

	log.Println("Starting the client")
	var res string
	var err error
	for {
		res, err = basicAuth(*Hostname,*Username,*Password)
		if err != nil{
			log.Println("Error on connection:",err.Error())
		}else{
			log.Println(res)
		}
		time.Sleep(time.Second * time.Duration(*Interval))
	}
}

//https://stackoverflow.com/q/16673766/4213397
func basicAuth(hostname,username,password string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://now-dns.com/update?hostname=" + hostname, nil)
	if err != nil{
		return "", err
	}
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil{
		return "", err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	return s , nil
}
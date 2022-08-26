package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/juanobligado/go-png/internal/png"
)

 var filenameFlag = "PNG File Name"
 var help = flag.Bool("help", false, "please run png-meam <<file name>> to get a png mean")
func main(){

	flag.StringVar(&filenameFlag,"file"," <<Path to PNG File >>","filename")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	
	data,err := png.ReadPng(filenameFlag)
	if err != nil{
		fmt.Print(err)
		os.Exit(0)
	}
	mean := data.Mean()
	fmt.Printf("%v",mean)
}
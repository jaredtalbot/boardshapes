package main

import "flag"
import "log"
import "os"
import "strings"
import "image"
import "image/png"
import "image/jpeg"





func main() {
	var simplified = flag.String("s", "./test_images/out.png", "For simplification to outputfile")
	var meshes = flag.Bool("m", false, "For mesh conversion")
	var resize = flag.Bool("r", false , "For Sizing" )
	var file_intake = flag.Arg(0) 

	if len(file_intake) == 0 {
		log.Println("No File Name provided")
		os.Exit(1)
	
	}else if !strings.Contains(file_intake, ".jpeg") && !strings.Contains(file_intake, ".png") {
		log.Println(".jpeg or .png wasnt found")
	

	}else 
		
	flag.Parse()
	


}
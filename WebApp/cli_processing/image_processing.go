package main

import (
	"codejester27/cmps401fa2024/web-app/processing"
	"flag"
	"fmt"
	"mime/quotedprintable"
	"path/filepath"

	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

var rFlag = flag.Bool("r", false, "this should be resizing if called")
var sFlag = flag.String("s", "", "Should allow a user to input a file path for output. Unsure on this")
var mFlag = flag.Bool("m", false, "Should meshify regions.")

func main() {
	flag.Parse()
	fileInput := flag.Args()
	fileToOutput, inputDir := file_opener(fileInput)



	if *rFlag {
		newResizedImage, err := processing.ResizeImage(fileToOutput)
			if err != nil {
				"fix your stuff bruh"
			}


		// umm image get image resize image.
		// should allow me to resize image but should be after file Open.
		//fileExtension := filepath.Ext(fileToOutput)
		//if fileExtension == ".png" || ".jpeg" {
			//fileDecoded, err = image.Decode(fileToOutput)
			//err, fileResized := processing.ResizeImage(fileDecoded)
		
	

	}

	image_output(fileToOutput, inputDir)

	fmt.Println("1:", *rFlag)
	fmt.Println("2:", *sFlag)
	fmt.Println("3:", *mFlag)
	fmt.Println("4:", fileInput)

	
}
	

func file_opener_decoder(fileInput []string) (*os.File, string) {
	
	joinedFileName := strings.Join(fileInput, "")
	fileTaken, err := os.Open(joinedFileName)
	if err != nil {
		panic(err)
	}
	//fmt.Println(fileTaken)
	inputDir := filepath.Dir(joinedFileName)

	fileExtension := filepath.Ext(fileTaken)
	if fileExtension == ".png" || ".jpeg" {
		img, extName, err :=image.Decode(fileTaken)
		if err != nil {
			panic(err)
		}


	return img, inputDir
}



func image_output(fileToOutput *os.File, inputDir string) {

	outputDir := filepath.Join(inputDir, "output_files")
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}
	outputFilePath := filepath.Join(outputDir, filepath.Base(fileToOutput.Name()))
	fmt.Printf("Output file path: %s\n", outputFilePath)
}



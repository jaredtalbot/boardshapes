package main

import (
	"codejester27/cmps401fa2024/web-app/processing"

	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

var rFlag = flag.Bool("r", false, "this should be resizing if called")
var sFlag = flag.String("s", "", "Should allow a user to input a file path for output. Unsure on this")
var mFlag = flag.Bool("m", false, "Should meshify regions.")

func main() {

	flag.Parse()
	fileInput := flag.Args()
	fileToOutput, err := fileOpenerDecoder(fileInput)

	if err != nil {
		panic(err)
	}

	if *rFlag {
		img, err := processing.ResizeImage(fileToOutput)
		if err != nil {
			fmt.Println("fix your stuff bruh")
		}
		if *mFlag {
			fileRegioned, regionCount,_  := processing.SimplifyImage(img, processing.RegionMapOptions{})
			fmt.Println(regionCount)
			outputFile := fileEncoder(fileRegioned)
			image_output(outputFile, filepath.Dir(fileInput[0]))
		} else {
			outputFile := fileEncoder(img)
			image_output(outputFile, filepath.Dir(fileInput[0]))
		}
	}

	fmt.Println("1:", *rFlag)
	fmt.Println("2:", *sFlag)
	fmt.Println("3:", *mFlag)
	fmt.Println("4:", fileInput)
	

}

func fileOpenerDecoder(fileInput []string) (image.Image, error) {

	joinedFileName := strings.Join(fileInput, "")

	fileTaken, err := os.Open(joinedFileName)
	if err != nil {
		panic(err)
	}
	defer fileTaken.Close()

	fileExtension := filepath.Ext(joinedFileName)
	if fileExtension == ".jpg" {
		fileExtension = ".jpeg"
	}
	if fileExtension == ".png" || fileExtension == ".jpeg"  {
		img, _, err := image.Decode(fileTaken)
		if err != nil {
			panic(err)
		}

		return img, nil
	}
	return nil, fmt.Errorf("unsuppported file format")
}

func fileEncoder(img image.Image) *os.File {

	var outputPath string
	if *sFlag != "" {
		outputPath = *sFlag
	} else {
		outputPath = "output.png"
	}
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	ext := strings.ToLower(filepath.Ext(outputPath))
	if ext == ".png" {
		err = png.Encode(outputFile, img)
	} else if ext == ".jpeg" || ext == ".jpg" {
		err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: 90})

	}
	if err != nil {
		panic(err)
	}
	return outputFile
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

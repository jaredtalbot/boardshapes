package main

import (
	"flag"
	"fmt"
	"path/filepath"

	//"image"
	//"image/png"
	//"image/jpeg"
	"os"
	"strings"
)

var rFlag = flag.Bool("r", false, "this should be resizing if called")
var sFlag = flag.String("s", "", "Should allow a user to input a file path. Unsure on this")
var mFlag = flag.Bool("m", false, "Should meshify regions.")

func main() {
	flag.Parse()
	var fileInput = flag.Args()
	fileToOutput, inputDir := image_open()
	image_output(fileToOutput, inputDir)

	fmt.Println("1:", *rFlag)
	fmt.Println("2:", *sFlag)
	fmt.Println("3:", *mFlag)
	fmt.Println("4:", fileInput)

}

func image_open() (*os.File, string) {
	var fileInput = flag.Args()
	joinedFileName := strings.Join(fileInput, "")
	fileTaken, err := os.Open(joinedFileName)
	if err != nil {
		panic(err)
	}
	fmt.Println(fileTaken)
	inputDir := filepath.Dir(joinedFileName)

	return fileTaken, inputDir
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

//var simplified = flag.String("s", "./test_images/out.png", "For simplification to outputfile")
//var meshes = flag.Bool("m", false, "For mesh conversion")
//var resize = flag.Bool("r", false, "For Sizing")
//var file_intake = flag.Arg(0)

//if len(file_intake) == 0 {
//	log.Println("No File Name provided")
//	os.Exit(1)

//} else if !strings.Contains(file_intake, ".jpeg") && !strings.Contains(file_intake, ".png") {
//	file, fileErr := os.Open(file_intake)
//
//		if fileErr != nil {
//			panic(fileErr)

//		}
//		defer file.Close()

//		img, format, err := image.Decode(file)
//		if err != nil {
//			fmt.Print(err)

//		} else {
//			flag.Parse()
//		}
//	}
//}

//else if !strings.Contains(file_intake, ".jpeg") && !strings.Contains(file_intake, ".png") {
//image.Decode(file_intake)

//}//else

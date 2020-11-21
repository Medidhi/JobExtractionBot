package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type jobsearch struct {
	Position string
	Location string
}

type email struct {
	EmailPassword string
	From          string
	TO            string
	Subject       string
	Body          string
}

type Config struct {
	Jobsearchs      []jobsearch
	OutputDirectory string
	Email           email
}

func LoadConfiguration() Config {
	var config Config

	viper.SetConfigName("config") // config file name without extension
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // config file path

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.UnmarshalKey("jobsearches", &config.Jobsearchs)
	viper.UnmarshalKey("outputFilePath", &config.OutputDirectory)
	viper.UnmarshalKey("email", &config.Email)
	return config
}

func ExtractJob(position string, localtion string, outputPath string) {

	position = strings.ReplaceAll(position, " ", "#")
	localtion = strings.ReplaceAll(localtion, " ", "#")
	outputPath = strings.ReplaceAll(outputPath, " ", "#")

	cmd := exec.Command("python3", "IndeedScraper.py", position, localtion, outputPath)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}

}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		fmt.Println(basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + "/"
			fmt.Println("Recursing and Adding SubDir: " + file.Name())
			fmt.Println("Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}

func ZipWriter(source string, destination string) {
	baseFolder := source

	// Get a Buffer to Write To
	outFile, err := os.Create(destination)
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	addFiles(w, baseFolder, "")

	if err != nil {
		fmt.Println(err)
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func sendEmails(config Config) {
	os.MkdirAll(config.OutputDirectory, os.ModePerm)
	ZipWriter(config.OutputDirectory, "./jobs.zip")
	cmd := exec.Command("python3", "sendEmails.py", config.Email.From, config.Email.EmailPassword, config.Email.TO, "./jobs.zip")
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Sent job listing email")
}

func createCalendarEvent(stringDate string) {
	cmd := exec.Command("python3", "createCalendarEvent.py", stringDate)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	c := LoadConfiguration()
	// parallel processing
	for _, ele := range c.Jobsearchs {
		fmt.Printf("pos ---- ( %+v at %+v) ----- started\n", ele.Position, ele.Location)
		ExtractJob(ele.Position, ele.Location, c.OutputDirectory)
	}

	fmt.Printf("Finished job parsing")
	sendEmails(c)

	// add two days to current date to get set the calenndar event
	t := time.Now().AddDate(0, 0, 2)
	createCalendarEvent(t.Format("2006-01-02"))

}

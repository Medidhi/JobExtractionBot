package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

type jobsearch struct {
	Position string
	Location string
}

type Config struct {
	Jobsearchs      []jobsearch
	OutputDirectory string
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

func main() {
	c := LoadConfiguration()
	// parallel processing
	for _, ele := range c.Jobsearchs {
		fmt.Printf("pos ---- ( %+v at %+v) ----- started\n", ele.Position, ele.Location)
		ExtractJob(ele.Position, ele.Location, c.OutputDirectory)
	}

	fmt.Printf("Finished job parsing")

}

package main

import (
   // "github.com/woanware/goutil"
)
import (
    "os"
    "log"
    "bufio"
   // "fmt"
    "regexp"
    "io/ioutil"
    "path"
    "io"
   // "fmt"
)

// ##### Types ###############################################################

// Encapsulates a Processor object and its properties
type Processor struct {
    id 						int
    config					*Config
}


// ##### Methods #############################################################

// Constructor/Initialiser for the Processor struct
func NewProcessor(id int, config *Config) *Processor {
    p := Processor{
        id:     id,
        config: config,
    }

    return &p
}

// Process an individual set of host data
func (p Processor) Process(filePath string) {
    defer wg.Done()

    inputFileDir := path.Dir(filePath)
    inputFilePath := filePath
    //outputFilePath := ""
    isFirstFile := true

    var r *regexp.Regexp
    for _, r = range regexps {
        inputFilePath = p.ProcessRegex(isFirstFile, r, inputFileDir, inputFilePath)
        isFirstFile = false
    }
}

//
func (p Processor) ProcessRegex(isFirstFile bool, regExp *regexp.Regexp, inputDir string, inputFilePath string) string {

    inputFile, err := os.Open(inputFilePath)
    if err != nil {
        log.Fatal(err)
    }
    defer inputFile.Close()

    tempFile, err := ioutil.TempFile(inputDir, "lfd-")
    if err != nil {
        log.Fatal(err)
    }
    defer tempFile.Close()

    reader := bufio.NewReader(inputFile)

    var line string
    for {
        line, err = reader.ReadString('\n')

        if err == io.EOF { break }
        if err != nil { break }

        if regExp.MatchString(line) == false {
            tempFile.WriteString(line)
        }
    }

    if isFirstFile == false {
        defer os.Remove(inputFile.Name())
    }

    // Return the temp file path, as that will be used in the next iteration
    return tempFile.Name()
}


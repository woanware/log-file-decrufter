package main

import (
    "os"
    "log"
    "bufio"
    "regexp"
    "io/ioutil"
    "path"
    "io"
    "path/filepath"
)

// ##### Types ###############################################################

// Encapsulates a Processor object and its properties
type Processor2 struct {
    id 						int
    config					*Config
}

// ##### Methods #############################################################

// Constructor/Initialiser for the Processor struct
func NewProcessor2(id int, config *Config) *Processor2 {
    p := Processor2{
        id:     id,
        config: config,
    }

    return &p
}

// Process an individual set of host data
func (p Processor2) Process(filePath string) {
    defer wg.Done()

    inputFileName := path.Base(filePath)
    inputFileDir := path.Dir(filePath)
    inputFilePath := filePath
    outputFilePath :=  filepath.Join(inputFileDir, FILE_PREFIX + inputFileName)

    inputFile, err := os.Open(inputFilePath)
    if err != nil {
        log.Fatal(err)
    }
    defer inputFile.Close()

    tempFile, err := ioutil.TempFile(inputFileDir, FILE_PREFIX)
    if err != nil {
        log.Fatal(err)
    }
    defer tempFile.Close()

    reader := bufio.NewReader(inputFile)
    var line string
    var matched bool
    var r *regexp.Regexp
    for {
        line, err = reader.ReadString('\n')

        if err == io.EOF { break }
        if err != nil { break }

        matched = false
        for _, r = range regexps {
            if r.MatchString(line) == true {
                matched = true
                break
            }
        }

        if matched == false {
            tempFile.WriteString(line)
        }
    }

    defer os.Rename(tempFile.Name(), outputFilePath)
}

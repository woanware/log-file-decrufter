package main

import (
    "github.com/voxelbrain/goptions"
    "fmt"
    "regexp"
    "path/filepath"
    "os"
    "runtime"
    "sync"
    "strings"
    "log"
)

// ##### Variables ############################################################

var (
    config      *Config
    regexps     []*regexp.Regexp
    workQueue   chan string
    wg          sync.WaitGroup
)

// ##### Constants ############################################################

const APP_TITLE string = "log-file-decrufter"
const APP_NAME string = "lfd"
const APP_VERSION string = "0.0.1"
const FILE_PREFIX string = "lfd-"

// ##### Methods ##############################################################

func main() {

    fmt.Printf("\n%s (%s) %s\n\n", APP_TITLE, APP_NAME, APP_VERSION)

    opt := struct {
        InputDir 	string        `goptions:"-i, --inputdir, obligatory, description='Input directory'"`
        Help		goptions.Help `goptions:"-h, --help, description='Show this help'"`
    }{}

    goptions.ParseAndFail(&opt)

    fileInfo, err := os.Stat(opt.InputDir)
    if err == nil {
        if fileInfo.IsDir() == false {
            fmt.Println("The input directory value is not a directory")
            os.Exit(-1)
        }
    } else {
        fmt.Printf("Error checking the input directory: %v \n", err)
        os.Exit(-1)
    }

    if os.IsNotExist(err) {
        fmt.Println("The input directory does not exist")
        os.Exit(-1)
    }

    config, err =  LoadConfig("log-file-decruft.config")
    if err != nil {
        fmt.Printf("%s", err.Error())
        os.Exit(-1)
    }

    if len(config.Regex.Regexes) == 0 {
        fmt.Println("No regexes loaded")
        os.Exit(-1)
    }

    fmt.Printf("Loaded %d regexps\n", len(config.Regex.Regexes))

    var regex *regexp.Regexp
    regexps = make([]*regexp.Regexp, 0)
    for _, v := range config.Regex.Regexes {
        fmt.Printf("Compiling regex: %s\n", v)
        regex, err = regexp.Compile(v)
        if err != nil {
            log.Fatalf("Error compiling regex: %v (%s)", err, v)
            return
        }

        regexps = append(regexps, regex)
    }

    createProcessors()

    err = filepath.Walk(opt.InputDir, func(path string, fi os.FileInfo, err error) error {
        if fi.IsDir() == true {
            return nil
        }

        if strings.HasPrefix(fi.Name(), FILE_PREFIX) == true {
            return nil
        }

        wg.Add(1)
        workQueue <- path
        return nil
    })

    wg.Wait()
}

// Initialise the channels for the cross process comms and then start the workers
func createProcessors() {
    processorCount := runtime.NumCPU()
    if config.Misc.ProcessorThreads > 0 {
        processorCount = config.Misc.ProcessorThreads
    }

    workQueue = make(chan string, 100)

    // Create the workers that perform the actual processing
    for i := 0; i < processorCount; i++ {
        p := NewProcessor2(i, config)
        go func (p *Processor2) {
            for j := range workQueue {
                p.Process(j)
            }
        } (p)
    }
}




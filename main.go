package main

import (
    "github.com/davecgh/go-spew/spew"
    jobLib "jobLib/jobLibrary"
)

func main() {
    job := new(jobLib.JobLibrary)
    job.Init()
    job.SetDomain("test")
    spew.Dump(job.GetJobActiveStatus())
}
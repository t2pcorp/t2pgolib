package main

import (
    jobLib "github.com/t2pcorp/t2pgolib.git/jobLibrary"
)

func main() {
    job := new(jobLib.JobLibrary)
    job.Init()
    job.SetDomain("Test_UpdateJobStatus")
    job.SetJobID("Test_UpdateJobStatus")
    job.SetName("Test_UpdateJobStatus")
    job.SetPeriodTypeMin()
    job.SetPeriodValue("1")
    job.SetScheduleTime("12:00:00")
    job.SetExecuteDuration("1")
    job.SetSkipCheck("true")
    job.SetLINENotification("ads")
    job.SetNotiFrequency("1")
    job.SetArchiveLogUnit("D")
    job.SetArchiveLogValue("1")

    job.UpdateJobRunningStatus()
}
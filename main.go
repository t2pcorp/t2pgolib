package main

import (
	// "fmt"
	jobLib "github.com/t2pcorp/t2pgolib/jobLibrary"
)

func main() {
	job := new(jobLib.JobLibrary)
	job.Init("DEVELOP")

	job.SetEmail("test@example.com")
	job.SetPassword("123456789")
	job.SetDomain("EXAMPLE")
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

	status := job.GetJobActiveStatus()
	if status != "N" {
		/*



		   TODO:: Do Job Process Here




		*/
		job.UpdateJobStatus()
		job.UpdateJobDashboard(10, jobLib.DashboardMatricKey{DimensionName: "default", Matric: "Monitor", CustomNamespace: "default"})
	} else {
		job.UpdateJobRunningStatus()
	}
}

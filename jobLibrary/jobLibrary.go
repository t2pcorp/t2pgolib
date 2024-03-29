package jobLibrary

import (
    "io/ioutil"
    "bytes"
    "fmt"
    "net/http"
    "encoding/json"
	"errors"
	"time"
    "github.com/aws/aws-sdk-go/aws"
    // "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudwatch"
	// "gitlab.t2p.co.th/central-library/t2plib/config"
)

type credentials struct {
	Email string 
	Password string
}

type DashboardMatricKey struct {
	DimensionName string `default`
	Matric	string `Monitor`
	CustomNamespace string `default`// if use this will ignore DimensionName and MatricName
}
type JobLibrary struct {
	JobConfig               JobConfig
	JobExecuteInfo          JobExecuteInfo
	JobUserPassword 		credentials
	Env                     string
}

type JobConfig struct {
	Domain            string
	JobID             string
	Name              string
	PeriodType        string
	PeriodValue       string
	ScheduleTime      string
	ExecuteDuration   string
	TimeZone          string
	AdditionCondition AdditionCondition
	SkipCheck         string
	Notification      Notification
	NotiFrequency     string
	ArchiveLogUnit    string
	ArchiveLogValue   string
}

type AdditionCondition struct {
	Success bool
}

type Notification struct {
	Line string
	Sms  string
	Call string
	Mail string
	hasSet bool
}

type JobExecuteInfo struct {
	Success bool
	Error   string
}

var user = credentials{}

func (job *JobLibrary) Init(env string) *JobLibrary {
	job.JobConfig.TimeZone = "Asia/Bangkok"
	job.JobConfig.AdditionCondition.Success = true
	job.JobExecuteInfo.Success = true
	job.Env = env
	return job
}

//Getter

func (j *JobLibrary) GetDomain() string {
	return j.JobConfig.Domain
}

func (j *JobLibrary) GetJobID() string {
	return j.JobConfig.JobID
}

func (j *JobLibrary) GetName() string {
	return j.JobConfig.Name
}

func (j *JobLibrary) GetPeriodType() string {
	return j.JobConfig.PeriodType
}

func (j *JobLibrary) GetPeriodValue() string {
	return j.JobConfig.PeriodValue
}

func (j *JobLibrary) GetScheduleTime() string {
	return j.JobConfig.ScheduleTime
}

func (j *JobLibrary) GetExecuteDuration() string {
	return j.JobConfig.ExecuteDuration
}

func (j *JobLibrary) GetTimeZone() string {
	return j.JobConfig.TimeZone
}

func (j *JobLibrary) GetAdditionCondition() AdditionCondition {
	return j.JobConfig.AdditionCondition
}

func (j *JobLibrary) GetSkipCheck() string {
	return j.JobConfig.SkipCheck
}

func (j *JobLibrary) GetNotification() Notification {
	return j.JobConfig.Notification
}

func (j *JobLibrary) GetNotiFrequency() string {
	return j.JobConfig.NotiFrequency
}

func (j *JobLibrary) GetArchiveLogUnit() string {
	return j.JobConfig.ArchiveLogUnit
}

func (j *JobLibrary) GetArchiveLogValue() string {
	return j.JobConfig.ArchiveLogValue
}

//Setter

func (j *JobLibrary) SetDomain(domain string) {
	j.JobConfig.Domain = domain
}

func (j *JobLibrary) SetJobID(jobID string) {
	j.JobConfig.JobID = jobID
}

func (j *JobLibrary) SetName(name string) {
	j.JobConfig.Name = name
}

func (j *JobLibrary) SetPeriodTypeMin() {
	j.JobConfig.PeriodType = "min"
}

func (j *JobLibrary) SetPeriodTypeDaily() {
	j.JobConfig.PeriodType = "daily"
}

func (j *JobLibrary) SetPeriodTypeDate() {
	j.JobConfig.PeriodType = "date"
}

func (j *JobLibrary) SetPeriodTypeDateMonth() {
	j.JobConfig.PeriodType = "datemonth"
}

func (j *JobLibrary) SetPeriodTypeOnce() {
	j.JobConfig.PeriodType = "once"
}

func (j *JobLibrary) SetPeriodValue(periodValue string) {
	j.JobConfig.PeriodValue = periodValue
}

func (j *JobLibrary) SetScheduleTime(scheduleTime string) {
	j.JobConfig.ScheduleTime = scheduleTime
}

func (j *JobLibrary) SetExecuteDuration(executeDuration string) {
	j.JobConfig.ExecuteDuration = executeDuration
}

func (j *JobLibrary) SetTimeZone(timeZone string) {
	j.JobConfig.TimeZone = timeZone
}

func (j *JobLibrary) SetAdditionCondition(success bool) {
	j.JobConfig.AdditionCondition.Success = success
}

func (j *JobLibrary) SetSkipCheck(skip string) {
	j.JobConfig.SkipCheck = skip
}

func (j *JobLibrary) SetLINENotification(token string) {
	j.JobConfig.Notification.Line = token
	j.JobConfig.Notification.hasSet = true
}

func (j *JobLibrary) SetSMSNotification(phoneNumber string) {
	j.JobConfig.Notification.Sms = phoneNumber
	j.JobConfig.Notification.hasSet = true
}

func (j *JobLibrary) SetPhoneNotification(phoneNumber string) {
	j.JobConfig.Notification.Call = phoneNumber
	j.JobConfig.Notification.hasSet = true
}

func (j *JobLibrary) SetMailNotification(mail string) {
	j.JobConfig.Notification.Mail = mail
	j.JobConfig.Notification.hasSet = true
}

func (j *JobLibrary) SetNotiFrequency(frequency string) {
	j.JobConfig.NotiFrequency = frequency
}

func (j *JobLibrary) SetArchiveLogUnit(archiveLogUnit string) {
	j.JobConfig.ArchiveLogUnit = archiveLogUnit
}

func (j *JobLibrary) SetArchiveLogValue(archiveLogValue string) {
	j.JobConfig.ArchiveLogValue = archiveLogValue
}

func (j *JobLibrary) SetSuccess(success bool) {
	j.JobExecuteInfo.Success = success
}

func (j *JobLibrary) SetErrorMessage(message string) {
	j.JobExecuteInfo.Error = message
}

func (j *JobLibrary) SetEmail(email string) {
	j.JobUserPassword.Email = email
}

func (j *JobLibrary) SetPassword(password string) {
	j.JobUserPassword.Password = password
}

func (j *JobLibrary) CheckField() error{
    if j.JobConfig.Domain == "" {
        return errors.New("field Domain need to assigned.")
    }
    if j.JobConfig.JobID == "" {
        return errors.New("field JobID need to assigned.")
    }
    if j.JobConfig.Name == "" {
        return errors.New("field Name need to assigned.")
    }
    if j.JobConfig.PeriodType == "" {
        return errors.New("field PeriodType need to assigned.")
    }
    if j.JobConfig.PeriodValue == "" {
        return errors.New("field PeriodValue need to assigned.")
    }
    if j.JobConfig.ScheduleTime == "" {
        return errors.New("field ScheduleTime need to assigned.")
    }
    if j.JobConfig.ExecuteDuration == "" {
        return errors.New("field ExecuteDuration need to assigned.")
    }
    if j.JobConfig.TimeZone == "" {
        return errors.New("field TimeZone need to assigned.")
    }
    if j.JobConfig.SkipCheck == "" {
        return errors.New("field SkipCheck need to assigned.")
    }
    if j.JobConfig.Notification.hasSet == false {
        return errors.New("field Notification need to assigned.")
    }
    if j.JobConfig.NotiFrequency == "" {
        return errors.New("field NotiFrequency need to assigned.")
    }
    if j.JobConfig.ArchiveLogUnit == "" {
        return errors.New("field ArchiveLogUnit need to assigned.")
    }
    if j.JobConfig.ArchiveLogValue == "" {
        return errors.New("field ArchiveLogValue need to assigned.")
    }
    return nil
}

func GetEnvUrl(env string) string {
	url := "https://job-api.t2p.co.th"
	if env == "LOCAL" {
		url = "http://localhost:7005"
	}
	if env == "DEVELOP" {
		url = "https://dev-job-api.t2p.co.th"
	}
	if env == "SIT" {
		url = "https://sit-job-api.t2p.co.th"
	}
	if env == "UAT" {
		url = "https://test-job-api.t2p.co.th"
	}
	return url
}

func GetToken(job *JobLibrary) string{
	env := job.Env
	urlEnv := GetEnvUrl(env)
	var jsonData = []byte(fmt.Sprintf(`{
		"email": "%s",
		"password": "%s"
	}`,job.JobUserPassword.Email,job.JobUserPassword.Password))
	request, err := http.NewRequest("POST", urlEnv+`/api/login`, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		fmt.Println("Login API doRequest Error : ", err)
		return ""
	}

	body, _ := ioutil.ReadAll(response.Body)
	var login map[string]interface{}
    json.Unmarshal([]byte(body), &login)
	if login["success"] == false {
		fmt.Println("Login API Error : ", login["message"])
		return ""
	}
	return login["token"].(string)
}

func (j *JobLibrary) GetJobActiveStatus() string{
	err := j.CheckField()
    if err != nil {
        fmt.Println(err.Error())
    }
    jsonData, err := json.Marshal(j)
    if err != nil {
        fmt.Println(err.Error())
    }
	
	bearer := "Bearer " + GetToken(j)
	env := j.Env
	urlEnv := GetEnvUrl(env)
    url := urlEnv + "/api/Job/getJobStatus/" + j.JobConfig.Domain + "/" + j.JobConfig.JobID
    request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    request.Header.Set("Authorization", bearer)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		fmt.Println("GetJobActiveStatus Error : ", err)
		return "Y"
	}

	data, _ := ioutil.ReadAll(response.Body)
	var job map[string]interface{}
    json.Unmarshal([]byte(data), &job)
    return job["getJobStatus"].(string)
}

func (j *JobLibrary) UpdateJobStatus(msg ...string) {
	err := j.CheckField()
    if err != nil {
        fmt.Println(err.Error())
    }
	if len(msg) > 0 {
		j.SetSuccess(false)
		j.SetErrorMessage(msg[0])
	}
    jsonData, err := json.Marshal(j)
    if err != nil {
        fmt.Println(err.Error())
    }

	bearer := "Bearer " + GetToken(j)
	env := j.Env
	urlEnv := GetEnvUrl(env)
    url := urlEnv + "/api/Job/updateJobStatus"
    request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    request.Header.Set("Authorization", bearer)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		fmt.Println("Update Job Status Error : ", err)
		return
	}

	data, _ := ioutil.ReadAll(response.Body)
	var job map[string]interface{}
    json.Unmarshal([]byte(data), &job)
}

func (j *JobLibrary) UpdateJobRunningStatus() {
	err := j.CheckField()
    if err != nil {
        fmt.Println(err.Error())
    }
	var jsonData = []byte(`{
		"status": "Y"
	}`)
    if err != nil {
        fmt.Println(err.Error())
    }

	bearer := "Bearer " + GetToken(j)
	env := j.Env
	urlEnv := GetEnvUrl(env)
    url := urlEnv + "/api/Job/updateJobRunningStatus/" + j.JobConfig.Domain + "/" + j.JobConfig.JobID
    request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    request.Header.Set("Authorization", bearer)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		fmt.Println("Login API doRequest Error : ", err)
		// panic(err)
	}
}

func (j *JobLibrary) UpdateJobDashboard(value float64, dashboardKey DashboardMatricKey) {
	DimensionName := dashboardKey.DimensionName
	Matric := dashboardKey.Matric
	customNamespace := dashboardKey.CustomNamespace
	namespace := j.JobConfig.Domain + `:` + j.JobConfig.JobID
	if (customNamespace != "default") {
		namespace = customNamespace
	}	
	now := time.Now()
	metricData := []*cloudwatch.MetricDatum{
        &cloudwatch.MetricDatum{
            MetricName: aws.String(Matric),
			Timestamp:	&now,
            Dimensions: []*cloudwatch.Dimension{
                &cloudwatch.Dimension{
                    Name:  aws.String(Matric),
                    Value: aws.String(DimensionName),
                },
            },
            Unit:       aws.String("Count"),
            Value:      aws.Float64(value),
        },
	}
	
	
	matricInput := cloudwatch.PutMetricDataInput{
		MetricData: metricData,
		Namespace: aws.String(namespace),
	}
	j.PutMetricData(matricInput)
}

func (j *JobLibrary) PutMetricData(matricInput cloudwatch.PutMetricDataInput) {
	// SET in ENV
	// $ export AWS_ACCESS_KEY_ID=YOUR_AKID
	// $ export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY
	config := &aws.Config{
		Region: aws.String("ap-southeast-1"),
	}
	mySession := session.Must(session.NewSession(config))

	// Create a CloudWatch client from just a session.
	svc := cloudwatch.New(mySession)
	rs, err := svc.PutMetricData(&matricInput)
	if err != nil {
		fmt.Println("Error adding metrics:", err.Error())
		return
	}
	fmt.Println("Put Matric Result: ",rs)
}
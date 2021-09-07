package jobLibrary

type JobLibrary struct {
	JobConfig               JobConfig
	JobExecuteInfo          JobExecuteInfo
	env                     string
	JobConfigHasChange      JobConfigHasChange
	JobExecuteInfoHasChange JobExecuteInfoHasChange
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
}

type JobExecuteInfo struct {
	Success bool
	Error   string
}

type JobConfigHasChange struct {
	Domain            bool
	JobID             bool
	Name              bool
	PeriodType        bool
	PeriodValue       bool
	ScheduleTime      bool
	ExecuteDuration   bool
	TimeZone          bool
	AdditionCondition bool
	SkipCheck         bool
	Notification      bool
	NotiFrequency     bool
	ArchiveLogUnit    bool
	ArchiveLogValue   bool
}

type JobExecuteInfoHasChange struct {
	Success bool
	Error   bool
}

func (job *JobLibrary) Init() *JobLibrary {
	job.JobConfig.TimeZone = "Asia/Bangkok"
	job.JobConfig.AdditionCondition.Success = true
	job.JobExecuteInfo.Success = true
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
}

func (j *JobLibrary) SetSMSNotification(phoneNumber string) {
	j.JobConfig.Notification.Sms = phoneNumber
}

func (j *JobLibrary) SetPhoneNotification(phoneNumber string) {
	j.JobConfig.Notification.Call = phoneNumber
}

func (j *JobLibrary) SetMailNotification(mail string) {
	j.JobConfig.Notification.Mail = mail
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
	if env == "TEST" {
		url = "https://test-job-api.t2p.co.th"
	}
	return url
}

func GetToken() string {
	config := "config"
	return config
}

func (j *JobLibrary) GetJobActiveStatus() (string, string, JobLibrary) {
	domain := j.JobConfig.Domain
	jobID := j.JobConfig.JobID

	config := JobLibrary{JobConfig: j.JobConfig, JobExecuteInfo: j.JobExecuteInfo}
	return domain, jobID, config
}

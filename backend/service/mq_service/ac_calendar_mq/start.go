package ac_calendar_mq

import "time"

func StartACProblemSenderHourly() {
	now := time.Now()
	nextHour := now.Truncate(time.Hour).Add(time.Hour)

	time.Sleep(time.Until(nextHour))

	SendACProblemsByTime(time.Now())

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		SendACProblemsByTime(time.Now())
	}
}

func StartACProblemSenderEveryMinute() {
	now := time.Now()
	nextMinute := now.Truncate(time.Minute).Add(time.Minute)

	time.Sleep(time.Until(nextMinute))

	SendACProblemsByTime(time.Now())

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		SendACProblemsByTime(time.Now())
	}
}

func SendACProblemsByTime(sendTime time.Time) {
	processDate := getProcessDate(sendTime)

	year, month, day := processDate.Date()

	SendPendingACProblems(year, int(month), day)
}

func getProcessDate(sendTime time.Time) time.Time {
	if sendTime.Hour() == 0 {
		return sendTime.AddDate(0, 0, -1)
	}

	return sendTime
}

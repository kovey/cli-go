package app

import (
	"fmt"
	"time"

	"github.com/kovey/cli-go/util"
)

var startTime time.Time

func StartTime() int64 {
	return startTime.Unix()
}

func GetRunTime() int64 {
	return time.Now().Unix() - startTime.Unix()
}

func StartTimestamp() string {
	return startTime.Format(time.DateTime)
}

func GetFormatRunTime() string {
	runTime := GetRunTime()
	days := runTime / util.Unit_Day
	runTime = runTime - days*util.Unit_Day
	hours := runTime / util.Unit_Hour
	runTime = runTime - hours*util.Unit_Hour
	mins := runTime / util.Unit_Minute
	secs := runTime - mins*util.Unit_Minute

	return fmt.Sprintf("%d days %d hours %d minutes %d seconds", days, hours, mins, secs)
}

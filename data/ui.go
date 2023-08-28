package data

import "time"

type UIStatistics struct {
	LastUpdate               int64
	NumberOfActivePrComments int
	NumberOfActivePrTasks    int
}

func ConvertActivePrResponseToUiStatistics(activePrs []PullRequest) UIStatistics {
	stats := UIStatistics{}
	for _, pr := range activePrs {
		stats.NumberOfActivePrComments += pr.Properties.CommentCount
		stats.NumberOfActivePrTasks += pr.Properties.OpenTaskCount
	}
	stats.LastUpdate = time.Now().Unix()
	return stats
}

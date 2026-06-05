package competition_service

import (
	"math"

	"GO1/models/competition_model"
)

const defaultCompetitionRating = 1000

func CalculateRank(rankingList []competition_model.CompetitionRankingListItem, ratingBefore map[int64]int) map[int64]int {
	ratingAfter := make(map[int64]int, len(rankingList))
	total := len(rankingList)
	if total == 0 {
		return ratingAfter
	}

	for start := 0; start < total; {
		end := start + 1
		for end < total && rankingList[end].CompetitionRank == rankingList[start].CompetitionRank {
			end++
		}

		delta := competitionRatingDelta(total, start, end)
		for i := start; i < end; i++ {
			before := ratingBefore[rankingList[i].UserID]
			if before == 0 {
				before = defaultCompetitionRating
			}
			ratingAfter[rankingList[i].UserID] = maxInt(0, before+delta)
		}
		start = end
	}

	return ratingAfter
}

func competitionRatingDelta(total, start, end int) int {
	expectedRank := float64(total+1) / 2
	averageRank := float64(start+1+end) / 2
	delta := int(math.Round((expectedRank - averageRank) * 10))

	return minInt(80, maxInt(-80, delta))
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

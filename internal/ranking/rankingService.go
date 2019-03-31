package ranking

type RankingService interface {
	GetRankings() (Rankings, error)
	SaveRankings(rankings Rankings) error
}

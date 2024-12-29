package repo

type Repo struct {
	DBRepo    *DBRepo
	CacheRepo *CacheRepo
}

func NewRepo(d *DBRepo, c *CacheRepo) *Repo {
	return &Repo{
		DBRepo:    d,
		CacheRepo: c,
	}
}

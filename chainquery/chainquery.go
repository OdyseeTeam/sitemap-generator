package chainquery

import (
	"database/sql"
	"fmt"
	"sync"

	"odysee-sitemap-generator/configs"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/null"
)

type CQApi struct {
	dbConn *sql.DB
	cache  *sync.Map
}

var instance *CQApi
var ClaimNotFoundErr = errors.Base("claim not found")

func Init() (*CQApi, error) {
	if instance != nil {
		return instance, nil
	}
	db, err := connect()
	if err != nil {
		return nil, err
	}
	instance = &CQApi{
		cache:  &sync.Map{},
		dbConn: db,
	}
	return instance, nil
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", configs.Configuration.Chainquery.User, configs.Configuration.Chainquery.Password, configs.Configuration.Chainquery.Host, configs.Configuration.Chainquery.Database))
	return db, errors.Err(err)
}

type Claim struct {
	ClaimID           string     `json:"claim_id"`
	Name              string     `json:"name"`
	ThumbnailURL      string     `json:"thumbnail_url"`
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	TransactionHashID string     `json:"transaction_hash_id"`
	Vout              int        `json:"vout"`
	SdHash            string     `json:"sd_hash"`
	Duration          null.Int   `json:"duration"`
	ReleaseTime       null.Int64 `json:"release_time"`
	TransactionTime   int64      `json:"transaction_time"`
}

func (c *CQApi) GetVideoStreams() ([]*Claim, error) {
	query := "SELECT claim_id, name, thumbnail_url, title, description, transaction_hash_id, vout, sd_hash, duration, release_time, transaction_time FROM claim where bid_state =? and claim_type = ? and title != ? and content_type like ? and fee_address =? and thumbnail_url != ?"
	rows, err := c.dbConn.Query(query, "controlling", 1, "", "video/%", "", "")
	if err != nil {
		return nil, errors.Err(err)
	}
	defer rows.Close()

	claims := make([]*Claim, 0, 100000)
	for rows.Next() {
		var claim Claim
		err = rows.Scan(&claim.ClaimID, &claim.Name, &claim.ThumbnailURL, &claim.Title, &claim.Description, &claim.TransactionHashID, &claim.Vout, &claim.SdHash, &claim.Duration, &claim.ReleaseTime, &claim.TransactionTime)
		if err != nil {
			return nil, errors.Err(err)
		}
		claims = append(claims, &claim)
	}
	return claims, nil
}

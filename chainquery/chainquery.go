package chainquery

import (
	"database/sql"
	"fmt"
	"sync"

	"odysee-sitemap-generator/configs"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/null"
	"github.com/sirupsen/logrus"
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
	ClaimID           string      `json:"claim_id"`
	Name              string      `json:"name"`
	ThumbnailURL      string      `json:"thumbnail_url"`
	Title             string      `json:"title"`
	Description       null.String `json:"description"`
	TransactionHashID string      `json:"transaction_hash_id"`
	Vout              int         `json:"vout"`
	SdHash            string      `json:"sd_hash"`
	Duration          null.Int    `json:"duration"`
	ReleaseTime       null.Int64  `json:"release_time"`
	TransactionTime   int64       `json:"transaction_time"`
}

func (c *CQApi) GetVideoStreams() ([]*Claim, error) {
	query := `SELECT claim_id,
       name,
       thumbnail_url,
       title,
       description,
       transaction_hash_id,
       vout,
       sd_hash,
       duration,
       release_time,
       transaction_time
FROM claim
where bid_state = ?
  and claim_type = ?
  and title != ?
  and content_type like ?
  and fee_address is null
  and thumbnail_url is not null`
	rows, err := c.dbConn.Query(query, "controlling", 1, "", "video/%")
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
	err = rows.Err()
	if err != nil {
		return nil, errors.Err(err)
	}
	return claims, nil
}

func (c *CQApi) GetVideoStreamsBatch(prevId, limit uint64) ([]*Claim, uint64, error) {
	logrus.Infof("processing batch with ids %d - %d", prevId, prevId+limit)
	query := `SELECT id, claim_id,
       name,
       thumbnail_url,
       title,
       description,
       transaction_hash_id,
       vout,
       sd_hash,
       duration,
       release_time,
       transaction_time
FROM claim
where id between ? and ?
  and bid_state = ?
  and claim_type = ?
  and title != ?
  and content_type like ?
  and fee_address is null
  and thumbnail_url is not null`
	rows, err := c.dbConn.Query(query, prevId, prevId+limit, "controlling", 1, "", "video/%")
	if err != nil {
		return nil, prevId, errors.Err(err)
	}
	defer rows.Close()

	claims := make([]*Claim, 0, limit)
	for rows.Next() {
		var claim Claim
		var lastId uint64
		err = rows.Scan(&lastId, &claim.ClaimID, &claim.Name, &claim.ThumbnailURL, &claim.Title, &claim.Description, &claim.TransactionHashID, &claim.Vout, &claim.SdHash, &claim.Duration, &claim.ReleaseTime, &claim.TransactionTime)
		if err != nil {
			return nil, prevId, errors.Err(err)
		}
		claims = append(claims, &claim)
		if lastId > prevId {
			prevId = lastId
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, prevId, errors.Err(err)
	}
	return claims, prevId, nil
}

func (c *CQApi) GetTopID() (uint64, error) {
	query := `SELECT id
FROM claim
order by id desc limit 1`
	row := c.dbConn.QueryRow(query)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return 0, errors.Err(err)
	}
	return id, nil
}

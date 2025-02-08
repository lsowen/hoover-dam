package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Credential struct {
	ID              int
	AccessKeyId     string
	CreationDate    time.Time
	SecretAccessKey string
	UserID          int64
	User            User
}

func (c Credential) OffsetKey() string {
	return c.AccessKeyId
}

type User struct {
	ID           int
	Username     string
	Email        *string
	ExternalId   *string
	FriendlyName *string
	Source       *string
	CreationDate time.Time
	Policies     []Policy `gorm:"many2many:user_policy;"`
	Groups       []Group  `gorm:"many2many:user_group;"`
	Credentials  []Credential
}

func (u User) OffsetKey() string {
	return u.Username
}

type Group struct {
	ID           int
	Name         string
	CreationDate time.Time
	Policies     []Policy `gorm:"many2many:group_policy;"`
	Users        []User   `gorm:"many2many:user_group;"`
}

func (g Group) OffsetKey() string {
	return g.Name
}

type PolicyEffect string

const (
	ALLOW PolicyEffect = "allow"
	DENY  PolicyEffect = "deny"
)

type PolicyStatement struct {
	Resource string
	Action   []string
	Effect   PolicyEffect
}

type PolicyDocument struct {
	Statement []PolicyStatement
}

func (pd PolicyDocument) Value() (driver.Value, error) {
	j, err := json.Marshal(pd)
	return j, err
}

func (pd *PolicyDocument) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &pd)
}

type Policy struct {
	ID           int
	Name         string
	Policy       PolicyDocument
	CreationDate time.Time `db:creation_date`
	Groups       []Group   `gorm:"many2many:group_policy;"`
	Users        []User    `gorm:"many2many:user_policy;"`
}

func (p Policy) OffsetKey() string {
	return p.Name
}

type PageRequest struct {
	Amount *int
	Prefix *string
	After  *string
}

func (pr PageRequest) Filter(column string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pr.Prefix != nil && *pr.Prefix != "" {
			db.Where(clause.Like{Column: column, Value: *pr.Prefix + "%"})
		}

		if pr.After != nil && *pr.After != "" {
			db.Where(clause.Gt{Column: column, Value: pr.After})
		}
		return db
	}
}

func (pr PageRequest) Limit(column string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Order(clause.OrderByColumn{Column: clause.Column{Name: column}})

		if pr.Amount != nil {
			db.Limit(*pr.Amount)
		}

		return db
	}
}

//func (pr PageRequest)

type Page struct {
	HasMore    bool
	NextOffset string
	Results    int
}

type Pageable interface {
	OffsetKey() string
}

func CalculatePage[T Pageable](db *gorm.DB, values []T) Page {
	var totalCount int64
	db.Count(&totalCount)

	if totalCount == 0 {
		return Page{
			Results:    0,
			NextOffset: "",
			HasMore:    false,
		}
	}

	count := len(values)
	hasMore := totalCount > int64(count)
	if !hasMore {
		return Page{
			Results:    count,
			NextOffset: "",
			HasMore:    hasMore,
		}
	}
	return Page{
		Results:    count,
		NextOffset: values[count-1].OffsetKey(),
		HasMore:    hasMore,
	}
}

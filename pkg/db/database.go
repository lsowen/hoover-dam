package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/lsowen/hoover-dam/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Database struct {
	orm *gorm.DB
}

func NewDatabase(ctx context.Context, cfg config.Config) (*Database, error) {

	orm, err := gorm.Open(postgres.Open(cfg.Database.Postgres.ConnectionString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "hoover_",
			SingularTable: true,
		},
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	db := &Database{
		orm: orm,
	}

	return db, nil
}

func (db *Database) GetCredential(ctx context.Context, accessKeyId string) (*Credential, error) {
	var cred Credential
	result := db.orm.Joins("User").First(&cred, "access_key_id = ?", accessKeyId)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &cred, nil
}

func (db *Database) DeleteUserCredential(ctx context.Context, username string, accessKeyId string) error {

	var cred Credential
	user, err := db.GetUser(ctx, username)
	if err != nil {
		return err
	}

	err = db.orm.Model(&user).Where("access_key_id = ?", accessKeyId).Association("Credentials").Find(&cred)
	if err != nil {
		return err
	}

	result := db.orm.Delete(&cred)
	return result.Error

}

func (db *Database) CreateUserCredential(ctx context.Context, username string, accessKeyId string, secretAccessKey string) (Credential, error) {

	cred := Credential{
		AccessKeyId:     accessKeyId,
		SecretAccessKey: secretAccessKey,
		CreationDate:    time.Now(),
	}

	user, err := db.GetUser(ctx, username)
	if err != nil {
		return cred, err
	}

	err = db.orm.Model(&user).Association("Credentials").Append(&cred)

	if err != nil {
		return cred, err
	}
	return cred, nil
}

func (db *Database) GetUser(ctx context.Context, userName string) (*User, error) {
	var user User
	result := db.orm.First(&user, "username = ?", userName)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

func (db *Database) CreateGroup(ctx context.Context, groupName string) (*Group, error) {
	group := &Group{
		Name:         groupName,
		CreationDate: time.Now(),
	}

	result := db.orm.Create(group)
	return group, result.Error
}

func (db *Database) DeleteGroup(ctx context.Context, groupName string) error {
	group, err := db.GetGroup(ctx, groupName)
	if err != nil {
		return err
	}

	return db.orm.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(group).Association("Users").Clear()
		if err != nil {
			return err
		}

		err = tx.Model(group).Association("Policies").Clear()
		if err != nil {
			return err
		}

		result := tx.Delete(group)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}

func (db *Database) GetGroup(ctx context.Context, groupName string) (*Group, error) {
	var group Group
	result := db.orm.First(&group, "name = ?", groupName)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &group, nil
}

func (db *Database) DeleteGroupMember(ctx context.Context, groupName string, username string) error {
	group, err := db.GetGroup(ctx, groupName)
	if err != nil {
		return err
	}
	user, err := db.GetUser(ctx, username)
	if err != nil {
		return err
	}

	err = db.orm.Model(group).Association("Users").Delete(user)
	return err
}

func (db *Database) AddGroupMember(ctx context.Context, groupName string, username string) error {
	group, err := db.GetGroup(ctx, groupName)
	if err != nil {
		return err
	}
	user, err := db.GetUser(ctx, username)
	if err != nil {
		return err
	}

	err = db.orm.Model(group).Association("Users").Append(user)
	return err
}

func (db *Database) DetachGroupPolicy(ctx context.Context, groupName string, policyName string) error {
	group, err := db.GetGroup(ctx, groupName)
	if err != nil {
		return err
	}
	policy, err := db.GetPolicy(ctx, policyName)
	if err != nil {
		return err
	}

	err = db.orm.Model(group).Association("Policies").Delete(policy)
	return err
}

func (db *Database) AttachGroupPolicy(ctx context.Context, groupName string, policyName string) error {
	group, err := db.GetGroup(ctx, groupName)
	if err != nil {
		return err
	}
	policy, err := db.GetPolicy(ctx, policyName)
	if err != nil {
		return err
	}

	err = db.orm.Model(group).Association("Policies").Append(policy)
	return err
}

const USER_GROUP_POLICY = `SELECT
    hoover_policy.*
  FROM
    hoover_policy
    JOIN hoover_group_policy ON hoover_policy.id = hoover_group_policy.policy_id
    JOIN hoover_user_group ON hoover_user_group.group_id = hoover_group_policy.group_id
    JOIN hoover_user ON hoover_user.id = hoover_user_group.user_id
  WHERE
   hoover_user.username = @name`

const USER_DIRECT_POLICY = `  SELECT
    hoover_policy.*
  FROM
    hoover_policy
    JOIN hoover_user_policy ON hoover_policy.id = hoover_user_policy.policy_id
    JOIN hoover_user ON hoover_user.id = hoover_user_policy.user_id
  WHERE
   hoover_user.username = @name`

func (db *Database) GetUserPolicies(ctx context.Context, username string, pr PageRequest, effective *bool) ([]Policy, Page, error) {
	var policies []Policy

	var query *gorm.DB
	if *effective {
		query = db.orm.Table(
			"(?) as u",
			db.orm.Raw(
				"? UNION ?",
				db.orm.Raw(USER_DIRECT_POLICY, sql.Named("name", username)),
				db.orm.Raw(USER_GROUP_POLICY, sql.Named("name", username)),
			),
		)
	} else {
		query = db.orm.Table(
			"(?) as u",
			db.orm.Raw(USER_DIRECT_POLICY, sql.Named("name", username)),
		)
	}
	query.Scopes(pr.Filter("name")).Session(&gorm.Session{})
	result := query.Scopes(pr.Limit("name")).Find(&policies)

	if result.Error != nil {
		return policies, Page{}, result.Error
	}

	return policies, CalculatePage(query, policies), nil
}

func (db *Database) GetGroupPolicies(ctx context.Context, groupName string, pr PageRequest) ([]Policy, Page, error) {
	var policies []Policy

	query := db.orm.
		Model(&Policy{}).
		Joins("JOIN hoover_group_policy ON hoover_group_policy.policy_id = hoover_policy.id").
		Joins("JOIN hoover_group ON hoover_group_policy.group_id = hoover_group.id").
		Where("hoover_group.name = ?", groupName).
		Scopes(pr.Filter("hoover_policy.name")).
		Session(&gorm.Session{})

	result := query.
		Scopes(pr.Limit("hoover_policy.name")).
		Find(&policies)
	if result.Error != nil {
		return policies, Page{}, result.Error
	}

	return policies, CalculatePage(query, policies), nil
}

func (db *Database) GetPolicies(ctx context.Context, pr PageRequest) ([]Policy, Page, error) {
	var policies []Policy
	query := db.orm.Model(&Policy{}).Scopes(pr.Filter("name")).Session(&gorm.Session{})
	result := query.Scopes(pr.Limit("name")).Find(&policies)
	if result.Error != nil {
		return policies, Page{}, result.Error
	}

	return policies, CalculatePage(query, policies), nil
}

func (db *Database) GetPolicy(ctx context.Context, name string) (*Policy, error) {
	var policy Policy
	result := db.orm.First(&policy, "name = ?", name)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &policy, nil
}

func (db *Database) CreatePolicy(ctx context.Context, policy *Policy) error {
	result := db.orm.Create(policy)
	return result.Error
}

func (db *Database) UpdatePolicy(ctx context.Context, policy *Policy) error {
	existingPolicy, err := db.GetPolicy(ctx, policy.Name)
	if err != nil {
		return err
	}

	result := db.orm.Model(existingPolicy).Update("policy", policy.Policy)
	return result.Error
}

func (db *Database) DeletePolicy(ctx context.Context, policyName string) error {
	policy, err := db.GetPolicy(ctx, policyName)
	if err != nil {
		return err
	}
	result := db.orm.Delete(policy)
	return result.Error
}

func (db *Database) GetGroups(ctx context.Context, pr PageRequest) ([]Group, Page, error) {
	var groups []Group
	query := db.orm.Model(&Group{}).Scopes(pr.Filter("name")).Session(&gorm.Session{})
	result := query.Scopes(pr.Limit("name")).Find(&groups)
	if result.Error != nil {
		return groups, Page{}, result.Error
	}

	return groups, CalculatePage(query, groups), nil
}

func (db *Database) GetUsers(ctx context.Context, pr PageRequest) ([]User, Page, error) {
	var users []User
	query := db.orm.Model(&User{}).Scopes(pr.Filter("username")).Session(&gorm.Session{})
	result := query.Scopes(pr.Limit("username")).Find(&users)
	if result.Error != nil {
		return users, Page{}, result.Error
	}

	return users, CalculatePage(query, users), nil
}

func (db *Database) CreateUser(ctx context.Context, user *User) error {
	result := db.orm.Create(user)
	return result.Error
}

func (db *Database) DeleteUser(ctx context.Context, username string) error {
	user, err := db.GetUser(ctx, username)
	if err != nil {
		return err
	}

	return db.orm.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(user).Association("Groups").Clear()
		if err != nil {
			return err
		}

		err = tx.Model(user).Association("Policies").Clear()
		if err != nil {
			return err
		}

		result := tx.Select("Credentials").Delete(user)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}

func (db *Database) DetachUserPolicy(ctx context.Context, username string, policyName string) error {
	user, err := db.GetUser(ctx, username)
	if err != nil {
		return err
	}
	policy, err := db.GetPolicy(ctx, policyName)
	if err != nil {
		return err
	}

	err = db.orm.Model(user).Association("Policies").Delete(policy)
	return err
}

func (db *Database) AttachUserPolicy(ctx context.Context, username string, policyName string) error {
	user, err := db.GetUser(ctx, username)
	if err != nil {
		return err
	}
	policy, err := db.GetPolicy(ctx, policyName)
	if err != nil {
		return err
	}

	err = db.orm.Model(user).Association("Policies").Append(policy)
	return err
}

func (db *Database) GetGroupMembers(ctx context.Context, pr PageRequest, groupName string) ([]User, Page, error) {
	var users []User
	sqlText := `SELECT
  hoover_user.*
FROM
  hoover_user
  JOIN hoover_user_group ON hoover_user.id = hoover_user_group.user_id
  JOIN hoover_group ON hoover_user_group.group_id = hoover_group.id
WHERE
   hoover_group.name = @name
`

	query := db.orm.
		Table(
			"(?) as u",
			db.orm.Raw(sqlText, sql.Named("name", groupName)),
		).
		Scopes(pr.Filter("username")).
		Session(&gorm.Session{})
	result := query.Scopes(pr.Limit("username")).Find(&users)
	if result.Error != nil {
		return users, Page{}, result.Error
	}

	return users, CalculatePage(query, users), nil
}

func (db *Database) GetUserGroups(ctx context.Context, pr PageRequest, username string) ([]Group, Page, error) {
	var groups []Group

	query := db.orm.
		Model(&Group{}).
		Joins("JOIN hoover_user_group ON hoover_user_group.group_id = hoover_group.id").
		Joins("JOIN hoover_user ON hoover_user_group.user_id = hoover_user.id").
		Where("username = ?", username).
		Scopes(pr.Filter("name")).
		Session(&gorm.Session{})
	result := query.Scopes(pr.Limit("name")).Find(&groups)
	if result.Error != nil {
		return groups, Page{}, result.Error
	}

	return groups, CalculatePage(query, groups), nil
}

func (db *Database) GetUserCredentials(ctx context.Context, pr PageRequest, username string) ([]Credential, Page, error) {
	var creds []Credential

	query := db.orm.
		Model(&Credential{}).
		Joins("JOIN hoover_user ON hoover_credential.user_id = hoover_user.id").
		Where("username = ?", username).
		Scopes(pr.Filter("access_key_id")).
		Session(&gorm.Session{})
	result := query.Scopes(pr.Limit("access_key_id")).Find(&creds)
	if result.Error != nil {
		return creds, Page{}, result.Error
	}

	return creds, CalculatePage(query, creds), nil
}

package api

import (
	"reflect"
	"time"

	"github.com/lsowen/hoover-dam/pkg/api/service"
	"github.com/lsowen/hoover-dam/pkg/db"
)

func ConvertCredential(cred db.Credential) service.Credentials {
	return service.Credentials{
		AccessKeyId:  cred.AccessKeyId,
		CreationDate: cred.CreationDate.Unix(),
	}
}

func ConvertCredentialWithSecret(cred db.Credential) service.CredentialsWithSecret {
	return service.CredentialsWithSecret{
		AccessKeyId:     cred.AccessKeyId,
		CreationDate:    cred.CreationDate.Unix(),
		SecretAccessKey: cred.SecretAccessKey,
		UserId:          cred.UserID,
		UserName:        &cred.User.Username,
	}
}

func ConvertUser(user db.User) service.User {
	return service.User{
		CreationDate:      user.CreationDate.Unix(),
		Email:             user.Email,
		ExternalId:        user.ExternalId,
		FriendlyName:      user.FriendlyName,
		Source:            user.Source,
		Username:          user.Username,
		EncryptedPassword: []byte{},
	}
}

func ConvertGroup(group db.Group) service.Group {
	return service.Group{
		Id:           &group.Name,
		Name:         group.Name,
		CreationDate: group.CreationDate.Unix(),
	}
}

func ConvertPolicy(policy db.Policy) service.Policy {
	statements := make([]service.Statement, len(policy.Policy.Statement))

	for idx, statement := range policy.Policy.Statement {
		statements[idx] = service.Statement{
			Action:   statement.Action,
			Effect:   string(statement.Effect),
			Resource: statement.Resource,
		}
	}

	creationDate := policy.CreationDate.Unix()
	return service.Policy{
		Name:         policy.Name,
		CreationDate: &creationDate,
		Statement:    statements,
	}
}

func ResolvePolicy(policy Policy) db.Policy {
	statements := make([]db.PolicyStatement, len(policy.Statement))

	for idx, statement := range policy.Statement {
		statements[idx] = db.PolicyStatement{
			Action:   statement.Action,
			Effect:   db.PolicyEffect(statement.Effect),
			Resource: statement.Resource,
		}
	}

	return db.Policy{
		Name:         policy.Name,
		CreationDate: time.Unix(*policy.CreationDate, 0),
		Policy: db.PolicyDocument{
			Statement: statements,
		},
	}
}

func ResolveUser(user UserCreation) db.User {
	return db.User{
		Username:     user.Username,
		ExternalId:   user.ExternalId,
		FriendlyName: user.FriendlyName,
		Source:       user.Source,
		CreationDate: time.Now(),
	}
}

func ConvertList[D interface{}, S interface{}](input []D, converter func(in D) S) []S {
	results := make([]S, len(input))
	for idx, value := range input {
		results[idx] = converter(value)
	}

	return results
}

func ConvertPagination(page db.Page) service.Pagination {
	return service.Pagination{
		HasMore:    page.HasMore,
		MaxPerPage: 1000,
		NextOffset: page.NextOffset,
		Results:    page.Results,
	}
}

func PreparePagination(input interface{}) db.PageRequest {
	inputValue := reflect.ValueOf(input)

	pr := db.PageRequest{}

	amount := inputValue.FieldByName("Amount")
	if !amount.IsNil() {
		value := (int)(amount.Elem().Interface().(service.PaginationAmount))
		pr.Amount = &value
	}

	after := inputValue.FieldByName("After")
	if !after.IsNil() {
		value := (string)(after.Elem().Interface().(service.PaginationAfter))
		pr.After = &value
	}

	prefix := inputValue.FieldByName("Prefix")
	if !prefix.IsNil() {
		value := (string)(prefix.Elem().Interface().(service.PaginationPrefix))
		pr.Prefix = &value
	}

	return pr
}

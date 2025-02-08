package api

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.5.6 -package service -generate "types,chi-server" -o service/server.gen.go ../../authorization.yml

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/lsowen/hoover-dam/pkg/api/service"
	"github.com/lsowen/hoover-dam/pkg/db"
	"github.com/treeverse/lakefs/pkg/auth/keys"
)

type APIService struct {
	Database *db.Database
}

func NewAPIService(database *db.Database) APIService {
	return APIService{
		Database: database,
	}
}

func (a APIService) GetCredentials(w http.ResponseWriter, r *http.Request, accessKeyId string) {
	cred, err := a.Database.GetCredential(r.Context(), accessKeyId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to lookup credentials: %s", err), http.StatusBadRequest)
		return
	}
	if cred == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	render.JSON(w, r, ConvertCredentialWithSecret(*cred))
}

func (a APIService) GetExternalPrincipal(w http.ResponseWriter, r *http.Request, params service.GetExternalPrincipalParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a APIService) ListGroups(w http.ResponseWriter, r *http.Request, params service.ListGroupsParams) {
	groups, page, err := a.Database.GetGroups(r.Context(), PreparePagination(params))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to lookup groups: %s", err), http.StatusBadRequest)
		return
	}

	render.JSON(w, r, service.GroupList{
		Results:    ConvertList(groups, ConvertGroup),
		Pagination: ConvertPagination(page),
	})
}

func (a APIService) CreateGroup(w http.ResponseWriter, r *http.Request) {
	payload := &GroupCreation{}
	err := render.Bind(r, payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create group: %s", err), http.StatusBadRequest)
		return
	}

	group, err := a.Database.CreateGroup(r.Context(), payload.Id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create group: %s", err), http.StatusBadRequest)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, ConvertGroup(*group))
}

func (a APIService) DeleteGroup(w http.ResponseWriter, r *http.Request, groupId string) {
	err := a.Database.DeleteGroup(r.Context(), groupId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete group: %s", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a APIService) GetGroup(w http.ResponseWriter, r *http.Request, groupId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a APIService) ListGroupMembers(w http.ResponseWriter, r *http.Request, groupId string, params service.ListGroupMembersParams) {
	users, page, err := a.Database.GetGroupMembers(r.Context(), PreparePagination(params), groupId)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to lookup group members: %s", err), http.StatusBadRequest)
		return
	}

	render.JSON(w, r, service.UserList{
		Results:    ConvertList(users, ConvertUser),
		Pagination: ConvertPagination(page),
	})
}

func (a APIService) DeleteGroupMembership(w http.ResponseWriter, r *http.Request, groupId string, userId string) {
	err := a.Database.DeleteGroupMember(r.Context(), groupId, userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete group membership: %s", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a APIService) AddGroupMembership(w http.ResponseWriter, r *http.Request, groupId string, userId string) {
	err := a.Database.AddGroupMember(r.Context(), groupId, userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add group membership: %s", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

func (a APIService) ListGroupPolicies(w http.ResponseWriter, r *http.Request, groupId string, params service.ListGroupPoliciesParams) {
	policies, page, err := a.Database.GetGroupPolicies(r.Context(), groupId, PreparePagination(params))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to lookup group policies: %s", err), http.StatusBadRequest)
		return
	}

	render.JSON(w, r, service.PolicyList{
		Results:    ConvertList(policies, ConvertPolicy),
		Pagination: ConvertPagination(page),
	})
}

func (a APIService) DetachPolicyFromGroup(w http.ResponseWriter, r *http.Request, groupId string, policyId string) {
	err := a.Database.DetachGroupPolicy(r.Context(), groupId, policyId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to detach group policy: %s", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a APIService) AttachPolicyToGroup(w http.ResponseWriter, r *http.Request, groupId string, policyId string) {
	err := a.Database.AttachGroupPolicy(r.Context(), groupId, policyId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add group policy: %s", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a APIService) ListPolicies(w http.ResponseWriter, r *http.Request, params service.ListPoliciesParams) {
	policies, page, err := a.Database.GetPolicies(r.Context(), PreparePagination(params))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to lookup policies: %s", err), http.StatusBadRequest)
		return
	}

	render.JSON(w, r, service.PolicyList{
		Results:    ConvertList(policies, ConvertPolicy),
		Pagination: ConvertPagination(page),
	})
}

func (a APIService) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	payload := &Policy{}
	err := render.Bind(r, payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create policy: %s", err), http.StatusBadRequest)
		return
	}

	policy := ResolvePolicy(*payload)
	err = a.Database.CreatePolicy(r.Context(), &policy)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create policy: %s", err), http.StatusBadRequest)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, ConvertPolicy(policy))
}

func (a APIService) DeletePolicy(w http.ResponseWriter, r *http.Request, policyId string) {
	err := a.Database.DeletePolicy(r.Context(), policyId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete policy: %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a APIService) GetPolicy(w http.ResponseWriter, r *http.Request, policyId string) {
	policy, err := a.Database.GetPolicy(r.Context(), policyId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to lookup policy: %s", err), http.StatusBadRequest)
		return
	}
	if policy == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	render.JSON(w, r, ConvertPolicy(*policy))
}

func (a APIService) UpdatePolicy(w http.ResponseWriter, r *http.Request, policyId string) {
	payload := &Policy{}
	err := render.Bind(r, payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create policy: %s", err), http.StatusBadRequest)
		return
	}

	policy := ResolvePolicy(*payload)
	policy.Name = policyId
	err = a.Database.UpdatePolicy(r.Context(), &policy)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update policy: %s", err), http.StatusBadRequest)
		return
	}

	render.JSON(w, r, ConvertPolicy(policy))
}

func (a APIService) ClaimTokenId(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a APIService) ListUsers(w http.ResponseWriter, r *http.Request, params service.ListUsersParams) {
	users, page, err := a.Database.GetUsers(r.Context(), PreparePagination(params))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to lookup users: %s", err), http.StatusBadRequest)
		return
	}

	render.JSON(w, r, service.UserList{
		Results:    ConvertList(users, ConvertUser),
		Pagination: ConvertPagination(page),
	})
}

func (a APIService) CreateUser(w http.ResponseWriter, r *http.Request) {
	payload := &UserCreation{}
	err := render.Bind(r, payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user: %s", err), http.StatusBadRequest)
		return
	}

	user := ResolveUser(*payload)
	err = a.Database.CreateUser(r.Context(), &user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user: %s", err), http.StatusBadRequest)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, ConvertUser(user))
}

func (a APIService) DeleteUser(w http.ResponseWriter, r *http.Request, userId string) {
	err := a.Database.DeleteUser(r.Context(), userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete user: %s", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a APIService) GetUser(w http.ResponseWriter, r *http.Request, userId string) {
	user, err := a.Database.GetUser(r.Context(), userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to lookup user: %s", err), http.StatusBadRequest)
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	render.JSON(w, r, ConvertUser(*user))
}

func (a APIService) ListUserCredentials(w http.ResponseWriter, r *http.Request, userId string, params service.ListUserCredentialsParams) {
	creds, page, err := a.Database.GetUserCredentials(r.Context(), PreparePagination(params), userId)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to lookup user creds: %s", err), http.StatusBadRequest)
		return
	}

	render.JSON(w, r, service.CredentialsList{
		Results:    ConvertList(creds, ConvertCredential),
		Pagination: ConvertPagination(page),
	})
}

func (a APIService) CreateCredentials(w http.ResponseWriter, r *http.Request, userId string, params service.CreateCredentialsParams) {
	accessKeyID := keys.GenAccessKeyID()
	secretAccessKey := keys.GenSecretAccessKey()

	cred, err := a.Database.CreateUserCredential(r.Context(), userId, accessKeyID, secretAccessKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user cred: %s", err), http.StatusBadRequest)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, ConvertCredentialWithSecret(cred))
}

func (a APIService) DeleteCredentials(w http.ResponseWriter, r *http.Request, userId string, accessKeyId string) {
	err := a.Database.DeleteUserCredential(r.Context(), userId, accessKeyId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete user cred: %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a APIService) GetCredentialsForUser(w http.ResponseWriter, r *http.Request, userId string, accessKeyId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a APIService) DeleteUserExternalPrincipal(w http.ResponseWriter, r *http.Request, userId string, params service.DeleteUserExternalPrincipalParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a APIService) CreateUserExternalPrincipal(w http.ResponseWriter, r *http.Request, userId string, params service.CreateUserExternalPrincipalParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a APIService) ListUserExternalPrincipals(w http.ResponseWriter, r *http.Request, userId string, params service.ListUserExternalPrincipalsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a APIService) UpdateUserFriendlyName(w http.ResponseWriter, r *http.Request, userId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a APIService) ListUserGroups(w http.ResponseWriter, r *http.Request, userId string, params service.ListUserGroupsParams) {
	groups, page, err := a.Database.GetUserGroups(r.Context(), PreparePagination(params), userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to lookup user groups: %s", err), http.StatusBadRequest)
		return
	}

	render.JSON(w, r, service.GroupList{
		Results:    ConvertList(groups, ConvertGroup),
		Pagination: ConvertPagination(page),
	})
}

func (a APIService) UpdatePassword(w http.ResponseWriter, r *http.Request, userId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a APIService) ListUserPolicies(w http.ResponseWriter, r *http.Request, userId string, params service.ListUserPoliciesParams) {
	policies, page, err := a.Database.GetUserPolicies(r.Context(), userId, PreparePagination(params), params.Effective)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to lookup user policies: %s", err), http.StatusBadRequest)
		return
	}

	render.JSON(w, r, service.PolicyList{
		Results:    ConvertList(policies, ConvertPolicy),
		Pagination: ConvertPagination(page),
	})
}

func (a APIService) DetachPolicyFromUser(w http.ResponseWriter, r *http.Request, userId string, policyId string) {
	err := a.Database.DetachUserPolicy(r.Context(), userId, policyId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to detach user policy: %s", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a APIService) AttachPolicyToUser(w http.ResponseWriter, r *http.Request, userId string, policyId string) {
	err := a.Database.AttachUserPolicy(r.Context(), userId, policyId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add user policy: %s", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a APIService) GetVersion(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, service.VersionConfig{
		Version: "something",
	})
}

func (a APIService) HealthCheck(w http.ResponseWriter, r *http.Request) {
	render.NoContent(w, r)
}

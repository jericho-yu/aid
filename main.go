package main

import (
	"fmt"
	"log"

	"github.com/jericho-yu/aid/httpClient"
)

var (
	rootUrl  = "http://127.0.0.1:30002/role-strategy/strategy"
	username = "yjz"
	token    = "11a7b66b7ca8894611ca2b7654ae6d7fe9"
)

// getAllRoles 获取所有角色
func getAllRoles() {
	client := httpClient.NewGet(fmt.Sprintf("%s/getAllRoles?type=projectRoles", rootUrl)).
		SetTimeoutSecond(5).
		SetAuthorization(username, token, "Basic")
	if client.Send().Err != nil {
		log.Fatalf("get all roles fail: %v", client.Err)
	}

	log.Printf("get all roles success: %s", client.GetResponseRawBody())
}

// addRole 增加角色
func addRole(role, pattern string) {
	client := httpClient.NewPost(fmt.Sprintf("%s/addRole", rootUrl)).
		SetTimeoutSecond(5).
		SetAuthorization(username, token, "Basic").
		SetFormBody(map[string]string{
			"type":          "projectRoles",
			"roleName":      role,
			"permissionIds": "hudson.model.Item.Discover,hudson.model.Item.Read",
			"overwrite":     "true",
			"pattern":       pattern,
		})
	if client.Send().Err != nil {
		log.Fatalf("add role fail: %v", client.Err)
	}

	log.Printf("add role success: %s", client.GetResponseRawBody())
}

// assignUserRole 分配用户到角色
func assignUserRole(role, user string) {
	client := httpClient.NewPost(fmt.Sprintf("%s/assignUserRole", rootUrl)).
		SetTimeoutSecond(5).
		SetAuthorization(username, token, "Basic").
		SetFormBody(map[string]string{
			"type":     "projectRoles",
			"roleName": role,
			"user":     user,
		})
	if client.Send().Err != nil {
		log.Fatalf("assign role fail: %v", client.Err)
	}

	log.Printf("assign role success: %s", client.GetResponseRawBody())
}

// unassignUserRole 取消用户分配到角色
func unassignUserRole(user string) {
	client := httpClient.NewPost(fmt.Sprintf("%s/unassignUserRole", rootUrl)).
		SetTimeoutSecond(5).
		SetAuthorization(username, token, "Basic").
		SetFormBody(map[string]string{
			"type":     "projectRoles",
			"roleName": "AMD",
			"user":     user,
		})
	if client.Send().Err != nil {
		log.Fatalf("assign role fail: %v", client.Err)
	}

	log.Printf("assign role success: %s", client.GetResponseRawBody())
}

// getRole 获取角色
func getRole() {
	client := httpClient.NewGet(fmt.Sprintf("%s/getRole?type=projectRoles&roleName=AMD", rootUrl)).
		SetTimeoutSecond(5).
		SetAuthorization(username, token, "Basic")
	if client.Send().Err != nil {
		log.Fatalf("get all roles fail: %v", client.Err)
	}

	log.Printf("get all roles success: %s", client.GetResponseRawBody())
}

// removeRoles 删除角色
func removeRoles(roles string) {
	client := httpClient.NewPost(fmt.Sprintf("%s/removeRoles", rootUrl)).
		SetTimeoutSecond(5).
		SetAuthorization(username, token, "Basic").
		SetFormBody(map[string]string{
			"type":      "projectRoles",
			"roleNames": roles,
		})
	if client.Send().Err != nil {
		log.Fatalf("assign role fail: %v", client.Err)
	}

	log.Printf("assign role success: %s", client.GetResponseRawBody())
}

func main() {
	// addRole("AMD1", "^AMD1$")
	// addRole("AMD2", "^AMD2$")
	// assignUserRole("AMD1", "demo1")
	// assignUserRole("AMD2", "demo2")
	removeRoles("AMD1,AMD2")
}

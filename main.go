package main

import (
	"log"

	"github.com/jericho-yu/aid/httpClient"
)

var (
	rootUrl  = "http://127.0.0.1:30002/role-strategy/strategy"
	username = "yjz"
	token    = "11a7b66b7ca8894611ca2b7654ae6d7fe9"
)

func getHttpClient(url string) *httpClient.HttpClient {
	return httpClient.NewGet(rootUrl+url).SetTimeoutSecond(5).SetAuthorization(username, token, "Basic")
}

func postHttpClient(url string) *httpClient.HttpClient {
	return httpClient.NewPost(rootUrl+url).SetTimeoutSecond(5).SetAuthorization(username, token, "Basic")
}

// getAllRoles 获取所有角色
func getAllRoles() {
	client := getHttpClient("/getAllRoles?type=projectRoles")
	if client.Send().Err != nil {
		log.Fatalf("get all roles fail: %v", client.Err)
	}

	log.Printf("get all roles success: %s", client.GetResponseRawBody())
}

// addRole 增加角色
func addRole(role, pattern string) {
	client := postHttpClient("/addRole").
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
	client := postHttpClient("/assignUserRole").
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
	client := postHttpClient("/unassignUserRole").
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
func getRole(role string) {
	client := getHttpClient("/getRole?type=projectRoles&roleName=" + role)
	if client.Send().Err != nil {
		log.Fatalf("get all roles fail: %v", client.Err)
	}

	log.Printf("get all roles success: %s", client.GetResponseRawBody())
}

// removeRoles 删除角色
func removeRoles(roles string) {
	client := postHttpClient("/removeRoles").
		SetFormBody(map[string]string{"type": "projectRoles", "roleNames": roles})
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

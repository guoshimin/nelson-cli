package main

import (
  "fmt"
  "errors"
  "strconv"
  "encoding/json"
  "github.com/parnurzeal/gorequest"
)

/////////////////// REDEPLOYMENT ///////////////////

func Redeploy(id int, http *gorequest.SuperAgent, cfg *Config) (str string, err []error){
  idAsStr := strconv.Itoa(id)
  r, body, errs := AugmentRequest(
    http.Post(cfg.Endpoint+"/v1/redeploy/"+idAsStr+""), cfg).SetDebug(false).EndBytes()

  if (r.StatusCode / 100 != 2){
    resp := string(body[:])
    errs = append(errs, errors.New("bad response from Nelson server"))
    return resp, errs
  } else {
    return "Redeployment requested.", errs
  }
}

/////////////////// LISTING STACKs ///////////////////

/**
 * {
 *   "workflow": "quasar",
 *   "guid": "67e04d28d6ab",
 *   "stack_name": "blobstore-testsuite--0-1-55--kbqg9nff",
 *   "deployed_at": 1467225866870,
 *   "unit": "blobstore-testsuite",
 *   "type": "job",
 *   "namespace": "devel"
 * }
 */
type Stack struct {
  Workflow string `json:"workflow"`
  Guid string `json:"guid"`
  StackName string `json:"stack_name"`
  DeployedAt int64 `json:"deployed_at"`
  UnitName string `json:"unit"`
  Type string `json:"type"`
  NamespaceRef string `json:"namespace"`
}

func ListStacks(dc string, http *gorequest.SuperAgent, cfg *Config) (list []Stack, err []error){
  r, bytes, errs := AugmentRequest(
    http.Get(cfg.Endpoint+"/v1/datacenters/"+dc+"/deployments"), cfg).EndBytes()

  if (r.StatusCode / 100 != 2){
    errs = append(errs, errors.New("Bad response from Nelson server"))
    return nil, errs
  } else {
    var list []Stack
    if err := json.Unmarshal(bytes, &list); err != nil {
      panic(err)
    }
    return list, errs
  }
}

func PrintListStacks(stacks []Stack){
  var tabulized = [][]string {}
  for _,s := range stacks {
    tabulized = append(tabulized,[]string{ s.Guid, s.NamespaceRef, s.StackName, s.Type, s.Workflow, JavaEpochToDateStr(s.DeployedAt) })
  }

  RenderTableToStdout([]string{ "GUID", "Namespace", "Stack", "Type", "Workflow", "Deployed At" }, tabulized)
}

/////////////////// DEPLOYMENT LOG ///////////////////

type StackLog struct {
  Content []string `json:"content"`
  Offset int `json:"offset"`
}

// v1/deployments/:id/log
func GetDeploymentLog(id int, http *gorequest.SuperAgent, cfg *Config){
  idAsStr := strconv.Itoa(id)
  _, bytes, errs := AugmentRequest(
    http.Get(cfg.Endpoint+"/v1/deployments/"+idAsStr+"/log"), cfg).EndBytes()

  if (len(errs) > 0) {
    panic(errs)
  }

  var logs StackLog
  if err := json.Unmarshal(bytes, &logs); err != nil {
    panic(err)
  }

  fmt.Println("===>> logs for deployment "+ idAsStr)

  for _,l := range logs.Content {
    fmt.Println(l)
  }
}
package wigo

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"strconv"
	"fmt"
)

func HttpRemotesHandler(params martini.Params) (int, string) {

	hostname := params["hostname"]

	if hostname != "" {
		remoteWigo := GetLocalWigo().FindRemoteWigoByHostname(hostname)
		if remoteWigo != nil {
			json, err := remoteWigo.ToJsonString()
			if err != nil {
				return 500, "Failed to encode remote wigo"
			} else {
				return 200, json
			}
		} else {
			return 404, ""
		}
	}

	// Return remotes list
	list := GetLocalWigo().ListRemoteWigosNames()
	json, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		return 500, ""
	} else {
		return 200, string(json)
	}
}

func HttpRemotesProbesHandler(params martini.Params) (int, string) {

	hostname := params["hostname"]
	probe := params["probe"]

	if hostname == "" {
		return 404, "No wigo name set in url"
	}

	// Get remote wigo
	remoteWigo := GetLocalWigo().FindRemoteWigoByHostname(hostname)
	if remoteWigo == nil {
		return 404, "Remote wigo " + hostname + " not found"
	}

	// Get probe or probes
	if probe != "" {
		if remoteWigo.LocalHost.Probes[probe] != nil {
			json, err := json.MarshalIndent(remoteWigo.LocalHost.Probes[probe], "", "    ")
			if err != nil {
				return 500, ""
			} else {
				return 200, string(json)
			}
		}
	} else {
		json, err := json.MarshalIndent(remoteWigo.ListProbes(), "", "    ")
		if err != nil {
			return 500, ""
		} else {
			return 200, string(json)
		}
	}

	return 200, ""
}

func HttpRemotesStatusHandler(params martini.Params) (int, string) {

	hostname := params["hostname"]

	if hostname == "" {
		return 404, "No wigo name set in url"
	}

	// Get remote wigo
	remoteWigo := GetLocalWigo().FindRemoteWigoByHostname(hostname)
	if remoteWigo == nil {
		return 404, "Remote wigo " + hostname + " not found"
	}

	return 200, strconv.Itoa(remoteWigo.GlobalStatus)
}

func HttpRemotesProbesStatusHandler(params martini.Params) (int, string) {

	hostname := params["hostname"]
	probe := params["probe"]

	if hostname == "" {
		return 404, "No wigo name set in url"
	}
	if probe == "" {
		return 404, "No probe name set in url"
	}

	// Get remote wigo
	remoteWigo := GetLocalWigo().FindRemoteWigoByHostname(hostname)
	if remoteWigo == nil {
		return 404, "Remote wigo " + hostname + " not found"
	}

	// Get probe
	if remoteWigo.LocalHost.Probes[probe] == nil {
		return 404, "Probe " + probe + " not found in remote wigo " + hostname
	}

	return 200, strconv.Itoa(remoteWigo.LocalHost.Probes[probe].Status)
}


func HttpLogsHandler(params martini.Params) (int, string) {

	hostname := params["hostname"]
	probe := params["probe"]
	group := params["group"]

	// Test hostname if present
	var remoteWigo *Wigo
	if hostname != "" {
		remoteWigo = GetLocalWigo().FindRemoteWigoByHostname(hostname)
		if remoteWigo == nil {
			return 404, "Remote wigo "+hostname+" not found"
		}
	}

	// Test probe
	if probe != "" {
		if hostname != "" {
			// Get probe
			if remoteWigo.LocalHost.Probes[probe] == nil {
				return 404, "Probe " + probe + " not found in remote wigo " + hostname
			}
		}
	}

	// Get logs
	logs := LocalWigo.SearchLogs(probe,hostname,group)

	// Json
	json, err := json.MarshalIndent(logs, "", "    ")
	if err != nil {
		return 500, ""
	}

	return 200, string(json)
}


func HttpGroupsHandler(params martini.Params) (int, string) {

	group := params["group"]

	result := make(map[string] interface {})
	result["Name"] = group

	if group != "" {
		gs, status := GetLocalWigo().GroupSummary(group)
		if gs != nil {

			result["Status"] = status
			result["Hosts"]  = gs

			json, err := json.MarshalIndent(result, "", "    ")
			if err != nil {
				return 500, fmt.Sprintf("Fail to encode summary : %s" ,err)
			} else {
				return 200, string(json)
			}
		} else {
			return 404, ""
		}
	}

	// Return remotes list
	list := GetLocalWigo().ListGroupsNames()
	json, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		return 500, ""
	} else {
		return 200, string(json)
	}
}

func HttpLogsIndexesHandler(params martini.Params) (int, string) {

	result := make(map[string][]string)
	result["probes"] 	= make([]string,0)
	result["hosts"] 	= make([]string,0)
	result["groups"] 	= make([]string,0)

	// Probes
	for probeName := range LocalWigo.logsProbeIndex {
		result["probes"] = append(result["probes"], probeName)
	}

	// Hosts
	for hostName := range LocalWigo.logsWigoIndex {
		result["hosts"] = append(result["hosts"], hostName)
	}

	// Groups
	for groupName := range LocalWigo.logsGroupIndex {
		result["groups"] = append(result["groups"], groupName)
	}

	// Return remotes list
	json, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return 500, fmt.Sprintf("Error while encoding to json : %s", err)
	} else {
		return 200, string(json)
	}
}

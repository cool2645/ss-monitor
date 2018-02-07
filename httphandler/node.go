package httphandler

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/cool2645/ss-monitor/model"
	"github.com/yanzay/log"
	"strconv"
	"github.com/cool2645/ss-monitor/manager"
	"encoding/json"
)

func NewNode(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !authAdmin(w, req) {
		return
	}
	var node model.Node
	if req.Header.Get("Content-Type") == "application/json" {
		err := json.NewDecoder(req.Body).Decode(&node)
		if err != nil {
			log.Error(err)
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Error occurred parsing json request.",
			}
			responseJson(w, res, http.StatusBadRequest)
			return
		}
	} else {
		req.ParseForm()
		if len(req.Form["name"]) != 1 {
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Invalid node name.",
			}
			responseJson(w, res, http.StatusBadRequest)
			return
		}
		node.Name = req.Form["name"][0]
		if len(req.Form["ipv4"]) == 1 {
			node.IPv4 = req.Form["ipv4"][0]
		}
		if len(req.Form["ipv6"]) == 1 {
			node.IPv6 = req.Form["ipv6"][0]
		}
		if len(req.Form["ss4_json"]) == 1 {
			node.Ss4Json = req.Form["ss4_json"][0]
			node.EnableIPv4Testing = true
		}
		if len(req.Form["ss6_json"]) == 1 {
			node.Ss6Json = req.Form["ss6_json"][0]
			node.EnableIPv6Testing = true
		}
		if len(req.Form["domain_prefix4"]) == 1 {
			node.DomainPrefix4 = req.Form["domain_prefix4"][0]
		}
		if len(req.Form["domain_prefix6"]) == 1 {
			node.DomainPrefix6 = req.Form["domain_prefix6"][0]
		}
		if len(req.Form["domain_root"]) == 1 {
			node.DomainRoot = req.Form["domain_root"][0]
			node.EnableWatching = true
		}
		if len(req.Form["provider"]) == 1 {
			node.Provider = req.Form["provider"][0]
			node.EnableCleaning = true
		}
		if len(req.Form["dns_provider"]) == 1 {
			node.DNSProvider = req.Form["dns_provider"][0]
		}
		if len(req.Form["os"]) == 1 {
			node.OS = req.Form["os"][0]
		}
		if len(req.Form["image"]) == 1 {
			node.Image = req.Form["image"][0]
		}
		if len(req.Form["data_center"]) == 1 {
			node.DataCenter = req.Form["data_center"][0]
		}
		if len(req.Form["plan"]) == 1 {
			node.Plan = req.Form["plan"][0]
		}
	}
	node, err := model.CreateNode(model.Db, node)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred creating node: " + err.Error(),
		}
		responseJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   node,
	}
	responseJson(w, res, http.StatusOK)
	manager.InitNodes()
}

func EditNode(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !auth(w, req) {
		return
	}
	var node model.Node
	if req.Header.Get("Content-Type") == "application/json" {
		err := json.NewDecoder(req.Body).Decode(&node)
		if err != nil {
			log.Error(err)
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Error occurred parsing json request.",
			}
			responseJson(w, res, http.StatusBadRequest)
			return
		}
		node, err = model.UpdateNodeAllFields(model.Db, node)
		if err != nil {
			log.Error(err)
			if err.Error() == "UpdateNodeAllFields: record not found" {
				res := map[string]interface{}{
					"code":   http.StatusNotFound,
					"result": false,
					"msg":    "Error occurred updating node: " + err.Error(),
				}
				responseJson(w, res, http.StatusNotFound)
				return
			}
			res := map[string]interface{}{
				"code":   http.StatusInternalServerError,
				"result": false,
				"msg":    "Error occurred updating node: " + err.Error(),
			}
			responseJson(w, res, http.StatusInternalServerError)
			return
		}
	} else {
		req.ParseForm()
		nodeID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
		if err != nil {
			log.Error(err)
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Error occurred parsing node id.",
			}
			responseJson(w, res, http.StatusBadRequest)
			return
		}
		node.ID = uint(nodeID64)
		if len(req.Form["name"]) == 1 {
			node.Name = req.Form["name"][0]
		}
		if len(req.Form["ipv4"]) == 1 {
			node.IPv4 = req.Form["ipv4"][0]
		}
		if len(req.Form["ipv6"]) == 1 {
			node.IPv6 = req.Form["ipv6"][0]
		}
		if len(req.Form["ss4_json"]) == 1 {
			node.Ss4Json = req.Form["ss4_json"][0]
		}
		if len(req.Form["ss6_json"]) == 1 {
			node.Ss6Json = req.Form["ss6_json"][0]
		}
		if len(req.Form["domain_prefix4"]) == 1 {
			node.DomainPrefix4 = req.Form["domain_prefix4"][0]
		}
		if len(req.Form["domain_prefix6"]) == 1 {
			node.DomainPrefix6 = req.Form["domain_prefix6"][0]
		}
		if len(req.Form["domain_root"]) == 1 {
			node.DomainRoot = req.Form["domain_root"][0]
		}
		if len(req.Form["provider"]) == 1 {
			node.Provider = req.Form["provider"][0]
		}
		if len(req.Form["dns_provider"]) == 1 {
			node.DNSProvider = req.Form["dns_provider"][0]
		}
		if len(req.Form["os"]) == 1 {
			node.OS = req.Form["os"][0]
		}
		if len(req.Form["image"]) == 1 {
			node.Image = req.Form["image"][0]
		}
		if len(req.Form["data_center"]) == 1 {
			node.DataCenter = req.Form["data_center"][0]
		}
		if len(req.Form["plan"]) == 1 {
			node.Plan = req.Form["plan"][0]
		}
		var fields = make(map[string]interface{})
		if len(req.Form["enable_watching"]) == 1 {
			node.EnableWatching, err = strconv.ParseBool(req.Form["enable_watching"][0])
			fields["EnableWatching"] = node.EnableWatching
			if err != nil {
				log.Error(err)
				res := map[string]interface{}{
					"code":   http.StatusBadRequest,
					"result": false,
					"msg":    "Error occurred parsing enable watching.",
				}
				responseJson(w, res, http.StatusBadRequest)
				return
			}
		}
		if len(req.Form["enable_ipv4_testing"]) == 1 {
			node.EnableIPv4Testing, err = strconv.ParseBool(req.Form["enable_ipv4_testing"][0])
			fields["EnableIPv4Testing"] = node.EnableIPv4Testing
			if err != nil {
				log.Error(err)
				res := map[string]interface{}{
					"code":   http.StatusBadRequest,
					"result": false,
					"msg":    "Error occurred parsing enable ipv4 testing.",
				}
				responseJson(w, res, http.StatusBadRequest)
				return
			}
		}
		if len(req.Form["enable_ipv6_testing"]) == 1 {
			node.EnableIPv6Testing, err = strconv.ParseBool(req.Form["enable_ipv6_testing"][0])
			fields["EnableIPv6Testing"] = node.EnableIPv6Testing
			if err != nil {
				log.Error(err)
				res := map[string]interface{}{
					"code":   http.StatusBadRequest,
					"result": false,
					"msg":    "Error occurred parsing enable ipv4 testing.",
				}
				responseJson(w, res, http.StatusBadRequest)
				return
			}
		}
		if len(req.Form["enable_cleaning"]) == 1 {
			node.EnableCleaning, err = strconv.ParseBool(req.Form["enable_cleaning"][0])
			fields["EnableCleaning"] = node.EnableCleaning
			if err != nil {
				log.Error(err)
				res := map[string]interface{}{
					"code":   http.StatusBadRequest,
					"result": false,
					"msg":    "Error occurred parsing enable cleaning.",
				}
				responseJson(w, res, http.StatusBadRequest)
				return
			}
		}
		node, err = model.UpdateNodeChangedFields(model.Db, node)
		if err != nil {
			log.Error(err)
			if err.Error() == "UpdateNode: record not found" {
				res := map[string]interface{}{
					"code":   http.StatusNotFound,
					"result": false,
					"msg":    "Error occurred updating node: " + err.Error(),
				}
				responseJson(w, res, http.StatusNotFound)
				return
			}
			res := map[string]interface{}{
				"code":   http.StatusInternalServerError,
				"result": false,
				"msg":    "Error occurred updating node: " + err.Error(),
			}
			responseJson(w, res, http.StatusInternalServerError)
			return
		}
		if len(fields) > 0 {
			err = model.UpdateNodeFields(model.Db, node.ID, fields)
			if err != nil {
				log.Error(err)
				if err.Error() == "UpdateNodeFields: Find node: record not found" {
					res := map[string]interface{}{
						"code":   http.StatusNotFound,
						"result": false,
						"msg":    "Error occurred updating node control bit: " + err.Error(),
					}
					responseJson(w, res, http.StatusNotFound)
					return
				}
				res := map[string]interface{}{
					"code":   http.StatusInternalServerError,
					"result": false,
					"msg":    "Error occurred updating node control bit: " + err.Error(),
				}
				responseJson(w, res, http.StatusInternalServerError)
				return
			}
		}
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   node,
	}
	responseJson(w, res, http.StatusOK)
	manager.InitNodes()
}

func SetNodeCleaning(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !authAdmin(w, req) {
		return
	}
	nodeID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred parsing node id.",
		}
		responseJson(w, res, http.StatusBadRequest)
		return
	}
	nodeID := uint(nodeID64)
	err = model.SetNodeCleaning(model.Db, nodeID)
	if err != nil {
		log.Error(err)
		if err.Error() == "SetNodeCleaning: Find node: record not found" {
			res := map[string]interface{}{
				"code":   http.StatusNotFound,
				"result": false,
				"msg":    "Error occurred setting node cleaning: " + err.Error(),
			}
			responseJson(w, res, http.StatusNotFound)
			return
		}
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred setting node cleaning: " + err.Error(),
		}
		responseJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"msg":    "success",
	}
	responseJson(w, res, http.StatusOK)
	manager.InitNodes()
}

func ResetNode(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !authAdmin(w, req) {
		return
	}
	nodeID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred parsing node id.",
		}
		responseJson(w, res, http.StatusBadRequest)
		return
	}
	nodeID := uint(nodeID64)
	err = model.ResetNode(model.Db, nodeID)
	if err != nil {
		log.Error(err)
		if err.Error() == "ResetNode: Find node: record not found" {
			res := map[string]interface{}{
				"code":   http.StatusNotFound,
				"result": false,
				"msg":    "Error occurred resetting node: " + err.Error(),
			}
			responseJson(w, res, http.StatusNotFound)
			return
		}
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred resetting node: " + err.Error(),
		}
		responseJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"msg":    "success",
	}
	responseJson(w, res, http.StatusOK)
	manager.InitNodes()
}

func DeleteNode(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !authAdmin(w, req) {
		return
	}
	nodeID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred parsing node id.",
		}
		responseJson(w, res, http.StatusBadRequest)
		return
	}
	nodeID := uint(nodeID64)
	err = model.DeleteNode(model.Db, nodeID)
	if err != nil {
		log.Error(err)
		if err.Error() == "DeleteNode: record not found" {
			res := map[string]interface{}{
				"code":   http.StatusNotFound,
				"result": false,
				"msg":    "Error occurred deleting node: " + err.Error(),
			}
			responseJson(w, res, http.StatusNotFound)
			return
		}
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred deleting node: " + err.Error(),
		}
		responseJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"msg":    "success",
	}
	responseJson(w, res, http.StatusOK)
	manager.InitNodes()
}

func GetNodes(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	nodes, err := model.GetNodes(model.Db)
	if err != nil {
		log.Error(err)
		if err.Error() == "GetNodes: sql: no rows in result set" {
			res := map[string]interface{}{
				"code":   http.StatusNotFound,
				"result": false,
				"msg":    "Error occurred querying nodes: " + err.Error(),
			}
			responseJson(w, res, http.StatusNotFound)
			return
		}
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred querying nodes: " + err.Error(),
		}
		responseJson(w, res, http.StatusInternalServerError)
		return
	}
	if !checkAccessKey(req) && !checkAdmin(w, req) {
		for i, _ := range nodes {
			nodes[i].Ss4Json = ""
			nodes[i].Ss6Json = ""
		}
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   nodes,
	}
	responseJson(w, res, http.StatusOK)
}

func GetNode(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	nodeID64, err := strconv.ParseUint(ps.ByName("id"), 10, 32)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred parsing node id.",
		}
		responseJson(w, res, http.StatusBadRequest)
		return
	}
	nodeID := uint(nodeID64)
	node, err := model.GetNode(model.Db, nodeID)
	if err != nil {
		log.Error(err)
		if err.Error() == "GetNode: record not found" {
			res := map[string]interface{}{
				"code":   http.StatusNotFound,
				"result": false,
				"msg":    "Error occurred querying nodes: " + err.Error(),
			}
			responseJson(w, res, http.StatusNotFound)
			return
		}
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred querying nodes: " + err.Error(),
		}
		responseJson(w, res, http.StatusInternalServerError)
		return
	}
	if !checkAccessKey(req) && !checkAdmin(w, req) {
		node.Ss4Json = ""
		node.Ss6Json = ""
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   node,
	}
	responseJson(w, res, http.StatusOK)
}

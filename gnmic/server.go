package gnmic

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	yson "github.com/ghodss/yaml"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/orb-community/orb-gnmic/config"
	"github.com/orb-community/orb-gnmic/runner"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type ReturnPolicyData struct {
	State runner.State `yaml:"status"`
	config.Policy
}

type ReturnValue struct {
	Message string `json:"message"`
}

func (o *Gnmic) startServer() error {
	gin.SetMode(gin.ReleaseMode)
	o.router = gin.New()

	o.router.Use(ginzap.Ginzap(o.logger, time.RFC3339, true))
	o.router.Use(ginzap.RecoveryWithZap(o.logger, true))

	// Routes
	o.router.GET("/api/v1/status", o.getStatus)
	o.router.GET("/api/v1/capabilities", o.getCapabilities)
	o.router.GET("/api/v1/policies", o.getPolicies)
	o.router.POST("/api/v1/policies", o.createPolicy)
	o.router.GET("/api/v1/policies/:policy", o.getPolicy)
	o.router.DELETE("/api/v1/policies/:policy", o.deletePolicy)

	serverHost := o.conf.ServerHost
	serverPort := strconv.FormatUint(o.conf.ServerPort, 10)
	go func() {
		serv := serverHost + ":" + serverPort
		o.logger.Info("starting gnmic server at: " + serv)
		if err := o.router.Run(serv); err != nil {
			o.logger.Fatal("shutting down the server", zap.Error(err))
		}
	}()
	return nil
}

func (o *Gnmic) getStatus(c *gin.Context) {
	o.stat.UpTime = time.Since(o.stat.StartTime)
	c.IndentedJSON(http.StatusOK, o.stat)
}

func (o *Gnmic) getCapabilities(c *gin.Context) {
	j, err := yson.YAMLToJSON(o.capabilities)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ReturnValue{err.Error()})
		return
	}
	var ret interface{}
	err = json.Unmarshal(j, &ret)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ReturnValue{err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, ret)
}

func (o *Gnmic) getPolicies(c *gin.Context) {
	policies := make([]string, 0, len(o.policies))
	for k := range o.policies {
		policies = append(policies, k)
	}
	c.IndentedJSON(http.StatusOK, policies)
}

func (o *Gnmic) getPolicy(c *gin.Context) {
	policy := c.Param("policy")
	rInfo, ok := o.policies[policy]
	if ok {
		c.YAML(http.StatusOK, map[string]ReturnPolicyData{policy: {rInfo.Instance.GetStatus(), rInfo.Policy}})
	} else {
		c.IndentedJSON(http.StatusNotFound, ReturnValue{"policy not found"})
	}
}

func (o *Gnmic) createPolicy(c *gin.Context) {
	if t := c.Request.Header.Get("Content-type"); t != "application/x-yaml" {
		c.IndentedJSON(http.StatusBadRequest, ReturnValue{"invalid Content-Type. Only 'application/x-yaml' is supported"})
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ReturnValue{err.Error()})
		return
	}
	var payload map[string]config.Policy
	err = yaml.Unmarshal(body, &payload)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ReturnValue{err.Error()})
		return
	}
	if len(payload) > 1 {
		c.IndentedJSON(http.StatusBadRequest, ReturnValue{"only single policy allowed per request"})
		return
	}
	var policy string
	var data config.Policy
	for policy, data = range payload {
		_, ok := o.policies[policy]
		if ok {
			c.IndentedJSON(http.StatusConflict, ReturnValue{"policy already exists"})
			return

		}
		if len(data.Config.Subscriptions) == 0 {
			c.IndentedJSON(http.StatusForbidden, ReturnValue{"subscriptions field is required"})
			return
		}
		if len(data.Config.Targets) == 0 {
			c.IndentedJSON(http.StatusForbidden, ReturnValue{"targets field is required"})
			return
		}
	}

	r := runner.New(o.logger, policy, o.policiesDir)
	if err := r.Configure(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, ReturnValue{err.Error()})
		return
	}
	runnerCtx := context.WithValue(o.ctx, "routine", policy)
	if err := r.Start(context.WithCancel(runnerCtx)); err != nil {
		c.IndentedJSON(http.StatusBadRequest, ReturnValue{err.Error()})
		return
	}
	o.policies[policy] = RunnerInfo{Policy: data, Instance: r}
	c.YAML(http.StatusCreated, map[string]ReturnPolicyData{policy: {r.GetStatus(), data}})
}

func (o *Gnmic) deletePolicy(c *gin.Context) {
	policy := c.Param("policy")
	r, ok := o.policies[policy]
	if ok {
		if err := r.Instance.Stop(o.ctx); err != nil {
			c.IndentedJSON(http.StatusBadRequest, ReturnValue{err.Error()})
			return
		}
		delete(o.policies, policy)
		c.IndentedJSON(http.StatusOK, ReturnValue{policy + " was deleted"})
	} else {
		c.IndentedJSON(http.StatusNotFound, ReturnValue{"policy not found"})
	}
}

package gnmic

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/orb-community/orb-gnmic/config"
	"github.com/orb-community/orb-gnmic/runner"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type RunnerInfo struct {
	Policy   config.Policy
	Instance runner.Runner
}

type Gnmic struct {
	logger         *zap.Logger
	conf           *config.Config
	stat           config.Status
	policies       map[string]RunnerInfo
	policiesDir    string
	ctx            context.Context
	cancelFunction context.CancelFunc
	router         *gin.Engine
	capabilities   []byte
}

func New(logger *zap.Logger, c *config.Config) (Gnmic, error) {
	return Gnmic{logger: logger, conf: c, policies: make(map[string]RunnerInfo)}, nil
}

func (o *Gnmic) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	o.stat.StartTime = time.Now()
	o.ctx = context.WithValue(ctx, "routine", "orbGnmicRoutine")
	o.cancelFunction = cancelFunc

	var err error
	o.policiesDir, err = os.MkdirTemp("", "policies")
	if err != nil {
		return err
	}
	o.capabilities, err = runner.GetCapabilities()
	if err != nil {
		return err
	}
	s := struct {
		Version string
	}{}
	err = yaml.Unmarshal(o.capabilities, &s)
	if err != nil {
		return err
	}
	o.stat.Version = s.Version
	if err = o.startServer(); err != nil {
		return err
	}
	return nil
}

func (o *Gnmic) Stop(ctx context.Context) {
	o.logger.Info("routine call for stop orb-gnmic", zap.Any("routine", ctx.Value("routine")))
	defer os.RemoveAll(o.policiesDir)
	o.cancelFunction()
}

func (o *Gnmic) RestartRunner(ctx context.Context, name string, reason string) error {
	return nil
}

func (o *Gnmic) RestartAll(ctx context.Context, reason string) error {
	return nil
}

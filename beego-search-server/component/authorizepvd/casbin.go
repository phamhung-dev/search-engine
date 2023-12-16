package authorizepvd

import (
	"beego-search-server/models"
	"fmt"
	"os"

	beeLogger "github.com/beego/bee/v2/logger"
	"github.com/casbin/casbin"
)

type casbinProvider struct {
	enforcer *casbin.Enforcer
}

func NewCasbinProvider() *casbinProvider {
	casbinConf := os.Getenv("CASBIN_CONF")
	casbinPolicy := os.Getenv("CASBIN_POLICY")
	if casbinConf == "" || casbinPolicy == "" {
		beeLogger.Log.Fatal(ErrProviderIsNotConfigured.Error())
	}

	enforcer := casbin.NewEnforcer(
		fmt.Sprintf("./component/authorizepvd/%s", casbinConf),
		fmt.Sprintf("./component/authorizepvd/%s", casbinPolicy),
	)

	return &casbinProvider{enforcer: enforcer}
}

func (provider *casbinProvider) ValidateRequest(user *models.User, path string, method string) (bool, error) {
	result, err := provider.enforcer.EnforceSafe(user.Role, path, method)

	if err != nil {
		return false, err
	}

	return result, nil
}

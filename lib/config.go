package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	KubeStaging    string `json:"kube_staging"`
	KubeProduction string `json:"kube_production"`
	KubeConfigPath string `json:"kube_config_path"`
}

var (
	UserHome        string
	ConfigPath      string
	ConfigDirectory string
	TmpDirectory    string
)

var ValidInspectComponents = map[string]string{"service": "kubernetes service"}
var ValidListComponents = map[string]string{"services": "kubernetes services"}
var ValidEnvironments = map[string]string{
	"staging":    "Staging environment",
	"production": "Production environment",
}
var RequiredDeploymentFiles = map[string][]string{
	BACKEND:  []string{"deployment.yaml.jinja2.tpl", "service.yaml.jinja2.tpl", BACKEND_DEPLOYMENT_GENERATION_COMMAND},
	FRONTEND: []string{"deployment.yaml.jinja2.tpl", "service.yaml.jinja2.tpl", FRONTEND_DEPLOYMENT_GENERATION_COMMAND},
}
var DeploymentMapping = map[string]string{BACKEND: "/", FRONTEND: "apps/"}

const APP_CONFIG_DIRECTORY = "/.toolchain/"
const APP_CONFIG_NAME = "config"
const APP_CONFIG_EXTENSION = "json"
const KUBECTL_CONFIG_DEFAULT = "/.kube/config" // could use contextOption.configAccess.GetDefaultFilename()
const KUBECTL_PRODUCTION_DEFAULT = "gke_andela-kube_us-east1-b_andela-prod"
const KUBECTL_STAGING_DEFAULT = "gke_microservices-kube_us-east1-c_staging"
const KUBECTL_CONFIG_PATH_KEY = "kube_config_path"
const GIT_ACCESS_KEY_STRING = "GIT_ACCESS_TOKEN"
const DEPLOYMENT_ENVIRONMENT_FOLDERNAME = "environments_fields/"
const DEPLOYMENT_SCRIPT_BRANCH = "develop"
const GITHUB_USER = "andela"
const DEPLOYMENT_SCRIPT_GITHUB_REPOSITORY = "micro-deployment-scripts"
const FRONTEND = "frontend"
const BACKEND = "backend"
const FRONTEND_DEPLOYMENT_GENERATION_COMMAND = "gen_apps.sh"
const BACKEND_DEPLOYMENT_GENERATION_COMMAND = "gen_services.sh"
const DEPLOYMENT_CONFIGURATION_ARCHIVE_BUCKET = "andela-deployment-configuration-archive"

func init() {
	UserHome = os.ExpandEnv("$HOME")
	ConfigDirectory = UserHome + APP_CONFIG_DIRECTORY
	TmpDirectory = UserHome + APP_CONFIG_DIRECTORY + "tmp/"
	ConfigPath = ConfigDirectory + APP_CONFIG_NAME + "." + APP_CONFIG_EXTENSION
	if FileExists(ConfigDirectory) == false {
		os.MkdirAll(ConfigDirectory, os.ModePerm)
	}
}

func (object *Config) SetDefaults() {
	object.KubeConfigPath = UserHome + KUBECTL_CONFIG_DEFAULT
	object.KubeProduction = KUBECTL_PRODUCTION_DEFAULT
	object.KubeStaging = KUBECTL_STAGING_DEFAULT
}

func (object *Config) Save() {
	data, err := json.MarshalIndent(*object, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(ConfigPath, data, 0644)
	if err != nil {
		panic(err)
	}
}

func NewConfigFromFile() *Config {
	config := Config{}
	b, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &config)
	if err != nil {
		panic(err)
	}
	return &config
}

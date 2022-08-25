package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
)

type ConfigMap struct {
	Name      string `conf:"default:nginx-load-balancer-microk8s-conf,env:NGINX_CONFIG_MAP_NAME"`
	Namespace string `conf:"default:ingress,env:NGINX_CONFIG_MAP_NAMESPACE"`
}
type Config struct {
	NginxConfigName string   `conf:"default:whitelist-source-range,env:NGINX_ANNOTATION"`
	SourceURIs      []string `conf:"default:https://www.cloudflare.com/ips-v4;https://www.cloudflare.com/ips-v6,env:NGINX_SOURCE_URIs"`
	ConfigMap
	conf.Version
}

const namespace = "INGRESS_WHITELIST"

func New(svn, desc string) (*Config, error) {
	if svn == "" || desc == "" {
		return nil, fmt.Errorf("missing version or description (both are required)")
	}

	cfg := Config{}
	cfg.Version.SVN = svn
	cfg.Version.Desc = desc

	if err := conf.Parse(os.Args[1:], namespace, &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage(namespace, &cfg)
			if err != nil {
				return nil, errors.Wrap(err, "generating config usage")
			}
			log.Panic(usage)

		case conf.ErrVersionWanted:
			version, err := conf.VersionString(namespace, &cfg)
			if err != nil {
				return nil, errors.Wrap(err, "generating config version")
			}
			log.Panic(version)
		}
		return nil, errors.Wrap(err, "parsing config")
	}

	return &cfg, nil
}

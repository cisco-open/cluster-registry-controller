// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"emperror.dev/errors"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"

	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/internal/config"
)

type Configuration config.Configuration

func configure() Configuration {
	p := flag.NewFlagSet(FriendlyServiceName, flag.ExitOnError)
	initConfiguration(viper.GetViper(), p)
	err := p.Parse(os.Args[1:])
	if err != nil {
		panic(errors.WrapIf(err, "failed to parse arguments"))
	}

	var config Configuration
	bindEnvs(config)
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(errors.WrapIf(err, "failed to unmarshal configuration"))
	}

	// Show version if asked for
	if viper.GetBool("version") {
		fmt.Printf("%s version %s (%s) built on %s\n", FriendlyServiceName, version, commitHash, buildDate)
		os.Exit(0)
	}

	// Dump config if asked for
	if viper.GetBool("dump-config") {
		t, err := yaml.Marshal(config)
		if err != nil {
			panic(errors.WrapIf(err, "failed to dump configuration"))
		}
		fmt.Print(string(t))
		os.Exit(0)
	}

	return config
}

func initConfiguration(v *viper.Viper, p *flag.FlagSet) {
	v.AllowEmptyEnv(true)
	p.Init(FriendlyServiceName, flag.ExitOnError)
	p.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", FriendlyServiceName)
		p.PrintDefaults()
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	p.String("metrics-addr", ":8080", "The address the metric endpoint binds to.")
	p.Bool("devel-mode", false, "Set development mode (mainly for logging).")
	p.Bool("version", false, "Show version information")
	p.Bool("dump-config", false, "Dump configuration to the console")

	p.Bool("enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	_ = viper.BindPFlag("leader-election.enabled", p.Lookup("enable-leader-election"))

	p.String("leader-election-name", "cluster-registry-leader-election", "Determines the name of the leader election configmap.")
	_ = viper.BindPFlag("leader-election.name", p.Lookup("leader-election-name"))
	p.String("leader-election-namespace", "", "Determines the namespace in which the leader election configmap will be created.")
	_ = viper.BindPFlag("leader-election.namespace", p.Lookup("leader-election-namespace"))

	p.Int("log-verbosity", 0, "Log verbosity")
	_ = viper.BindPFlag("log.verbosity", p.Lookup("log-verbosity"))
	p.String("log-format", "json", "Log format (console, json)")
	_ = viper.BindPFlag("log.format", p.Lookup("log-format"))

	p.String("namespace", "cluster-registry", "Namespace where the controller is running")
	_ = viper.BindPFlag("namespace", p.Lookup("namespace"))

	p.String("provision-local-cluster", "", "Name of the default local cluster to provision (if not specified no provisioning occurs)")
	_ = viper.BindPFlag("provision-local-cluster", p.Lookup("provision-local-cluster"))

	p.Bool("manage-local-cluster-secret", true, "Whether to manage secret for the local cluster")
	_ = viper.BindPFlag("manage-local-cluster-secret", p.Lookup("manage-local-cluster-secret"))

	p.String("reader-service-account-name", "cluster-registry-controller-reader", "Name of the reader service account. Used for managed cluster secret")
	_ = viper.BindPFlag("reader-service-account-name", p.Lookup("reader-service-account-name"))

	p.String("network-name", "default", "Name of the network this controller belongs to. It is used to determine api server endpoint address")
	_ = viper.BindPFlag("network-name", p.Lookup("network-name"))

	p.String("apiserver-endpoint-address", "", "Endpoint address of the API server of the cluster the controller is running on. It is used in the managed cluster secret and/or in the provisioned local cluster resource if one or both of those features are turned on.")
	_ = viper.BindPFlag("apiserver-endpoint-address", p.Lookup("apiserver-endpoint-address"))

	p.Bool("core-resources-source-enabled", true, "Whether to act as a source for core cluster api resources")
	_ = viper.BindPFlag("core-resources-source-enabled", p.Lookup("core-resources-source-enabled"))

	v.SetDefault("syncController.workerCount", 1)
	v.SetDefault("syncController.rateLimit.maxKeys", 1024)
	v.SetDefault("syncController.rateLimit.maxRatePerSecond", 5)
	v.SetDefault("syncController.rateLimit.maxBurst", 10)
	v.SetDefault("clusterController.workerCount", 2)
	v.SetDefault("clusterController.refreshIntervalSeconds", 0)

	_ = v.BindPFlags(p)
}

func bindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() { //nolint:exhaustive
		case reflect.Struct:
			bindEnvs(v.Interface(), append(parts, tv)...)
		default:
			err := viper.BindEnv(strings.Join(append(parts, tv), "."))
			if err != nil {
				panic(errors.WrapIf(err, "could not bind env variable"))
			}
		}
	}
}

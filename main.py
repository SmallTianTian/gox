#coding=utf-8
#!/usr/bin/env python3
import os
import platform
import argparse

organization = 'github.com'
projectAbsP  = None
projectName  = None
gitPath      = None

class Model():
    '''
    How to transmit data?

    Your model.getParam() should return a dict. Like this:
    {
        "<Need data model name>": {
            "model child": {
                "<Your model name, like:`logger`>": {
                    "<The param key>": "<default value>",
                    "<The param key>": "<default value>",
                }
            }
        }
    }

    For example, grpc auto set config:
    {
        "Base": {
            "config": {
                "grpc": {
                    "port": "5001"
                }
            }
        }
    }
    '''

    def render(self, data):
        ''' rener file.

        @param data Map
        '''
        pass

    def checkDependents(self, models):
        ''' check this model's dependent.

        @param  models []Model
        @throws exception, if not exist dependent
        '''
        pass

    def getParam(self):
        ''' return this model data. For dependents model.

        @return map
        '''
        return {}

    def name(self):
        ''' return name about this model.
        @return String
        '''
        return ''

class Git(Model):
    def render(self, data):
        global gitPath
        for line in os.popen('cd %s && git init && git remote add -m master origin %s && git fetch origin && git checkout origin/master && git checkout master' % (projectAbsP, gitPath)):
            print(line)

    def name(self):
        return 'git'

class Base():
    def checkDependents(self, models):
        return

    def name(self):
        return "Base"

    def render(self, data):
        self._main()
        self._logPkg()
        self._config(data)
        self._cmd()

    def getParam(self):
        result = {}
        # for log config.
        result['Base'] = {'config': {'logger':{'level': 'info', 'env': 'dev'}}}
        return result

    def _cmd(self):
        root_content = '''// Auto generate code. Thanks ['tianxuxin@126.com']
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"{organization}/{project}/config"
	"{organization}/{project}/pkg/logger"
)

var rootCmd = &cobra.Command{{
	Use:   "{project}",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {{
		cmd.Help()
	}},
}}

var baseCmd = &cobra.Command{{
	Use:   "base",
	Short: "Start base server.",
	Run: func(cmd *cobra.Command, args []string) {{
		RunBase(cmd, args)
	}},
}}

func Execute() {{
	if err := rootCmd.Execute(); err != nil {{
		fmt.Println(err)
		os.Exit(1)
	}}
}}

func init() {{
	rootCmd.AddCommand(baseCmd)
	cobra.OnInitialize(onInitialize)
}}

func onInitialize() {{
	cfg := config.InitConfig()
	logger.InitLogger(cfg.Logger.Level, cfg.Logger.Env)
}}
'''
        WriteContentToFile('root.go', root_content.format(**{'organization': organization, 'project': projectName}), 'cmd')
        user_content = '''// Auto generate code. Thanks ['tianxuxin@126.com']
package cmd

import (
    "github.com/spf13/cobra"
    "fmt"
    "os"
    "os/signal"
)

func RunBase(cmd *cobra.Command, args []string) {
    // TODO Please start your server with non-blocking

    // end
    fmt.Println("All server model is started.")
    flag := make(chan os.Signal)
    signal.Notify(flag, os.Interrupt, os.Kill)
    <-flag
    // TODO Please stop your server and clean up

    // end
    fmt.Println("\\nStop server. Bye~")
}
'''
        WriteContentToFile('base.go', user_content, 'cmd')

    def _main(self):
        template = '''// Auto generate code. Thanks ['tianxuxin@126.com']
package main
// code from auto build. 

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"{organization}/{project}/cmd"
)

func main() {{
	go func() {{
	    if err := http.ListenAndServe(":6060", nil); err != nil {{
            log.Println(err)
        }}
	}}()
	cmd.Execute()
}}
'''
        global organization, projectName
        WriteContentToFile('main.go', template.format(**{'organization': organization, 'project': projectName}))

    def _logPkg(self):
        log_template = '''// Auto generate code. Thanks ['tianxuxin@126.com']
package logger

import (
	"github.com/sirupsen/logrus"
)

var me *logrus.Entry

func InitLogger(logLevel string, env string) {
	logEntry := logrus.New()
	//logEntry.Level = logrus.WarnLevel
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.WithError(err).Error("Error parsing log level, using: info")
		level = logrus.InfoLevel
	}
	logEntry.Level = level
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "20060102T15:04:05.999"
	Formatter.FullTimestamp = true
	logEntry.SetFormatter(Formatter)
	me = logrus.NewEntry(logEntry).WithField("env", env)
}

// Pre is used to get a prepared logrus.Entry
func Pre() *logrus.Entry {
	return me
}
'''
        WriteContentToFile('logger.go', log_template, 'pkg/logger')

    def _config(self, cMap):
        global projectName
        goConfig = '''// Auto generate code. Thanks ['tianxuxin@126.com']
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var (
	DefaultConfig Config
)

type Config struct {{
{config}
	// TODO more config
}}

func InitConfig() *Config {{
	// load config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {{
		panic(err)
	}}
	// bind config value from env
	viper.SetEnvPrefix("{project}")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&DefaultConfig); err != nil {{
		fmt.Println(err)
		panic("Init config failed.")
	}}
	return &DefaultConfig
}}

func GetConfig() *Config {{
	return &DefaultConfig
}}
'''
        goData   = ''
        yamlData = ''
        if cMap['config']:
            for k in cMap['config']:
                value   = []
                maxSize = 0
                for i in cMap['config'][k]:
                    maxSize = maxSize if len(i) <= maxSize else len(i)
                    value.append(i)
                goData   += '\t{} struct {{\n'.format(k.capitalize())
                yamlData += '{}:\n'.format(k)
                for i in value:
                    goData   += '\t\t{}{} string\n'.format(i.capitalize(), (maxSize - len(i)) * ' ')
                    yamlData += '    {}: {}\n'.format(i, cMap['config'][k][i])
                goData += '\t}\n'

        WriteContentToFile('config.go', goConfig.format(**{'project': projectName, 'config': goData}), 'config')
        yamlConfig = '''# Auto generate code. Thanks ['tianxuxin@126.com']
{}
'''
        WriteContentToFile('config.yaml', yamlConfig.format(yamlData), 'config')

class Grpc():
    def render(self, data):
        self._grpcMiddleware()
        self._grpcServer()
        self._grpcFolder()

    def checkDependents(self, models):
        for item in models:
            if isinstance(item, Base):
                return True
        return False

    def getParam(self):
        result = {}
        # for log config.
        result['Base'] = {'config': {'grpc':{'port': '5001'}}}
        return result

    def name(self):
        return 'Grpc'

    def _grpcServer(self):
        grpcContent = '''// Auto generate code. Thanks ['tianxuxin@126.com']
package server

import (
	"fmt"
	"math"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"{organization}/{project}/config"
	"{organization}/{project}/pkg/logger"
	"{organization}/{project}/pkg/middleware"

	"github.com/davecgh/go-spew/spew"
	"github.com/facebookgo/stack"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pkg/errors"
)

func onPanic(p interface{{}}) error {{
	stack := stack.Callers(1)
	logger.Pre().WithField("values", spew.Sdump(p)).WithField("stack", stack).Errorln("paniced in grpc")
	return errors.WithStack(errors.New("recovered from grpc panic"))
}}

var (
	loggerOpts = []grpclogrus.Option{{
		grpclogrus.WithDurationField(func(duration time.Duration) (key string, value interface{{}}) {{
			return "grpc.duration", duration.String()
		}}),
	}}

	recoveryOpts = []grpc_recovery.Option{{
		grpc_recovery.WithRecoveryHandler(onPanic),
	}}

	grpcServer *grpc.Server
)

func GrpcStart() error {{
	port := config.DefaultConfig.Grpc.Port
	// init grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {{
		return err
	}}

	grpcServer = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
			grpclogrus.UnaryServerInterceptor(logger.Pre(), loggerOpts...),
			middleware.UnaryLogDurationServerInterceptor,
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(recoveryOpts...),
			grpclogrus.StreamServerInterceptor(logger.Pre(), loggerOpts...),
		),
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc.MaxSendMsgSize(math.MaxInt32),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{{
                        MinTime:             time.Nanosecond,
			PermitWithoutStream: true,
		}}),
	)
	// TODO set yourself protobuf server
	// Example: schemaPB.RegisterSchemaServiceServer(grpcServer, &schema.SchemaServer{{}})

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	logger.Pre().Infof("GRPC server start, listen at %d\\n", port)
	go grpcServer.Serve(lis)
	return nil
}}

func GrpcStop() {{
	grpcServer.Stop()
	logger.Pre().Infoln("GRPC server stoped.")
}}
'''
        WriteContentToFile('grpc.go', grpcContent.format(**{'organization': organization, 'project': projectName}), 'server')

    def _grpcFolder(self):
        hintContent ='''// Auto generate code. Thanks ['tianxuxin@126.com']
package grpc
// Write grpc code to this folder
'''
        WriteContentToFile('root.go', hintContent, 'grpc')

    def _grpcMiddleware(self):
        middlewareContent = '''// Auto generate code. Thanks ['tianxuxin@126.com']
package middleware

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"{organization}/{project}/pkg/logger"
)

func UnaryLogDurationServerInterceptor(ctx context.Context, req interface{{}}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{{}}, error) {{
	start := time.Now()
	resp, invokeErr := handler(ctx, req)
	duration := fmt.Sprintf("%d", int64(time.Since(start)/time.Millisecond))
	md := metadata.Pairs("duration", duration)
	err := grpc.SetHeader(ctx, md)
	if err != nil {{
		logger.Pre().Error("set duration header fail")
	}}

	logger.Pre().Info("duration: ", duration)
	return resp, invokeErr
}}
'''
        WriteContentToFile('grpc_duration.go', middlewareContent.format(**{'organization': organization, 'project': projectName}), 'pkg/middleware')

def WriteContentToFile(fileName, content, abstractPath = None):
    global projectAbsP
    if not abstractPath:
        abstractPath = projectAbsP
    else:
        abstractPath = os.path.join(projectAbsP, abstractPath)

    if not os.path.exists(abstractPath):
        os.makedirs(abstractPath)

    fn = os.path.join(abstractPath, fileName)
    with open(fn, 'w') as f:
        f.write(content)

def getGoHome():
    basePath = os.environ.get('HOME', os.path.expanduser('~')) + '/go'
    gopath   = os.environ.get('GOPATH', basePath)
    splitStr = ';' if platform.system().lower() == 'windows' else ':'
    paths    = gopath.split(splitStr)
    for p in paths:
        if p.lower() == basePath or p.lower() == '~/go' or p.lower() == '~/go/':
            return basePath
    return paths[0]

def joinMap(m1, m2):
    if isinstance(m1, dict) and isinstance(m2, dict):
        for k in m2:
            if k in m1:
                if isinstance(m1[k], dict):
                    m1[k] = joinMap(m1[k], m2[k])
            else:
                m1[k] = m2[k]
        return m1
    raise Exception("Not dict.", m1, m2)

def main():
    global projectAbsP, projectName, organization, gitPath
    parser = argparse.ArgumentParser()

    parser.add_argument('project', help='project name')
    parser.add_argument("--base", help='Create base model, like: main.go cmd...; Default use this model. If you don\'t want use base. Set: --base disable.')
    parser.add_argument('-o', '--organization', help='Reset organization. Default is `github.com`', metavar='')
    parser.add_argument("--go_path", help='Choose go path. Default use your path value.')
    parser.add_argument("--grpc", help='Create GRPC model.', action='store_true')
    parser.add_argument("--git", help='Add remote git.', metavar='')

    # 解析参数步骤
    args = parser.parse_args()
    projectName = args.project
    if args.organization:
        organization = args.organization

    goHome = getGoHome()
    if args.go_path:
        goHome = args.go_path

    projectAbsP = os.path.join(goHome, 'src', organization, projectName)
    if os.path.exists(projectAbsP):
        print("Project `{}` is exists. Absolute path: {}.".format(args.project, projectAbsP))
        _in = input("Do you want to overriding this project? [Y/N]:")
        if _in.upper() != 'Y':
            print('Give up!')
            return
        print('continue...')
    else:
        os.mkdir(projectAbsP)

    models = ([] if args.base == 'disable' else [Base()])
    if args.grpc:
        models.append(Grpc())

    if args.git:
        gitPath = args.git
        models.append(Git())

    config = {}
    for item in models:
        item.checkDependents(models)
        config = joinMap(config, item.getParam())

    for item in models:
        item.render(config.get(item.name(), {}))
    print('executed!')

if __name__ == '__main__':
    main()

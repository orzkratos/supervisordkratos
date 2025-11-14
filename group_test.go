package supervisorkratos_test

import (
	"testing"

	"github.com/orzkratos/supervisorkratos"
	"github.com/stretchr/testify/require"
)

func TestNewGroupConfig(t *testing.T) {
	// Test new GroupConfig structure with multiple programs
	// 测试新的 GroupConfig 结构与多个程序
	program1 := supervisorkratos.NewProgramConfig(
		"api-server",
		"/opt/api-server",
		"deploy",
		"/var/log/services",
	).WithStartRetries(3)

	program2 := supervisorkratos.NewProgramConfig(
		"worker",
		"/opt/worker",
		"deploy",
		"/var/log/services",
	).WithAutoStart(false)

	group := supervisorkratos.NewGroupConfig("microservices").
		AddProgram(program1).
		AddProgram(program2)

	content := supervisorkratos.GenerateGroupConfig(group)
	t.Log("=== New GroupConfig Structure ===")
	t.Log(content)

	const expected = `[group:microservices]
programs=api-server,worker


[program:api-server]
user            = deploy
directory       = /opt/api-server
command         = /opt/api-server/bin/api-server
startretries    = 3
stdout_logfile  = /var/log/services/api-server.log
stderr_logfile  = /var/log/services/api-server.err

[program:worker]
user            = deploy
directory       = /opt/worker
command         = /opt/worker/bin/worker
autostart       = false
stdout_logfile  = /var/log/services/worker.log
stderr_logfile  = /var/log/services/worker.err
`

	require.Equal(t, expected, content)
}

func TestLargeScaleGroupConfig(t *testing.T) {
	// Test large-scale group configuration
	// 测试大规模组配置
	group := supervisorkratos.NewGroupConfig("mega-cluster")

	// Create multiple programs with different configurations
	// 创建多个不同配置的程序
	for i := 1; i <= 3; i++ {
		name := "service" + string(rune('0'+i))
		program := supervisorkratos.NewProgramConfig(
			name,
			"/opt/"+name,
			"cluster-user",
			"/var/log/cluster",
		).WithPriority(50).
			WithNumProcs(2).
			WithProcessName("%(program_name)s-%(process_num)02d").
			WithEnvironment(map[string]string{
				"CLUSTER_MODE": "production",
			})

		group.AddProgram(program)
	}

	content := supervisorkratos.GenerateGroupConfig(group)
	t.Log("=== Large-scale group configuration ===")
	t.Log(content)

	const expected = `[group:mega-cluster]
programs=service1,service2,service3


[program:service1]
user            = cluster-user
directory       = /opt/service1
command         = /opt/service1/bin/service1
environment     = CLUSTER_MODE=production
stdout_logfile  = /var/log/cluster/service1.log
stderr_logfile  = /var/log/cluster/service1.err
priority        = 50
numprocs        = 2
process_name    = %(program_name)s-%(process_num)02d

[program:service2]
user            = cluster-user
directory       = /opt/service2
command         = /opt/service2/bin/service2
environment     = CLUSTER_MODE=production
stdout_logfile  = /var/log/cluster/service2.log
stderr_logfile  = /var/log/cluster/service2.err
priority        = 50
numprocs        = 2
process_name    = %(program_name)s-%(process_num)02d

[program:service3]
user            = cluster-user
directory       = /opt/service3
command         = /opt/service3/bin/service3
environment     = CLUSTER_MODE=production
stdout_logfile  = /var/log/cluster/service3.log
stderr_logfile  = /var/log/cluster/service3.err
priority        = 50
numprocs        = 2
process_name    = %(program_name)s-%(process_num)02d
`

	require.Equal(t, expected, content)
}

func TestMicroserviceGroupConfig(t *testing.T) {
	// Test microservice cluster with different service types
	// 测试微服务集群，包含不同类型的服务
	gateway := supervisorkratos.NewProgramConfig(
		"api-gateway",
		"/opt/gateway",
		"deploy",
		"/var/log/cluster",
	).WithPriority(1).
		WithNumProcs(2).
		WithProcessName("%(program_name)s-%(process_num)02d").
		WithEnvironment(map[string]string{
			"SERVICE_TYPE": "gateway",
		})

	userService := supervisorkratos.NewProgramConfig(
		"user-service",
		"/opt/user-service",
		"deploy",
		"/var/log/cluster",
	).WithStartRetries(5).
		WithStopWaitSecs(30)

	orderService := supervisorkratos.NewProgramConfig(
		"order-service",
		"/opt/order-service",
		"deploy",
		"/var/log/cluster",
	).WithAutoRestart(false).
		WithLogMaxBytes("200MB")

	cluster := supervisorkratos.NewGroupConfig("microservice-cluster").
		AddProgram(gateway).
		AddProgram(userService).
		AddProgram(orderService)

	content := supervisorkratos.GenerateGroupConfig(cluster)
	t.Log("=== Microservice cluster configuration ===")
	t.Log(content)

	const expected = `[group:microservice-cluster]
programs=api-gateway,user-service,order-service


[program:api-gateway]
user            = deploy
directory       = /opt/gateway
command         = /opt/gateway/bin/api-gateway
environment     = SERVICE_TYPE=gateway
stdout_logfile  = /var/log/cluster/api-gateway.log
stderr_logfile  = /var/log/cluster/api-gateway.err
priority        = 1
numprocs        = 2
process_name    = %(program_name)s-%(process_num)02d

[program:user-service]
user            = deploy
directory       = /opt/user-service
command         = /opt/user-service/bin/user-service
startretries    = 5
stdout_logfile  = /var/log/cluster/user-service.log
stderr_logfile  = /var/log/cluster/user-service.err
stopwaitsecs    = 30

[program:order-service]
user            = deploy
directory       = /opt/order-service
command         = /opt/order-service/bin/order-service
autorestart     = false
stdout_logfile  = /var/log/cluster/order-service.log
stdout_logfile_maxbytes = 200MB
stderr_logfile  = /var/log/cluster/order-service.err
stderr_logfile_maxbytes = 200MB
`

	require.Equal(t, expected, content)
}

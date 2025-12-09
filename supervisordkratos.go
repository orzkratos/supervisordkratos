// Package supervisordkratos: Generate supervisord configuration with Kratos microservices integration
// Provides fluent API to build supervisord program and group configs with sensible defaults
// Supports single program, multi-instance deployment, and group management scenarios
//
// supervisordkratos: 生成 supervisord 配置，集成 Kratos 微服务
// 提供流畅的 API 来构建 supervisord 程序和组配置，包含合理的默认值
// 支持单程序、多实例部署和组管理场景
package supervisordkratos

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/yyle88/must"
	"github.com/yyle88/must/mustslice"
	"github.com/yyle88/printgo"
)

// ProgramConfig single program configuration
// 单个程序配置
type ProgramConfig struct {
	// Basic program information // 基本程序信息
	Name     string // Program name // 程序名称
	UserName string // Account name to run programs // 运行程序的账户名称
	Root     string // Program root DIR // 程序根目录
	SlogRoot string // Standard output log root DIR // 标准输出日志根目录

	// Environment variables // 环境变量
	Environment *Opt[map[string]string] // Environment variables // 环境变量

	// Process settings // 进程设置
	AutoStart    *Opt[bool] // Auto start on supervisord startup // supervisord 启动时自动启动
	AutoRestart  *Opt[any]  // Auto restart on failure (bool/string: "false"/"true"/"unexpected") // 失败时自动重启（布尔值/字符串）
	StartRetries *Opt[int]  // Max start attempts // 最大启动尝试次数
	StartSecs    *Opt[int]  // Seconds to wait to confirm start success // 启动成功前等待秒数

	// Log settings // 日志设置
	LogMaxBytes    *Opt[string] // Max log file size // 最大日志文件大小
	LogBackups     *Opt[int]    // Log backup files count // 日志备份文件数量
	RedirectStderr *Opt[bool]   // Redirect stderr to stdout // 重定向 stderr 到 stdout

	// Advanced process settings // 高级进程设置
	StopAsGroup  *Opt[bool]   // Stop processes as group // 作为组停止进程
	StopWaitSecs *Opt[int]    // Stop timeout seconds // 停止超时秒数
	KillAsGroup  *Opt[bool]   // Terminate child processes as group // 作为组终止子进程
	StopSignal   *Opt[string] // Signal to stop process (TERM/INT/QUIT) // 停止进程的信号（TERM/INT/QUIT）
	Priority     *Opt[int]    // Start rank (low starts first) // 启动顺序（小值先启动）
	ExitCodes    *Opt[[]int]  // Expected exit codes // 预期退出码

	// Multi-instance settings // 多实例设置
	NumProcs    *Opt[int]    // Process instance count // 进程实例数量
	ProcessName *Opt[string] // Process name template // 进程名称模板
}

// NewProgramConfig create new ProgramConfig with required fields
// Name, Root, UserName, SlogRoot are required parameters
//
// 创建新的 ProgramConfig，需要提供必填字段
// Name、Root、UserName、SlogRoot 是必填参数
func NewProgramConfig(name string, root string, userName string, slogRoot string) *ProgramConfig {
	return &ProgramConfig{
		// Basic program information // 基本程序信息
		Name:     must.Nice(name),
		UserName: must.Nice(userName),
		Root:     must.Nice(root),
		SlogRoot: must.Nice(slogRoot),

		// Environment variables // 环境变量
		Environment: NewOpt(make(map[string]string)),

		// Set supervisord standard default values
		// 设置 supervisord 标准默认值

		// Process settings // 进程设置
		AutoStart:    NewOpt(true),
		AutoRestart:  NewOpt[any]("unexpected"), // supervisord standard default
		StartRetries: NewOpt(3),
		StartSecs:    NewOpt(1),

		// Log settings // 日志设置
		LogMaxBytes:    NewOpt("50MB"),
		LogBackups:     NewOpt(10),
		RedirectStderr: NewOpt(false),

		// Advanced process settings defaults
		// 高级进程设置默认值
		StopAsGroup:  NewOpt(false),
		StopWaitSecs: NewOpt(10),
		KillAsGroup:  NewOpt(false),
		StopSignal:   NewOpt("TERM"),
		Priority:     NewOpt(999),
		ExitCodes:    NewOpt([]int{0}),

		// Multi-instance defaults
		// 多实例默认值
		NumProcs:    NewOpt(1),
		ProcessName: NewOpt("%(program_name)s"),
	}
}

// ProgramConfig chain methods for configuration customization
// ProgramConfig 链式配置方法

// WithAutoStart set auto start flag
// 设置自动启动标志
func (p *ProgramConfig) WithAutoStart(autoStart bool) *ProgramConfig {
	p.AutoStart.Set(autoStart)
	return p
}

// WithAutoRestart set auto restart flag
// 设置自动重启标志
func (p *ProgramConfig) WithAutoRestart(autoRestart bool) *ProgramConfig {
	p.AutoRestart.Set(autoRestart)
	return p
}

// WithAutoRestartMode set auto restart mode with string value
// Accepts: "false", "true", "unexpected"
// 设置自动重启模式（字符串值）
// 接受："false"、"true"、"unexpected"
func (p *ProgramConfig) WithAutoRestartMode(mode string) *ProgramConfig {
	mustslice.In(mode, []string{"false", "true", "unexpected"})
	p.AutoRestart.Set(mode)
	return p
}

// WithStartRetries set start retries count
// 设置启动重试次数
func (p *ProgramConfig) WithStartRetries(startRetries int) *ProgramConfig {
	p.StartRetries.Set(startRetries)
	return p
}

// WithStartSecs set start seconds
// 设置启动成功等待时间
func (p *ProgramConfig) WithStartSecs(startSecs int) *ProgramConfig {
	p.StartSecs.Set(startSecs)
	return p
}

// WithLogMaxBytes set log file max bytes
// 设置日志文件最大字节数
func (p *ProgramConfig) WithLogMaxBytes(logMaxBytes string) *ProgramConfig {
	p.LogMaxBytes.Set(logMaxBytes)
	return p
}

// WithLogBackups set log backup count
// 设置日志备份文件数量
func (p *ProgramConfig) WithLogBackups(logBackups int) *ProgramConfig {
	p.LogBackups.Set(logBackups)
	return p
}

// WithRedirectStderr set stderr redirect flag
// 设置标准错误重定向标志
func (p *ProgramConfig) WithRedirectStderr(redirectStderr bool) *ProgramConfig {
	p.RedirectStderr.Set(redirectStderr)
	return p
}

// WithStopAsGroup set stop as group flag
// 设置作为组停止标志
func (p *ProgramConfig) WithStopAsGroup(stopAsGroup bool) *ProgramConfig {
	p.StopAsGroup.Set(stopAsGroup)
	return p
}

// WithKillAsGroup set terminate as group flag
// 设置作为组终止标志
func (p *ProgramConfig) WithKillAsGroup(killAsGroup bool) *ProgramConfig {
	p.KillAsGroup.Set(killAsGroup)
	return p
}

// WithStopWaitSecs set stop wait seconds
// 设置停止等待时间
func (p *ProgramConfig) WithStopWaitSecs(stopWaitSecs int) *ProgramConfig {
	p.StopWaitSecs.Set(stopWaitSecs)
	return p
}

// WithStopSignal configure the stop signal (TERM/INT/QUIT)
// 配置停止信号（TERM/INT/QUIT）
func (p *ProgramConfig) WithStopSignal(stopSignal string) *ProgramConfig {
	p.StopSignal.Set(stopSignal)
	return p
}

// WithPriority set process start rank (low starts first)
// 设置进程启动顺序（小值先启动）
func (p *ProgramConfig) WithPriority(priority int) *ProgramConfig {
	p.Priority.Set(priority)
	return p
}

// WithEnvironment set environment variables
// 设置环境变量
func (p *ProgramConfig) WithEnvironment(environment map[string]string) *ProgramConfig {
	p.Environment.Set(environment)
	return p
}

// WithExitCodes set expected exit codes
// 设置期望的退出码
func (p *ProgramConfig) WithExitCodes(exitCodes []int) *ProgramConfig {
	p.ExitCodes.Set(exitCodes)
	return p
}

// WithNumProcs set process instance count
// 设置进程实例数量
func (p *ProgramConfig) WithNumProcs(numProcs int) *ProgramConfig {
	p.NumProcs.Set(numProcs)
	return p
}

// WithProcessName set process name pattern
// 设置进程名称模式
func (p *ProgramConfig) WithProcessName(processName string) *ProgramConfig {
	p.ProcessName.Set(processName)
	return p
}

// GenerateProgramConfig generate single program configuration from ProgramConfig
// Creates supervisord INI format config with explicit values (no spacing inside)
// Includes basic info, process settings, log paths, and advanced settings
// Omits default values to keep config concise and focused on what matters
//
// GenerateProgramConfig 从 ProgramConfig 生成单个程序配置
// 创建 supervisord INI 格式配置，包含显式值（内部无空行）
// 包括基础信息、进程控制、日志路径和高级设置
// 省略默认值以保持配置简洁，专注于用户设置
func GenerateProgramConfig(program *ProgramConfig) string {
	must.Full(program)
	must.Nice(program.Name)
	must.Nice(program.Root)
	must.Nice(program.UserName)
	must.Nice(program.SlogRoot)

	ptx := printgo.NewPTX()

	// Generate program section and basic required settings
	// 生成程序段落和基本必需设置
	ptx.Println("[program:" + program.Name + "]")
	ptx.Println("user            = " + program.UserName)
	ptx.Println("directory       = " + program.Root)
	ptx.Println("command         = " + filepath.Join(program.Root, "bin", program.Name))
	// Add environment variables if set
	// 添加环境变量（如果已设置）
	if program.Environment.IsSet() {
		if env := combineSsMap(program.Environment.Get(), ","); env != "" {
			ptx.Println("environment     = " + env)
		}
	}
	// Process settings - just print explicit values
	// 进程设置 - 只打印显式设置的值
	if program.AutoStart.IsSet() {
		ptx.Println("autostart       = " + strconv.FormatBool(program.AutoStart.Get()))
	}
	if program.AutoRestart.IsSet() {
		value := program.AutoRestart.Get()
		switch v := value.(type) {
		case bool:
			ptx.Println("autorestart     = " + strconv.FormatBool(v))
		case string:
			ptx.Println("autorestart     = " + v)
		default:
			panic(errors.New("IMPOSSIBLE: INVALID TYPE"))
		}
	}
	if program.StartRetries.IsSet() {
		ptx.Println("startretries    = " + strconv.Itoa(program.StartRetries.Get()))
	}
	if program.StartSecs.IsSet() {
		ptx.Println("startsecs       = " + strconv.Itoa(program.StartSecs.Get()))
	}
	// Log settings always show (required for paths)
	// 日志设置始终显示（路径必需）
	ptx.Println("stdout_logfile  = " + filepath.Join(program.SlogRoot, program.Name+".log"))
	if program.LogMaxBytes.IsSet() {
		ptx.Println("stdout_logfile_maxbytes = " + program.LogMaxBytes.Get())
	}
	if program.LogBackups.IsSet() {
		ptx.Println("stdout_logfile_backups = " + strconv.Itoa(program.LogBackups.Get()))
	}
	ptx.Println("stderr_logfile  = " + filepath.Join(program.SlogRoot, program.Name+".err"))
	if program.LogMaxBytes.IsSet() {
		ptx.Println("stderr_logfile_maxbytes = " + program.LogMaxBytes.Get())
	}
	if program.LogBackups.IsSet() {
		ptx.Println("stderr_logfile_backups = " + strconv.Itoa(program.LogBackups.Get()))
	}
	if program.RedirectStderr.IsSet() {
		ptx.Println("redirect_stderr = " + strconv.FormatBool(program.RedirectStderr.Get()))
	}
	// Advanced process settings - just non-defaults
	// 高级进程设置 - 只显示非默认值
	if program.StopAsGroup.IsSet() {
		ptx.Println("stopasgroup     = " + strconv.FormatBool(program.StopAsGroup.Get()))
	}
	if program.StopWaitSecs.IsSet() {
		ptx.Println("stopwaitsecs    = " + strconv.Itoa(program.StopWaitSecs.Get()))
	}
	if program.KillAsGroup.IsSet() {
		ptx.Println("killasgroup     = " + strconv.FormatBool(program.KillAsGroup.Get()))
	}
	if program.StopSignal.IsSet() {
		ptx.Println("stopsignal      = " + program.StopSignal.Get())
	}
	if program.Priority.IsSet() {
		ptx.Println("priority        = " + strconv.Itoa(program.Priority.Get()))
	}
	if program.ExitCodes.IsSet() {
		ptx.Println("exitcodes       = " + combineInts(program.ExitCodes.Get(), ","))
	}
	if program.NumProcs.IsSet() {
		ptx.Println("numprocs        = " + strconv.Itoa(program.NumProcs.Get()))
	}
	if program.ProcessName.IsSet() {
		ptx.Println("process_name    = " + program.ProcessName.Get())
	}

	return ptx.String()
}

// combineInts converts int slice to comma-separated string
// Returns blank string if input is blank
//
// combineInts 将整数切片转换为逗号分隔的字符串
// 输入为空时返回空字符串
func combineInts(items []int, sep string) string {
	if len(items) == 0 {
		return ""
	}
	results := make([]string, len(items))
	for i, item := range items {
		results[i] = strconv.Itoa(item)
	}
	return strings.Join(results, sep)
}

// combineSsMap converts string map to name=value pairs joined with sep
// Used to format environment variables as KEY1=VALUE1,KEY2=VALUE2
// Returns blank string if input is blank
//
// combineSsMap 将字符串映射转换为由分隔符连接的键值对
// 用于格式化环境变量为 KEY1=VALUE1,KEY2=VALUE2 格式
// 输入为空时返回空字符串
func combineSsMap(items map[string]string, sep string) string {
	if len(items) == 0 {
		return ""
	}
	pairs := make([]string, 0, len(items))
	for key, value := range items {
		pairs = append(pairs, key+"="+value)
	}
	return strings.Join(pairs, sep)
}

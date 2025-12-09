package supervisordkratos

import (
	"strings"

	"github.com/yyle88/must"
	"github.com/yyle88/printgo"
)

// GroupConfig supervisord group configuration
// supervisord 组配置
type GroupConfig struct {
	Name     string           // Group name // 组名称
	Programs []*ProgramConfig // Program configs // 程序配置列表
}

// NewGroupConfig create new GroupConfig
// 创建新的 GroupConfig
func NewGroupConfig(name string) *GroupConfig {
	return &GroupConfig{
		Name:     must.Nice(name),
		Programs: make([]*ProgramConfig, 0),
	}
}

// AddProgram add program to group
// 添加程序到组
func (g *GroupConfig) AddProgram(program *ProgramConfig) *GroupConfig {
	g.Programs = append(g.Programs, program)
	return g
}

// GenerateGroupConfig generate supervisord group configuration in INI format
// Creates complete group config with name section and programs
// Outputs group section then program sections with spacing
//
// GenerateGroupConfig 生成 INI 格式的 supervisord 组配置
// 创建包含名称段和程序的完整组配置
// 输出组段落然后输出程序段落，使用间距
func GenerateGroupConfig(group *GroupConfig) string {
	must.Full(group)
	must.Nice(group.Name)
	must.Have(group.Programs)

	ptx := printgo.NewPTX()

	// Generate group name section
	// 生成组名称段
	ptx.Println(`[group:` + group.Name + `]`)
	programs := make([]string, 0, len(group.Programs))
	for _, p := range group.Programs {
		programs = append(programs, p.Name)
	}
	ptx.Println(`programs=` + strings.Join(programs, ","))
	ptx.Println()

	// Generate each program config
	// 生成每个程序配置
	for _, program := range group.Programs {
		ptx.Println()
		cfs := GenerateProgramConfig(program)
		ptx.Println(strings.TrimSpace(cfs))
	}

	return ptx.String()
}

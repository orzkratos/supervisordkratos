package supervisordkratos

// Opt represents a config value that can track if it has been set
// Helps distinguish between defaults and custom-specified values
// Generic type T allows flexible usage across config fields
//
// Opt 表示可跟踪是否已设置的可选配置值
// 帮助区分默认值和自定义指定的值
// 泛型类型 T 允许在配置字段中灵活使用
type Opt[T any] struct {
	Value T    // Stored value // 存储的值
	isSet bool // Track if value was set // 跟踪值是否已被设置
}

// NewOpt creates new Opt with default value (not marked as set)
// The value is stored but isSet flag stays false
// Use Set() to mark as custom-configured value
//
// NewOpt 创建新的 Opt 并设置默认值（未标记为已设置）
// 值已存储但 isSet 标志保持 false
// 使用 Set() 标记为自定义配置的值
func NewOpt[T any](v T) *Opt[T] {
	return &Opt[T]{Value: v, isSet: false}
}

// Get returns current stored value
// Returns value regardless of isSet flag status
//
// Get 返回当前存储的值
// 不管 isSet 标志状态都返回值
func (sv *Opt[T]) Get() T {
	return sv.Value
}

// Set updates value and marks as custom-set
// Indicates this value has been configured
//
// Set 更新值并标记为自定义设置
// 表示此值已被配置
func (sv *Opt[T]) Set(v T) {
	sv.Value = v
	sv.isSet = true
}

// IsSet checks if value has been set via Set()
// Returns false when using defaults from NewOpt()
// Returns true when Set() has been invoked
//
// IsSet 检查值是否已通过 Set() 设置
// 使用 NewOpt() 默认值时返回 false
// 调用 Set() 后返回 true
func (sv *Opt[T]) IsSet() bool {
	return sv.isSet
}

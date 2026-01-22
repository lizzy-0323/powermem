package intelligence

import (
	"math"
	"time"
)

// EbbinghausManager Ebbinghaus 遗忘曲线管理器
type EbbinghausManager struct {
	decayRate           float64 // 衰减率
	reinforcementFactor float64 // 强化因子
}

// NewEbbinghausManager 创建 Ebbinghaus 管理器
func NewEbbinghausManager(decayRate, reinforcementFactor float64) *EbbinghausManager {
	return &EbbinghausManager{
		decayRate:           decayRate,
		reinforcementFactor: reinforcementFactor,
	}
}

// CalculateRetention 计算记忆保持强度
// retention = e^(-decay_rate * time_elapsed)
func (m *EbbinghausManager) CalculateRetention(createdAt time.Time, lastAccessedAt *time.Time) float64 {
	now := time.Now()
	var timeElapsed time.Duration

	if lastAccessedAt != nil {
		timeElapsed = now.Sub(*lastAccessedAt)
	} else {
		timeElapsed = now.Sub(createdAt)
	}

	// 将时间转换为小时
	hoursElapsed := timeElapsed.Hours()

	// 应用 Ebbinghaus 公式
	retention := math.Exp(-m.decayRate * hoursElapsed / 24.0)

	return retention
}

// Reinforce 强化记忆
// 每次访问时，增加记忆强度
func (m *EbbinghausManager) Reinforce(currentStrength float64) float64 {
	// 强化公式：new_strength = min(1.0, current_strength + reinforcement_factor * (1 - current_strength))
	newStrength := currentStrength + m.reinforcementFactor*(1.0-currentStrength)
	if newStrength > 1.0 {
		return 1.0
	}
	return newStrength
}

// ShouldArchive 判断记忆是否应该归档
// 如果保持强度低于阈值，则应该归档
func (m *EbbinghausManager) ShouldArchive(retentionStrength float64, threshold float64) bool {
	if threshold == 0 {
		threshold = 0.2 // 默认阈值
	}
	return retentionStrength < threshold
}

// CalculateNextReview 计算下次复习时间
// 基于记忆强度，强度越高，复习间隔越长
func (m *EbbinghausManager) CalculateNextReview(retentionStrength float64) time.Time {
	// 复习间隔（小时） = 24 * (1 + strength * 10)
	// strength 越高，间隔越长
	hoursUntilReview := 24.0 * (1.0 + retentionStrength*10.0)
	return time.Now().Add(time.Duration(hoursUntilReview) * time.Hour)
}

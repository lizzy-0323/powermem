package intelligence_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/oceanbase/powermem-go/pkg/intelligence"
)

func TestEbbinghausManager(t *testing.T) {
	decayRate := 0.1
	reinforcementFactor := 0.3
	
	manager := intelligence.NewEbbinghausManager(decayRate, reinforcementFactor)
	assert.NotNil(t, manager)
}

func TestCalculateRetention(t *testing.T) {
	decayRate := 0.1
	reinforcementFactor := 0.3
	
	manager := intelligence.NewEbbinghausManager(decayRate, reinforcementFactor)
	
	// 测试初始强度（刚创建）
	createdAt := time.Now()
	retention := manager.CalculateRetention(createdAt, nil)
	assert.Greater(t, retention, 0.0, "保持强度应该大于 0")
	assert.LessOrEqual(t, retention, 1.0, "保持强度应该不超过 1.0")
	
	// 测试时间衰减（1 天后）
	createdAt = time.Now().Add(-24 * time.Hour)
	retention = manager.CalculateRetention(createdAt, nil)
	assert.Less(t, retention, 1.0, "时间衰减应该降低强度")
	assert.Greater(t, retention, 0.0, "强度应该大于 0")
	
	// 测试访问强化
	currentStrength := 0.5
	reinforced := manager.Reinforce(currentStrength)
	assert.Greater(t, reinforced, currentStrength, "强化应该增加强度")
	assert.LessOrEqual(t, reinforced, 1.0, "强度应该不超过 1.0")
}

func TestEbbinghausDecay(t *testing.T) {
	decayRate := 0.1
	reinforcementFactor := 0.3
	
	manager := intelligence.NewEbbinghausManager(decayRate, reinforcementFactor)
	
	// 测试不同时间点的衰减
	now := time.Now()
	testCases := []struct {
		hoursAgo  float64
		wantLower bool
	}{
		{0, false},
		{1, true},
		{24, true},
		{168, true}, // 1 周
	}
	
	for _, tc := range testCases {
		createdAt := now.Add(-time.Duration(tc.hoursAgo) * time.Hour)
		retention := manager.CalculateRetention(createdAt, nil)
		if tc.wantLower {
			assert.Less(t, retention, 1.0, 
				"时间 %v 小时后强度应该降低", tc.hoursAgo)
		}
		assert.Greater(t, retention, 0.0, "强度应该始终大于 0")
		assert.LessOrEqual(t, retention, 1.0, "强度应该不超过 1.0")
	}
}

func TestReinforcementFactor(t *testing.T) {
	decayRate := 0.1
	reinforcementFactor := 0.3
	
	manager := intelligence.NewEbbinghausManager(decayRate, reinforcementFactor)
	
	// 测试强化功能
	currentStrength := 0.5
	reinforced := manager.Reinforce(currentStrength)
	
	assert.Greater(t, reinforced, currentStrength, 
		"强化应该增加记忆强度")
	assert.LessOrEqual(t, reinforced, 1.0, "强度应该不超过 1.0")
}

func TestEbbinghausEdgeCases(t *testing.T) {
	decayRate := 0.1
	reinforcementFactor := 0.3
	
	manager := intelligence.NewEbbinghausManager(decayRate, reinforcementFactor)
	
	// 测试极端情况
	now := time.Now()
	
	// 很长时间前创建
	oldCreatedAt := now.Add(-1000 * time.Hour)
	retention := manager.CalculateRetention(oldCreatedAt, nil)
	assert.Greater(t, retention, 0.0)
	assert.Less(t, retention, 1.0)
	
	// 测试强化上限
	highStrength := 0.99
	reinforced := manager.Reinforce(highStrength)
	assert.LessOrEqual(t, reinforced, 1.0, "强化后不应该超过 1.0")
}

package postgres

import (
	"fmt"
	"strings"
)

// buildWhereClause 构建 WHERE 子句（修复 DeleteAll）
func buildWhereClause(userID, agentID string, filters map[string]interface{}) (string, []interface{}) {
	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if userID != "" {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, userID)
		argIndex++
	}

	if agentID != "" {
		conditions = append(conditions, fmt.Sprintf("agent_id = $%d", argIndex))
		args = append(args, agentID)
		argIndex++
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}

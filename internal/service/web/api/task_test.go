package api

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"encoding/json"
	"testing"
)

func TestMaskSSHParams(t *testing.T) {
	params := `{"host":"192.168.1.100","port":22,"username":"root","password":"secret123","mode":"script"}`
	masked := maskSSHParams(global.TaskProtocolSSH, params)

	var config map[string]interface{}
	if err := json.Unmarshal([]byte(masked), &config); err != nil {
		t.Fatalf("掩码后的参数不是合法JSON: %v", err)
	}
	if config["password"] != sshPasswordMask {
		t.Fatalf("密码未被掩码, 期望 %q, 实际 %v", sshPasswordMask, config["password"])
	}
	if config["host"] != "192.168.1.100" {
		t.Fatalf("非敏感字段被误改: %v", config["host"])
	}
}

func TestMaskSSHParamsSkipsNonSSH(t *testing.T) {
	params := `{"url":"http://example.com","method":"GET"}`
	if got := maskSSHParams(global.TaskProtocolHttp, params); got != params {
		t.Fatalf("非SSH任务的参数不应被修改, 期望 %q, 实际 %q", params, got)
	}
}

func TestMaskSSHParamsEmpty(t *testing.T) {
	if got := maskSSHParams(global.TaskProtocolSSH, ""); got != "" {
		t.Fatalf("空参数应原样返回, 实际 %q", got)
	}
}

func TestSSHPasswordNeedsPreserve(t *testing.T) {
	cases := []struct {
		name   string
		params string
		want   bool
	}{
		{"空参数需要保留", "", true},
		{"掩码占位符需要保留", `{"host":"h","password":"********"}`, true},
		{"空密码需要保留", `{"host":"h","password":""}`, true},
		{"缺少密码字段需要保留", `{"host":"h"}`, true},
		{"真实密码不需要保留", `{"host":"h","password":"real-pass"}`, false},
		{"非法JSON不保留", `not-json`, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := sshPasswordNeedsPreserve(tc.params); got != tc.want {
				t.Fatalf("sshPasswordNeedsPreserve(%q) = %v, 期望 %v", tc.params, got, tc.want)
			}
		})
	}
}

func TestTaskDependencyWouldCycle(t *testing.T) {
	dep := func(taskID, dependentID uint) models.TaskDependency {
		return models.TaskDependency{TaskID: taskID, DependentID: dependentID}
	}

	cases := []struct {
		name        string
		taskID      uint
		dependentID uint
		deps        []models.TaskDependency
		want        bool
	}{
		{"依赖自己", 1, 1, nil, true},
		{"直接循环: 1依赖2, 2已依赖1", 1, 2, []models.TaskDependency{dep(2, 1)}, true},
		{"多级循环: A->B->C->A", 1, 3, []models.TaskDependency{dep(3, 2), dep(2, 1)}, true},
		{"无环: 3依赖2, 2依赖1", 3, 2, []models.TaskDependency{dep(2, 1)}, false},
		{"空依赖无环", 1, 2, nil, false},
		{"无关链路无环", 5, 4, []models.TaskDependency{dep(2, 1), dep(3, 2)}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := taskDependencyWouldCycle(tc.taskID, tc.dependentID, tc.deps); got != tc.want {
				t.Fatalf("taskDependencyWouldCycle(%d, %d) = %v, 期望 %v", tc.taskID, tc.dependentID, got, tc.want)
			}
		})
	}
}

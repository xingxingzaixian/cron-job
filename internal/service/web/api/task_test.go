package api

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/viper"
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

func TestEncryptSSHParams(t *testing.T) {
	old := viper.GetString("security.secret")
	viper.Set("security.secret", "ssh-test-key")
	defer viper.Set("security.secret", old)

	params := `{"host":"h","port":22,"username":"u","password":"real-pass","mode":"script"}`
	enc := encryptSSHParams(global.TaskProtocolSSH, params)
	if enc == params {
		t.Fatal("SSH密码应被加密")
	}

	var config map[string]interface{}
	if err := json.Unmarshal([]byte(enc), &config); err != nil {
		t.Fatalf("加密后参数不是合法JSON: %v", err)
	}
	if pw, _ := config["password"].(string); !strings.HasPrefix(pw, "enc:v1:") {
		t.Fatalf("密码应为密文: %v", config["password"])
	}

	// 已加密的值不再重复加密
	enc2 := encryptSSHParams(global.TaskProtocolSSH, enc)
	if enc2 != enc {
		t.Fatal("已加密的密码不应重复加密")
	}

	// 非SSH任务原样返回
	httpParams := `{"url":"http://x","method":"GET"}`
	if got := encryptSSHParams(global.TaskProtocolHttp, httpParams); got != httpParams {
		t.Fatal("非SSH任务参数不应被修改")
	}
}

func TestNewTaskHandler(t *testing.T) {
	tests := []struct {
		name     string
		protocol global.TaskProtocol
		wantErr  bool
	}{
		{"http", global.TaskProtocolHttp, false},
		{"shell", global.TaskProtocolShell, false},
		{"ssh", global.TaskProtocolSSH, false},
		{"invalid", global.TaskProtocol(99), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := newTaskHandler(tt.protocol)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for protocol %d, got handler %T", tt.protocol, h)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for protocol %d: %v", tt.protocol, err)
			}
			if h == nil {
				t.Fatalf("expected non-nil handler for protocol %d", tt.protocol)
			}
		})
	}
}

package handler

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"encoding/json"
	"testing"
)

// TestSSHConfigParsing 测试SSH配置解析
func TestSSHConfigParsing(t *testing.T) {
	// 测试用例1：密码认证
	config1 := SSHConfig{
		Host:     "192.168.1.100",
		Port:     22,
		Username: "root",
		Password: "password123",
		Mode:     "sequential",
	}

	jsonData1, err := json.Marshal(config1)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	task1 := &models.Task{
		Params:   string(jsonData1), // SSH配置存储在Params字段
		Command:  "whoami\npwd",     // 命令存储在Command字段
		Protocol: global.TaskProtocolSSH,
		Timeout:  300,
	}

	handler := &SSHHandler{}
	parsedConfig1, err := handler.parseSSHConfig(task1)
	if err != nil {
		t.Fatalf("SSH配置解析失败: %v", err)
	}

	if parsedConfig1.Host != config1.Host {
		t.Errorf("主机地址解析错误，期望: %s, 实际: %s", config1.Host, parsedConfig1.Host)
	}

	// 测试命令解析
	commands1, err := handler.parseCommands(task1.Command)
	if err != nil {
		t.Fatalf("命令解析失败: %v", err)
	}

	if len(commands1) != 2 {
		t.Errorf("命令数量解析错误，期望: 2, 实际: %d", len(commands1))
	}

	// 测试用例2：私钥认证
	config2 := SSHConfig{
		Host:       "192.168.1.101",
		Port:       2222,
		Username:   "ubuntu",
		Mode:       "script",
	}

	jsonData2, err := json.Marshal(config2)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	task2 := &models.Task{
		Params:   string(jsonData2), // SSH配置存储在Params字段
		Command:  "ls -la\ndf -h",   // 命令存储在Command字段
		Protocol: global.TaskProtocolSSH,
		Timeout:  600,
	}

	parsedConfig2, err := handler.parseSSHConfig(task2)
	if err != nil {
		t.Fatalf("SSH配置解析失败: %v", err)
	}

	if parsedConfig2.Port != config2.Port {
		t.Errorf("端口解析错误，期望: %d, 实际: %d", config2.Port, parsedConfig2.Port)
	}

	// 测试命令解析
	commands2, err := handler.parseCommands(task2.Command)
	if err != nil {
		t.Fatalf("命令解析失败: %v", err)
	}

	if len(commands2) != 2 {
		t.Errorf("命令数量解析错误，期望: 2, 实际: %d", len(commands2))
	}
}

// TestSSHConfigValidation 测试SSH配置验证
func TestSSHConfigValidation(t *testing.T) {
	handler := &SSHHandler{}

	// 测试用例1：有效配置
	validConfig := &SSHConfig{
		Host:     "192.168.1.100",
		Port:     22,
		Username: "root",
		Password: "password123",
		Mode:     "sequential",
	}
	validCommands := []string{"whoami"}

	if err := handler.validateSSHConfig(validConfig, validCommands); err != nil {
		t.Errorf("有效配置验证失败: %v", err)
	}

	// 测试用例2：缺少主机地址
	invalidConfig1 := &SSHConfig{
		Port:     22,
		Username: "root",
		Password: "password123",
	}
	testCommands := []string{"whoami"}

	if err := handler.validateSSHConfig(invalidConfig1, testCommands); err == nil {
		t.Error("应该检测到缺少主机地址的错误")
	}

	// 测试用例3：缺少认证信息
	invalidConfig2 := &SSHConfig{
		Host:     "192.168.1.100",
		Port:     22,
		Username: "root",
	}

	if err := handler.validateSSHConfig(invalidConfig2, testCommands); err == nil {
		t.Error("应该检测到缺少认证信息的错误")
	}

	// 测试用例4：端口号无效
	invalidConfig3 := &SSHConfig{
		Host:     "192.168.1.100",
		Port:     70000,
		Username: "root",
		Password: "password123",
	}

	if err := handler.validateSSHConfig(invalidConfig3, testCommands); err == nil {
		t.Error("应该检测到无效端口号的错误")
	}

	// 测试用例5：命令列表为空
	invalidConfig4 := &SSHConfig{
		Host:     "192.168.1.100",
		Port:     22,
		Username: "root",
		Password: "password123",
	}
	emptyCommands := []string{}

	if err := handler.validateSSHConfig(invalidConfig4, emptyCommands); err == nil {
		t.Error("应该检测到命令列表为空的错误")
	}
}

// TestGetSSHConfigExample 测试获取SSH配置示例
func TestGetSSHConfigExample(t *testing.T) {
	configExample, commandExample := GetSSHConfigExample()
	if configExample == "" {
		t.Error("SSH配置示例不应该为空")
	}
	if commandExample == "" {
		t.Error("SSH命令示例不应该为空")
	}

	// 验证示例是否为有效的JSON
	var config SSHConfig
	if err := json.Unmarshal([]byte(configExample), &config); err != nil {
		t.Errorf("SSH配置示例不是有效的JSON: %v", err)
	}

	// 解析命令示例
	handler := &SSHHandler{}
	commands, err := handler.parseCommands(commandExample)
	if err != nil {
		t.Errorf("命令示例解析失败: %v", err)
	}

	// 验证示例配置是否有效
	if err := handler.validateSSHConfig(&config, commands); err != nil {
		t.Errorf("SSH配置示例验证失败: %v", err)
	}
}

// BenchmarkSSHConfigParsing 性能测试：SSH配置解析
func BenchmarkSSHConfigParsing(b *testing.B) {
	config := SSHConfig{
		Host:     "192.168.1.100",
		Port:     22,
		Username: "root",
		Password: "password123",
		Mode:     "sequential",
	}

	jsonData, _ := json.Marshal(config)
	task := &models.Task{
		Params:   string(jsonData),      // SSH配置存储在Params字段
		Command:  "whoami\npwd\nls -la", // 命令存储在Command字段
		Protocol: global.TaskProtocolSSH,
		Timeout:  300,
	}

	handler := &SSHHandler{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handler.parseSSHConfig(task)
		if err != nil {
			b.Fatalf("SSH配置解析失败: %v", err)
		}
	}
}

// ExampleSSHHandler_Run 示例：如何使用SSH Handler
func ExampleSSHHandler_Run() {
	// 创建SSH配置
	config := SSHConfig{
		Host:     "192.168.1.100",
		Port:     22,
		Username: "root",
		Password: "your_password",
		Mode:     "sequential",
	}

	// 序列化为JSON
	jsonData, _ := json.Marshal(config)

	// 创建任务模型
	task := &models.Task{
		Params:   string(jsonData),           // SSH配置存储在Params字段
		Command:  "whoami\npwd\nls -la /tmp", // 命令存储在Command字段
		Protocol: global.TaskProtocolSSH,
		Timeout:  300,
	}

	// 创建SSH Handler并执行
	handler := &SSHHandler{}
	result, err := handler.Run(task, 12345)
	if err != nil {
		// 处理错误
		return
	}

	// 处理结果
	_ = result
}

// 提供一些常用的SSH配置模板
func GetSSHConfigTemplates() map[string]map[string]string {
	templates := make(map[string]map[string]string)

	// 密码认证模板
	passwordAuth := SSHConfig{
		Host:     "your.server.com",
		Port:     22,
		Username: "root",
		Password: "your_password",
		Mode:     "sequential",
	}
	passwordData, _ := json.MarshalIndent(passwordAuth, "", "  ")
	passwordCommands := `whoami
uptime
df -h`
	templates["password_auth"] = map[string]string{
		"params":  string(passwordData),
		"command": passwordCommands,
	}


	// 脚本执行模板
	scriptMode := SSHConfig{
		Host:     "your.server.com",
		Port:     22,
		Username: "root",
		Password: "your_password",
		Mode:     "script",
	}
	scriptData, _ := json.MarshalIndent(scriptMode, "", "  ")
	scriptCommands := `cd /tmp
echo 'Hello World' > test.txt
cat test.txt
rm test.txt`
	templates["script_mode"] = map[string]string{
		"params":  string(scriptData),
		"command": scriptCommands,
	}

	return templates
}

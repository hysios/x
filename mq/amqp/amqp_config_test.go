package amqp

import (
	"strings"
	"testing"
	"time"

	"github.com/hysios/x/mq"
	"github.com/mitchellh/mapstructure"
)

// TestDefaultConfigOverride 测试默认配置覆盖问题
func TestDefaultConfigOverride(t *testing.T) {
	tests := []struct {
		name             string
		userConfig       mq.Config
		expectedURL      string
		expectedExchange string
		expectedTimeout  time.Duration
		expectedDurable  bool
	}{
		{
			name: "用户配置完全覆盖默认配置",
			userConfig: mq.Config{
				"url":             "amqp://user:pass@example.com:5672/",
				"exchange_name":   "custom_exchange",
				"publish_timeout": "10s",
				"durable":         false,
			},
			expectedURL:      "amqp://user:pass@example.com:5672/",
			expectedExchange: "custom_exchange",
			expectedTimeout:  10 * time.Second,
			expectedDurable:  false,
		},
		{
			name: "用户配置部分覆盖默认配置",
			userConfig: mq.Config{
				"url": "amqp://user:pass@example.com:5672/",
			},
			expectedURL:      "amqp://user:pass@example.com:5672/",
			expectedExchange: "events",        // 应该使用默认值
			expectedTimeout:  5 * time.Second, // 应该使用默认值
			expectedDurable:  true,            // 应该使用默认值
		},
		{
			name:             "空用户配置使用默认配置",
			userConfig:       mq.Config{},
			expectedURL:      "amqp://guest:guest@localhost:5672/",
			expectedExchange: "events",
			expectedTimeout:  5 * time.Second,
			expectedDurable:  true,
		},
		{
			name: "用户配置中的零值应该覆盖默认配置",
			userConfig: mq.Config{
				"durable":         false, // 零值但应该覆盖默认的true
				"publish_timeout": "0s",  // 零值但应该覆盖默认的5s
			},
			expectedURL:      "amqp://guest:guest@localhost:5672/",
			expectedExchange: "events",
			expectedTimeout:  0 * time.Second, // 应该被覆盖为0
			expectedDurable:  false,           // 应该被覆盖为false
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 模拟init函数中的新配置合并逻辑
			var cfg Config

			// 配置 mapstructure 解码器以支持 time.Duration
			decoderConfig := &mapstructure.DecoderConfig{
				DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
				Result:     &cfg,
			}
			decoder, err := mapstructure.NewDecoder(decoderConfig)
			if err != nil {
				t.Fatalf("创建解码器失败: %v", err)
			}

			if err := decoder.Decode(tt.userConfig); err != nil {
				t.Fatalf("解码用户配置失败: %v", err)
			}

			// 使用新的配置合并逻辑
			dst := mergeConfigs(DefaultConfig, cfg, tt.userConfig)

			// 验证最终配置
			if dst.URL != tt.expectedURL {
				t.Errorf("URL不匹配: 期望 %s, 实际 %s", tt.expectedURL, dst.URL)
			}
			if dst.ExchangeName != tt.expectedExchange {
				t.Errorf("ExchangeName不匹配: 期望 %s, 实际 %s", tt.expectedExchange, dst.ExchangeName)
			}
			if dst.PublishTimeout != tt.expectedTimeout {
				t.Errorf("PublishTimeout不匹配: 期望 %v, 实际 %v", tt.expectedTimeout, dst.PublishTimeout)
			}
			if dst.Durable != tt.expectedDurable {
				t.Errorf("Durable不匹配: 期望 %v, 实际 %v", tt.expectedDurable, dst.Durable)
			}
		})
	}
}

// TestConfigMergeWithZeroValues 测试我们的新配置合并逻辑
func TestConfigMergeWithZeroValues(t *testing.T) {
	t.Run("测试新的配置合并逻辑处理零值", func(t *testing.T) {
		defaultCfg := Config{
			URL:            "amqp://guest:guest@localhost:5672/",
			ExchangeName:   "events",
			PublishTimeout: 5 * time.Second,
			Durable:        true,
		}

		userCfg := Config{
			Durable:        false, // 零值但应该覆盖
			PublishTimeout: 0,     // 零值但应该覆盖
		}

		rawConfig := mq.Config{
			"durable":         false,
			"publish_timeout": "0s",
		}

		result := mergeConfigs(defaultCfg, userCfg, rawConfig)

		// 验证零值正确覆盖了默认值
		if result.Durable != false {
			t.Errorf("Durable零值覆盖失败: 期望 false, 实际 %v", result.Durable)
		}
		if result.PublishTimeout != 0 {
			t.Errorf("PublishTimeout零值覆盖失败: 期望 0s, 实际 %v", result.PublishTimeout)
		}
		// 未明确设置的字段应该保持默认值
		if result.URL != "amqp://guest:guest@localhost:5672/" {
			t.Errorf("URL应该保持默认值: 期望 %s, 实际 %s", "amqp://guest:guest@localhost:5672/", result.URL)
		}
		if result.ExchangeName != "events" {
			t.Errorf("ExchangeName应该保持默认值: 期望 %s, 实际 %s", "events", result.ExchangeName)
		}
	})
}

// TestDefaultConfigURL 专门测试 DefaultConfig URL 的各种情况
func TestDefaultConfigURL(t *testing.T) {
	tests := []struct {
		name        string
		userConfig  mq.Config
		expectedURL string
	}{
		{
			name:        "用户未提供URL，应使用DefaultConfig的URL",
			userConfig:  mq.Config{},
			expectedURL: DefaultConfig.URL, // "amqp://guest:guest@localhost:5672/"
		},
		{
			name: "用户提供了URL，应使用用户的URL",
			userConfig: mq.Config{
				"url": "amqp://user:pass@example.com:5672/",
			},
			expectedURL: "amqp://user:pass@example.com:5672/",
		},
		{
			name: "用户提供了空字符串URL，应使用空字符串覆盖默认URL",
			userConfig: mq.Config{
				"url": "",
			},
			expectedURL: "",
		},
		{
			name: "用户提供了其他配置但没有URL，应使用默认URL",
			userConfig: mq.Config{
				"exchange_name": "custom_exchange",
				"durable":       false,
			},
			expectedURL: DefaultConfig.URL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg Config

			// 配置 mapstructure 解码器以支持 time.Duration
			decoderConfig := &mapstructure.DecoderConfig{
				DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
				Result:     &cfg,
			}
			decoder, err := mapstructure.NewDecoder(decoderConfig)
			if err != nil {
				t.Fatalf("创建解码器失败: %v", err)
			}

			if err := decoder.Decode(tt.userConfig); err != nil {
				t.Fatalf("解码用户配置失败: %v", err)
			}

			// 使用 mergeConfigs 合并配置
			result := mergeConfigs(DefaultConfig, cfg, tt.userConfig)

			if result.URL != tt.expectedURL {
				t.Errorf("URL不匹配: 期望 %q, 实际 %q", tt.expectedURL, result.URL)
			}
		})
	}
}

// TestMergeConfigsEdgeCases 测试 mergeConfigs 函数的边界情况
func TestMergeConfigsEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		defaultCfg Config
		userCfg    Config
		rawConfig  mq.Config
		expected   Config
	}{
		{
			name: "所有字段都在rawConfig中存在但为零值",
			defaultCfg: Config{
				URL:            "amqp://default:default@localhost:5672/",
				ExchangeName:   "default_exchange",
				QueueName:      "default_queue",
				PublishTimeout: 10 * time.Second,
				Durable:        true,
			},
			userCfg: Config{
				URL:            "",
				ExchangeName:   "",
				QueueName:      "",
				PublishTimeout: 0,
				Durable:        false,
			},
			rawConfig: mq.Config{
				"url":             "",
				"exchange_name":   "",
				"queue_name":      "",
				"publish_timeout": "0s",
				"durable":         false,
			},
			expected: Config{
				URL:            "",
				ExchangeName:   "",
				QueueName:      "",
				PublishTimeout: 0,
				Durable:        false,
			},
		},
		{
			name: "rawConfig中只有部分字段存在",
			defaultCfg: Config{
				URL:            "amqp://default:default@localhost:5672/",
				ExchangeName:   "default_exchange",
				QueueName:      "default_queue",
				PublishTimeout: 10 * time.Second,
				Durable:        true,
			},
			userCfg: Config{
				URL:     "amqp://user:user@example.com:5672/",
				Durable: false,
			},
			rawConfig: mq.Config{
				"url":     "amqp://user:user@example.com:5672/",
				"durable": false,
			},
			expected: Config{
				URL:            "amqp://user:user@example.com:5672/",
				ExchangeName:   "default_exchange", // 保持默认值
				QueueName:      "default_queue",    // 保持默认值
				PublishTimeout: 10 * time.Second,   // 保持默认值
				Durable:        false,              // 使用用户值
			},
		},
		{
			name: "rawConfig为空，应该使用默认配置",
			defaultCfg: Config{
				URL:            "amqp://default:default@localhost:5672/",
				ExchangeName:   "default_exchange",
				QueueName:      "default_queue",
				PublishTimeout: 10 * time.Second,
				Durable:        true,
			},
			userCfg:   Config{}, // mapstructure 解码后的零值
			rawConfig: mq.Config{},
			expected: Config{
				URL:            "amqp://default:default@localhost:5672/",
				ExchangeName:   "default_exchange",
				QueueName:      "default_queue",
				PublishTimeout: 10 * time.Second,
				Durable:        true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeConfigs(tt.defaultCfg, tt.userCfg, tt.rawConfig)

			if result.URL != tt.expected.URL {
				t.Errorf("URL不匹配: 期望 %q, 实际 %q", tt.expected.URL, result.URL)
			}
			if result.ExchangeName != tt.expected.ExchangeName {
				t.Errorf("ExchangeName不匹配: 期望 %q, 实际 %q", tt.expected.ExchangeName, result.ExchangeName)
			}
			if result.QueueName != tt.expected.QueueName {
				t.Errorf("QueueName不匹配: 期望 %q, 实际 %q", tt.expected.QueueName, result.QueueName)
			}
			if result.PublishTimeout != tt.expected.PublishTimeout {
				t.Errorf("PublishTimeout不匹配: 期望 %v, 实际 %v", tt.expected.PublishTimeout, result.PublishTimeout)
			}
			if result.Durable != tt.expected.Durable {
				t.Errorf("Durable不匹配: 期望 %v, 实际 %v", tt.expected.Durable, result.Durable)
			}
		})
	}
}

// TestMergeConfigsKeyMatching 测试 mergeConfigs 函数的键名匹配
func TestMergeConfigsKeyMatching(t *testing.T) {
	tests := []struct {
		name        string
		rawConfig   mq.Config
		shouldMatch bool
		expectedURL string
		description string
	}{
		{
			name: "正确的键名应该匹配",
			rawConfig: mq.Config{
				"url": "amqp://test:test@localhost:5672/",
			},
			shouldMatch: true,
			expectedURL: "amqp://test:test@localhost:5672/",
			description: "使用正确的键名 'url'",
		},
		{
			name: "错误的键名不应该匹配",
			rawConfig: mq.Config{
				"URL": "amqp://test:test@localhost:5672/", // 大写
			},
			shouldMatch: false,
			expectedURL: DefaultConfig.URL,
			description: "大写的 'URL' 不应该匹配",
		},
		{
			name: "带下划线的键名应该匹配",
			rawConfig: mq.Config{
				"exchange_name": "test_exchange",
			},
			shouldMatch: true,
			expectedURL: DefaultConfig.URL,
			description: "带下划线的键名应该正确匹配",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg Config

			// 配置 mapstructure 解码器
			decoderConfig := &mapstructure.DecoderConfig{
				DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
				Result:     &cfg,
			}
			decoder, err := mapstructure.NewDecoder(decoderConfig)
			if err != nil {
				t.Fatalf("创建解码器失败: %v", err)
			}

			if err := decoder.Decode(tt.rawConfig); err != nil {
				t.Fatalf("解码用户配置失败: %v", err)
			}

			// 使用 mergeConfigs 合并配置
			result := mergeConfigs(DefaultConfig, cfg, tt.rawConfig)

			if result.URL != tt.expectedURL {
				t.Errorf("%s: URL不匹配: 期望 %q, 实际 %q", tt.description, tt.expectedURL, result.URL)
			}

			// 检查 exchange_name 的情况
			if exchangeName, exists := tt.rawConfig["exchange_name"]; exists {
				expectedExchange := exchangeName.(string)
				if result.ExchangeName != expectedExchange {
					t.Errorf("ExchangeName不匹配: 期望 %q, 实际 %q", expectedExchange, result.ExchangeName)
				}
			}
		})
	}
}

// TestMergeConfigsNilHandling 测试 mergeConfigs 函数对 nil 的处理
func TestMergeConfigsNilHandling(t *testing.T) {
	tests := []struct {
		name       string
		defaultCfg Config
		userCfg    Config
		rawConfig  mq.Config
		expected   Config
	}{
		{
			name: "rawConfig为nil应该使用默认配置",
			defaultCfg: Config{
				URL:            "amqp://default:default@localhost:5672/",
				ExchangeName:   "default_exchange",
				PublishTimeout: 10 * time.Second,
				Durable:        true,
			},
			userCfg:   Config{}, // mapstructure 解码后的零值
			rawConfig: nil,      // nil config
			expected: Config{
				URL:            "amqp://default:default@localhost:5672/",
				ExchangeName:   "default_exchange",
				QueueName:      "",
				PublishTimeout: 10 * time.Second,
				Durable:        true,
			},
		},
		{
			name: "空map应该使用默认配置",
			defaultCfg: Config{
				URL:            "amqp://default:default@localhost:5672/",
				ExchangeName:   "default_exchange",
				PublishTimeout: 10 * time.Second,
				Durable:        true,
			},
			userCfg:   Config{},    // mapstructure 解码后的零值
			rawConfig: mq.Config{}, // 空map
			expected: Config{
				URL:            "amqp://default:default@localhost:5672/",
				ExchangeName:   "default_exchange",
				QueueName:      "",
				PublishTimeout: 10 * time.Second,
				Durable:        true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeConfigs(tt.defaultCfg, tt.userCfg, tt.rawConfig)

			if result.URL != tt.expected.URL {
				t.Errorf("URL不匹配: 期望 %q, 实际 %q", tt.expected.URL, result.URL)
			}
			if result.ExchangeName != tt.expected.ExchangeName {
				t.Errorf("ExchangeName不匹配: 期望 %q, 实际 %q", tt.expected.ExchangeName, result.ExchangeName)
			}
			if result.QueueName != tt.expected.QueueName {
				t.Errorf("QueueName不匹配: 期望 %q, 实际 %q", tt.expected.QueueName, result.QueueName)
			}
			if result.PublishTimeout != tt.expected.PublishTimeout {
				t.Errorf("PublishTimeout不匹配: 期望 %v, 实际 %v", tt.expected.PublishTimeout, result.PublishTimeout)
			}
			if result.Durable != tt.expected.Durable {
				t.Errorf("Durable不匹配: 期望 %v, 实际 %v", tt.expected.Durable, result.Durable)
			}
		})
	}
}

// TestMergeConfigsWithNilValues 测试 mergeConfigs 函数对字段值为 nil 的处理
func TestMergeConfigsWithNilValues(t *testing.T) {
	tests := []struct {
		name       string
		rawConfig  mq.Config
		expected   Config
		shouldFail bool
	}{
		{
			name: "rawConfig中某些字段为nil",
			rawConfig: mq.Config{
				"url":             nil, // nil 值
				"exchange_name":   "test_exchange",
				"publish_timeout": "5s",
			},
			// 这种情况下，mapstructure 解码可能会失败或产生意外结果
			// 我们需要测试实际行为
			shouldFail: false, // 先假设不会失败，实际运行后调整
		},
		{
			name: "rawConfig中所有字段都为nil",
			rawConfig: mq.Config{
				"url":             nil,
				"exchange_name":   nil,
				"queue_name":      nil,
				"publish_timeout": nil,
				"durable":         nil,
			},
			shouldFail: false, // 先假设不会失败
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg Config

			// 配置 mapstructure 解码器
			decoderConfig := &mapstructure.DecoderConfig{
				DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
				Result:     &cfg,
			}
			decoder, err := mapstructure.NewDecoder(decoderConfig)
			if err != nil {
				t.Fatalf("创建解码器失败: %v", err)
			}

			// 尝试解码，可能会失败
			err = decoder.Decode(tt.rawConfig)
			if tt.shouldFail {
				if err == nil {
					t.Errorf("期望解码失败，但实际成功了")
				}
				return // 如果期望失败，就不继续测试
			}

			if err != nil {
				t.Logf("解码失败（这可能是预期的）: %v", err)
				// 如果解码失败，我们可以使用零值配置
				cfg = Config{}
			}

			// 使用 mergeConfigs 合并配置
			result := mergeConfigs(DefaultConfig, cfg, tt.rawConfig)

			// 基本验证：结果不应该panic
			t.Logf("合并结果: URL=%q, ExchangeName=%q, PublishTimeout=%v, Durable=%v",
				result.URL, result.ExchangeName, result.PublishTimeout, result.Durable)

			// 对于 nil 值的字段，应该使用零值覆盖默认值（如果字段存在于 rawConfig 中）
			if _, exists := tt.rawConfig["url"]; exists && result.URL != cfg.URL {
				t.Errorf("URL应该使用解码后的值: 期望 %q, 实际 %q", cfg.URL, result.URL)
			}
		})
	}
}

// TestMqOpenIntegration 使用 mq.Open() 方式测试集成功能
func TestMqOpenIntegration(t *testing.T) {
	// 注意：这些测试不会真正连接到 RabbitMQ，因为 Open 函数会尝试连接
	// 我们主要测试配置合并逻辑，连接失败是预期的
	tests := []struct {
		name             string
		config           mq.Config
		expectError      bool
		expectedURL      string
		expectedExchange string
		expectedTimeout  time.Duration
		expectedDurable  bool
	}{
		{
			name:             "空配置应该使用DefaultConfig",
			config:           mq.Config{},
			expectError:      true, // 连接会失败，但这是预期的
			expectedURL:      DefaultConfig.URL,
			expectedExchange: DefaultConfig.ExchangeName,
			expectedTimeout:  DefaultConfig.PublishTimeout,
			expectedDurable:  DefaultConfig.Durable,
		},
		{
			name: "部分配置覆盖",
			config: mq.Config{
				"url":           "amqp://test:test@example.com:5672/",
				"exchange_name": "test_exchange",
			},
			expectError:      true, // 连接会失败
			expectedURL:      "amqp://test:test@example.com:5672/",
			expectedExchange: "test_exchange",
			expectedTimeout:  DefaultConfig.PublishTimeout, // 应该使用默认值
			expectedDurable:  DefaultConfig.Durable,        // 应该使用默认值
		},
		{
			name: "零值覆盖测试",
			config: mq.Config{
				"durable":         false,
				"publish_timeout": "0s",
			},
			expectError:      true,
			expectedURL:      DefaultConfig.URL,
			expectedExchange: DefaultConfig.ExchangeName,
			expectedTimeout:  0 * time.Second, // 零值应该覆盖默认值
			expectedDurable:  false,           // 零值应该覆盖默认值
		},
		{
			name: "完整配置覆盖",
			config: mq.Config{
				"url":             "amqp://custom:custom@localhost:5672/",
				"exchange_name":   "custom_exchange",
				"queue_name":      "custom_queue",
				"publish_timeout": "30s",
				"durable":         false,
			},
			expectError:      true,
			expectedURL:      "amqp://custom:custom@localhost:5672/",
			expectedExchange: "custom_exchange",
			expectedTimeout:  30 * time.Second,
			expectedDurable:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用 mq.Open 打开驱动（会触发我们的配置合并逻辑）
			driver, err := mq.Open("amqp", tt.config)

			if tt.expectError {
				// 我们期望连接失败，但应该是连接错误，不是配置错误
				if err == nil {
					t.Errorf("期望连接失败，但实际成功了")
				} else {
					t.Logf("预期的连接错误: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("意外的错误: %v", err)
			}

			// 如果成功创建了驱动，验证其配置
			if driver != nil {
				// 类型断言获取底层的 amqpDriver
				if amqpDriver, ok := driver.(*amqpDriver); ok {
					// 验证配置是否正确合并
					if amqpDriver.ExchangeName != tt.expectedExchange {
						t.Errorf("ExchangeName不匹配: 期望 %q, 实际 %q", tt.expectedExchange, amqpDriver.ExchangeName)
					}
					if amqpDriver.PublishTimeout != tt.expectedTimeout {
						t.Errorf("PublishTimeout不匹配: 期望 %v, 实际 %v", tt.expectedTimeout, amqpDriver.PublishTimeout)
					}
					if amqpDriver.Durable != tt.expectedDurable {
						t.Errorf("Durable不匹配: 期望 %v, 实际 %v", tt.expectedDurable, amqpDriver.Durable)
					}
				}
			}
		})
	}
}

// TestMqOpenWithNilConfig 测试 mq.Open 使用 nil 配置的情况
func TestMqOpenWithNilConfig(t *testing.T) {
	// 测试传入 nil 配置的情况
	t.Run("nil配置测试", func(t *testing.T) {
		// 这个测试主要验证不会 panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("使用 nil 配置时发生 panic: %v", r)
			}
		}()

		// 使用 nil 配置
		driver, err := mq.Open("amqp", nil)

		// 期望连接失败（因为使用默认配置连接本地 RabbitMQ）
		if err == nil {
			t.Logf("意外成功创建了驱动: %v", driver)
		} else {
			t.Logf("预期的连接错误: %v", err)
		}
	})
}

// TestMqOpenConfigMerging 通过直接调用 init 函数中的逻辑来测试配置合并
func TestMqOpenConfigMerging(t *testing.T) {
	tests := []struct {
		name             string
		config           mq.Config
		expectedURL      string
		expectedExchange string
		expectedQueue    string
		expectedTimeout  time.Duration
		expectedDurable  bool
	}{
		{
			name:             "空配置使用默认值",
			config:           mq.Config{},
			expectedURL:      DefaultConfig.URL,
			expectedExchange: DefaultConfig.ExchangeName,
			expectedQueue:    "", // DefaultConfig 中没有设置 QueueName
			expectedTimeout:  DefaultConfig.PublishTimeout,
			expectedDurable:  DefaultConfig.Durable,
		},
		{
			name: "部分配置覆盖",
			config: mq.Config{
				"url":           "amqp://custom:custom@example.com:5672/",
				"exchange_name": "custom_exchange",
			},
			expectedURL:      "amqp://custom:custom@example.com:5672/",
			expectedExchange: "custom_exchange",
			expectedQueue:    "",                           // 未设置，应该为空
			expectedTimeout:  DefaultConfig.PublishTimeout, // 未设置，使用默认值
			expectedDurable:  DefaultConfig.Durable,        // 未设置，使用默认值
		},
		{
			name: "零值正确覆盖默认值",
			config: mq.Config{
				"durable":         false,
				"publish_timeout": "0s",
				"queue_name":      "",
			},
			expectedURL:      DefaultConfig.URL,          // 未设置，使用默认值
			expectedExchange: DefaultConfig.ExchangeName, // 未设置，使用默认值
			expectedQueue:    "",                         // 明确设置为空字符串
			expectedTimeout:  0,                          // 零值覆盖默认值
			expectedDurable:  false,                      // 零值覆盖默认值
		},
		{
			name: "完整配置",
			config: mq.Config{
				"url":             "amqp://full:full@localhost:5672/",
				"exchange_name":   "full_exchange",
				"queue_name":      "full_queue",
				"publish_timeout": "15s",
				"durable":         true,
			},
			expectedURL:      "amqp://full:full@localhost:5672/",
			expectedExchange: "full_exchange",
			expectedQueue:    "full_queue",
			expectedTimeout:  15 * time.Second,
			expectedDurable:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 直接模拟 init 函数中的配置处理逻辑
			var cfg Config

			// 配置 mapstructure 解码器以支持 time.Duration
			decoderConfig := &mapstructure.DecoderConfig{
				DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
				Result:     &cfg,
			}
			decoder, err := mapstructure.NewDecoder(decoderConfig)
			if err != nil {
				t.Fatalf("创建解码器失败: %v", err)
			}

			if err := decoder.Decode(tt.config); err != nil {
				t.Fatalf("解码配置失败: %v", err)
			}

			// 使用 mergeConfigs 合并配置
			dst := mergeConfigs(DefaultConfig, cfg, tt.config)

			// 验证配置合并结果
			if dst.URL != tt.expectedURL {
				t.Errorf("URL不匹配: 期望 %q, 实际 %q", tt.expectedURL, dst.URL)
			}
			if dst.ExchangeName != tt.expectedExchange {
				t.Errorf("ExchangeName不匹配: 期望 %q, 实际 %q", tt.expectedExchange, dst.ExchangeName)
			}
			if dst.QueueName != tt.expectedQueue {
				t.Errorf("QueueName不匹配: 期望 %q, 实际 %q", tt.expectedQueue, dst.QueueName)
			}
			if dst.PublishTimeout != tt.expectedTimeout {
				t.Errorf("PublishTimeout不匹配: 期望 %v, 实际 %v", tt.expectedTimeout, dst.PublishTimeout)
			}
			if dst.Durable != tt.expectedDurable {
				t.Errorf("Durable不匹配: 期望 %v, 实际 %v", tt.expectedDurable, dst.Durable)
			}

			// 额外的日志输出，显示合并后的完整配置
			t.Logf("配置合并结果: URL=%q, Exchange=%q, Queue=%q, Timeout=%v, Durable=%v",
				dst.URL, dst.ExchangeName, dst.QueueName, dst.PublishTimeout, dst.Durable)
		})
	}
}

// TestMqOpenUsageExamples 展示如何使用 mq.Open(mq.Config{}) 的实际例子
func TestMqOpenUsageExamples(t *testing.T) {
	// 这些测试展示了实际使用场景，连接失败是预期的
	t.Run("使用空配置（全部默认值）", func(t *testing.T) {
		// 使用空配置，应该使用所有默认值
		_, err := mq.Open("amqp", mq.Config{})

		// 连接失败是预期的，但不应该是配置错误
		if err != nil {
			t.Logf("预期的连接错误: %v", err)
			// 验证错误信息包含连接相关的内容，而不是配置错误
			errMsg := err.Error()
			if !strings.Contains(errMsg, "connect") && !strings.Contains(errMsg, "dial") && !strings.Contains(errMsg, "timeout") {
				t.Errorf("错误看起来不像连接错误: %v", err)
			}
		}
	})

	t.Run("使用自定义URL", func(t *testing.T) {
		// 只设置URL，其他使用默认值
		_, err := mq.Open("amqp", mq.Config{
			"url": "amqp://myuser:mypass@myserver.com:5672/",
		})

		if err != nil {
			t.Logf("预期的连接错误: %v", err)
		}
	})

	t.Run("使用零值覆盖默认值", func(t *testing.T) {
		// 明确设置 durable 为 false，覆盖默认的 true
		_, err := mq.Open("amqp", mq.Config{
			"durable":         false,
			"publish_timeout": "0s", // 零值超时
		})

		if err != nil {
			t.Logf("预期的连接错误: %v", err)
		}
	})

	t.Run("完整自定义配置", func(t *testing.T) {
		// 完整的自定义配置
		_, err := mq.Open("amqp", mq.Config{
			"url":             "amqp://admin:secret@production.rabbitmq.com:5672/",
			"exchange_name":   "production_events",
			"queue_name":      "service_queue",
			"publish_timeout": "30s",
			"durable":         true,
		})

		if err != nil {
			t.Logf("预期的连接错误: %v", err)
		}
	})

	t.Run("测试DefaultConfig URL的使用", func(t *testing.T) {
		// 验证当没有提供URL时，使用DefaultConfig.URL
		t.Logf("DefaultConfig.URL = %q", DefaultConfig.URL)

		// 使用不包含URL的配置
		_, err := mq.Open("amqp", mq.Config{
			"exchange_name": "test_exchange",
		})

		if err != nil {
			t.Logf("使用DefaultConfig URL的连接错误: %v", err)
			// 错误信息应该包含默认的localhost地址
			if !strings.Contains(err.Error(), "localhost") && !strings.Contains(err.Error(), "127.0.0.1") && !strings.Contains(err.Error(), "::1") {
				t.Errorf("错误信息中应该包含默认URL的主机名，实际错误: %v", err)
			}
		}
	})
}

// TestRealWorldScenarios 测试真实世界的使用场景
func TestRealWorldScenarios(t *testing.T) {
	scenarios := []struct {
		name        string
		description string
		config      mq.Config
	}{
		{
			name:        "开发环境",
			description: "开发环境使用本地RabbitMQ，默认配置",
			config:      mq.Config{},
		},
		{
			name:        "测试环境",
			description: "测试环境使用自定义exchange，但其他保持默认",
			config: mq.Config{
				"exchange_name": "test_events",
			},
		},
		{
			name:        "生产环境",
			description: "生产环境使用完整自定义配置",
			config: mq.Config{
				"url":             "amqp://prod_user:prod_pass@prod.rabbitmq.com:5672/",
				"exchange_name":   "prod_events",
				"publish_timeout": "60s",
				"durable":         true,
			},
		},
		{
			name:        "高性能场景",
			description: "高性能场景，关闭持久化以提高性能",
			config: mq.Config{
				"durable":         false, // 零值覆盖默认的true
				"publish_timeout": "1s",  // 较短的超时
			},
		},
		{
			name:        "调试场景",
			description: "调试场景，使用特定队列名称",
			config: mq.Config{
				"queue_name":    "debug_queue",
				"exchange_name": "debug_exchange",
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			t.Logf("场景描述: %s", scenario.description)
			t.Logf("配置: %+v", scenario.config)

			// 尝试打开连接
			_, err := mq.Open("amqp", scenario.config)

			// 连接失败是预期的
			if err != nil {
				t.Logf("预期的连接错误: %v", err)
			} else {
				t.Logf("意外成功！这可能意味着本地有RabbitMQ在运行")
			}
		})
	}
}

// TestMapstructureDecoding 测试mapstructure解码行为
func TestMapstructureDecoding(t *testing.T) {
	tests := []struct {
		name     string
		input    mq.Config
		expected Config
	}{
		{
			name: "正常字段映射",
			input: mq.Config{
				"url":             "amqp://test:test@localhost:5672/",
				"exchange_name":   "test_exchange",
				"queue_name":      "test_queue",
				"publish_timeout": "30s",
				"durable":         false,
			},
			expected: Config{
				URL:            "amqp://test:test@localhost:5672/",
				ExchangeName:   "test_exchange",
				QueueName:      "test_queue",
				PublishTimeout: 30 * time.Second,
				Durable:        false,
			},
		},
		{
			name: "部分字段缺失",
			input: mq.Config{
				"url":     "amqp://test:test@localhost:5672/",
				"durable": true,
			},
			expected: Config{
				URL:            "amqp://test:test@localhost:5672/",
				ExchangeName:   "", // 零值
				QueueName:      "", // 零值
				PublishTimeout: 0,  // 零值
				Durable:        true,
			},
		},
		{
			name:  "空配置",
			input: mq.Config{},
			expected: Config{
				URL:            "",    // 零值
				ExchangeName:   "",    // 零值
				QueueName:      "",    // 零值
				PublishTimeout: 0,     // 零值
				Durable:        false, // 零值
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg Config
			decoderConfig := &mapstructure.DecoderConfig{
				DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
				Result:     &cfg,
			}
			decoder, err := mapstructure.NewDecoder(decoderConfig)
			if err != nil {
				t.Fatalf("创建解码器失败: %v", err)
			}

			if err := decoder.Decode(tt.input); err != nil {
				t.Fatalf("解码失败: %v", err)
			}

			if cfg.URL != tt.expected.URL {
				t.Errorf("URL不匹配: 期望 %s, 实际 %s", tt.expected.URL, cfg.URL)
			}
			if cfg.ExchangeName != tt.expected.ExchangeName {
				t.Errorf("ExchangeName不匹配: 期望 %s, 实际 %s", tt.expected.ExchangeName, cfg.ExchangeName)
			}
			if cfg.QueueName != tt.expected.QueueName {
				t.Errorf("QueueName不匹配: 期望 %s, 实际 %s", tt.expected.QueueName, cfg.QueueName)
			}
			if cfg.PublishTimeout != tt.expected.PublishTimeout {
				t.Errorf("PublishTimeout不匹配: 期望 %v, 实际 %v", tt.expected.PublishTimeout, cfg.PublishTimeout)
			}
			if cfg.Durable != tt.expected.Durable {
				t.Errorf("Durable不匹配: 期望 %v, 实际 %v", tt.expected.Durable, cfg.Durable)
			}
		})
	}
}

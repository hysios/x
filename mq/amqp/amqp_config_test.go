package amqp

import (
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

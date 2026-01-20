const defaultConfig = {
  input: {
    name: '',
    type: 1,
    key: '',
    base_url: '',
    other: '',
    proxy: '',
    test_model: '',
    model_mapping: [],
    model_headers: [],
    custom_parameter: '',
    models: [],
    groups: ['default'],
    plugin: {},
    tag: '',
    only_chat: false,
    pre_cost: 1,
    disabled_stream: [],
    compatible_response: false,
    allow_extra_body: false
  },
  inputLabel: {
    name: '渠道名稱',
    type: '渠道類型',
    base_url: '渠道API地址',
    key: '密鑰',
    other: '其他參數',
    proxy: '代理地址',
    test_model: '測速模型',
    models: '模型',
    model_mapping: '模型映射關係',
    model_headers: '自定義模型請求頭',
    custom_parameter: '額外參數',
    groups: '用戶組',
    only_chat: '僅支持聊天',
    tag: '標籤',
    provider_models_list: '',
    pre_cost: '預計費選項',
    disabled_stream: '禁用流式的模型',
    compatible_response: '兼容Response API',
    allow_extra_body: '允許額外字段透傳'
  },
  prompt: {
    type: '請選擇渠道類型',
    name: '請為渠道命名',
    base_url: '可空，請輸入中轉API地址，例如通過cloudflare中轉',
    key: '請輸入渠道對應的鑑權密鑰',
    other: '',
    proxy:
      '單獨設置代理地址，支持http和socks5，例如：http://127.0.0.1:1080,代理地址中可以通過 `%s` 作為會話標識佔位符，程序中檢測到有佔位符會根據Key生成唯一會話標識符進行替換',
    test_model: '用於測試使用的模型，為空時無法測速,如：gpt-3.5-turbo，僅支持chat模型',
    models:
      '請選擇該渠道所支持的模型,你也可以輸入通配符*來匹配模型，例如：gpt-3.5*，表示支持所有gpt-3.5開頭的模型，*號只能在最後一位使用，前面必須有字符，例如：gpt-3.5*是正確的，*gpt-3.5是錯誤的',
    model_mapping: '模型映射關係：例如用戶請求A模型，實際轉發給渠道的模型為B。在B模型加前綴+，表示使用傳入模型計費，例如：+gpt-3.5-turbo',
    model_headers: '自定義模型請求頭，例如：{"key": "value"}',
    custom_parameter:
      '支持通過 JSON 注入額外參數（可嵌套）。可用控制項：overwrite：設為 true 覆蓋同名字段，未設置或 false 時僅補充缺失字段；per_model：設為 true 後按模型名進行參數覆蓋，如 {"per_model":true,"gpt-3.5-turbo":{"temperature": 0.7},"gpt-4":{"temperature": 0.5}}；pre_add：設為 true 時在請求入口即完成參數覆蓋，否則會在發送請求前再進行參數覆蓋，適用於所有渠道（含 Claude、Gemini），如 {"pre_add":true,"overwrite":true,"stream":false}。',
    groups: '請選擇該渠道所支持的用戶組',
    only_chat: '如果選擇了僅支持聊天，那麼遇到有函數調用的請求會跳過該渠道',
    provider_models_list: '必須填寫所有數據後才能獲取模型列表',
    tag: '你可以為你的渠道打一個標籤，打完標籤後，可以通過標籤進行批量管理渠道，注意：設置標籤後某些設置只能通過渠道標籤修改，無法在渠道列表中修改。',
    pre_cost:
      '這裡選擇預計費選項，用於預估費用，如果你覺得計算圖片佔用太多資源，可以選擇關閉圖片計費。但是請注意：有些渠道在stream下是不會返回tokens的，這會導致輸入tokens計算錯誤。',
    disabled_stream: '這裡填寫禁用流式的模型，注意：如果填寫了禁用流式的模型，那麼這些模型在流式請求時會跳過該渠道',
    compatible_response: '兼容Response API',
    allow_extra_body: '開啟後，將會透傳用戶請求中的額外字段（如OpenAI SDK的extra_body參數），適用於需要傳遞自定義參數到上游API的場景'
  },
  modelGroup: 'OpenAI'
};

const typeConfig = {
  1: {
    inputLabel: {
      provider_models_list: '從OpenAI獲取模型列表'
    }
  },
  8: {
    inputLabel: {
      provider_models_list: '從渠道獲取模型列表'
    },
    prompt: {
      other: ''
    }
  },
  3: {
    inputLabel: {
      base_url: 'AZURE_OPENAI_ENDPOINT',
      other: '默認 API 版本',
      provider_models_list: '從Azure獲取已部署模型列表'
    },
    prompt: {
      base_url: '請填寫AZURE_OPENAI_ENDPOINT',
      other: '請輸入默認API版本，例如：2024-05-01-preview'
    }
  },
  55: {
    inputLabel: {
      base_url: 'AZURE_OPENAI_ENDPOINT',
      other: '默認 API 版本',
      provider_models_list: '從Azure獲取已部署模型列表'
    },
    prompt: {
      base_url: '請填寫AZURE_OPENAI_ENDPOINT',
      other: '請輸入默認API版本，例如：preview OR latest'
    }
  },
  11: {
    input: {
      models: ['PaLM-2'],
      test_model: 'PaLM-2'
    },
    modelGroup: 'Google PaLM'
  },
  14: {
    inputLabel: {
      provider_models_list: '從Claude獲取模型列表'
    },
    input: {
      models: [
        'claude-instant-1.2',
        'claude-2.0',
        'claude-2.1',
        'claude-3-opus-20240229',
        'claude-3-sonnet-20240229',
        'claude-3-haiku-20240307'
      ],
      test_model: 'claude-3-haiku-20240307'
    },
    modelGroup: 'Anthropic'
  },
  15: {
    input: {
      models: [
        'ERNIE-4.0-Turbo-8K',
        'ERNIE-4.0-8K-Latest',
        'ERNIE-4.0-8K-0613',
        'ERNIE-3.5-8K-0613',
        'ERNIE-Bot-turbo',
        'ERNIE-Lite-8K-0922',
        'ERNIE-Lite-8K',
        'ERNIE-Lite-8K-0308',
        'ERNIE-3.5-8K',
        'ERNIE-Bot',
        'ERNIE-4.0-8K',
        'ERNIE-4.0-8K-Preview',
        'ERNIE-4.0-8K-Preview-0518',
        'ERNIE-4.0-8K-0329',
        'ERNIE-4.0-8K-0104',
        'ERNIE-Bot-4',
        'ERNIE-Bot-8k',
        'ERNIE-3.5-128K',
        'ERNIE-3.5-8K-preview',
        'ERNIE-3.5-8K-0329',
        'ERNIE-3.5-4K-0205',
        'ERNIE-3.5-8K-0205',
        'ERNIE-3.5-8K-1222',
        'ERNIE-Speed',
        'ERNIE-Speed-8K',
        'ERNIE-Speed-128K',
        'ERNIE-Tiny-8K',
        'ERNIE-Function-8K',
        'ERNIE-Character-8K',
        'ERNIE-Character-Fiction-8K',
        'ERNIE-Bot-turbo-AI',
        'Embedding-V1'
      ],
      test_model: 'ERNIE-Speed'
    },
    prompt: {
      key: '按照如下格式輸入：APIKey|SecretKey, 如果開啟了OpenAI API，請直接輸入APIKEY'
    },
    modelGroup: 'Baidu'
  },
  16: {
    input: {
      models: ['glm-3-turbo', 'glm-4', 'glm-4v', 'embedding-2', 'cogview-3'],
      test_model: 'glm-3-turbo'
    },
    modelGroup: 'Zhipu'
  },
  17: {
    inputLabel: {
      other: '插件參數',
      provider_models_list: '從Ali獲取模型列表'
    },
    input: {
      models: ['qwen-turbo', 'qwen-plus', 'qwen-max', 'qwen-max-longcontext', 'text-embedding-v1'],
      test_model: 'qwen-turbo'
    },
    prompt: {
      other: '請輸入插件參數，即 X-DashScope-Plugin 請求頭的取值'
    },
    modelGroup: 'Ali'
  },
  18: {
    inputLabel: {
      other: '版本號'
    },
    input: {
      models: ['SparkDesk', 'SparkDesk-v1.1', 'SparkDesk-v2.1', 'SparkDesk-v3.1', 'SparkDesk-v3.5']
    },
    prompt: {
      key: '按照如下格式輸入：APPID|APISecret|APIKey',
      other: '請輸入版本號，例如：v3.1'
    },
    modelGroup: 'Xunfei'
  },
  19: {
    input: {
      models: ['360GPT_S2_V9', 'embedding-bert-512-v1', 'embedding_s1_v1', 'semantic_similarity_s1_v1'],
      test_model: '360GPT_S2_V9'
    },
    modelGroup: '360'
  },
  22: {
    prompt: {
      key: '按照如下格式輸入：APIKey-AppId，例如：fastgpt-0sp2gtvfdgyi4k30jwlgwf1i-64f335d84283f05518e9e041'
    }
  },
  23: {
    input: {
      models: ['ChatStd', 'ChatPro'],
      test_model: 'ChatStd'
    },
    prompt: {
      key: '按照如下格式輸入：AppId|SecretId|SecretKey'
    },
    modelGroup: 'Tencent'
  },
  25: {
    inputLabel: {
      other: '版本號',
      provider_models_list: '從Gemini獲取模型列表'
    },
    input: {
      models: ['gemini-pro', 'gemini-pro-vision', 'gemini-1.0-pro', 'gemini-1.5-pro'],
      test_model: 'gemini-pro'
    },
    prompt: {
      other: '請輸入版本號，例如：v1'
    },
    modelGroup: 'Google Gemini'
  },
  26: {
    input: {
      models: ['Baichuan2-Turbo', 'Baichuan2-Turbo-192k', 'Baichuan2-53B', 'Baichuan-Text-Embedding'],
      test_model: 'Baichuan2-Turbo'
    },
    modelGroup: 'Baichuan'
  },
  24: {
    inputLabel: {
      other: '位置/區域'
    },
    input: {
      models: ['tts-1', 'tts-1-hd']
    },
    prompt: {
      test_model: '',
      base_url: '',
      other: '請輸入你 Speech Studio 的位置/區域，例如：eastasia'
    }
  },
  27: {
    input: {
      models: ['abab6.5s-chat', 'MiniMax-Text-01', 'speech-01-turbo', 'speech-01-240228', 'speech-01-turbo-240228'],
      test_model: 'abab6.5s-chat'
    },
    modelGroup: 'MiniMax'
  },
  28: {
    input: {
      models: ['deepseek-coder', 'deepseek-chat'],
      test_model: 'deepseek-chat'
    },
    inputLabel: {
      provider_models_list: '從Deepseek獲取模型列表'
    },
    modelGroup: 'Deepseek'
  },
  29: {
    inputLabel: {
      provider_models_list: '從Moonshot獲取模型列表'
    },
    input: {
      models: ['moonshot-v1-8k', 'moonshot-v1-32k', 'moonshot-v1-128k'],
      test_model: 'moonshot-v1-8k'
    },
    modelGroup: 'Moonshot'
  },
  30: {
    input: {
      models: [
        'open-mistral-7b',
        'open-mixtral-8x7b',
        'mistral-small-latest',
        'mistral-medium-latest',
        'mistral-large-latest',
        'mistral-embed'
      ],
      test_model: 'open-mistral-7b'
    },
    inputLabel: {
      provider_models_list: '從Mistral獲取模型列表'
    },
    modelGroup: 'Mistral'
  },
  31: {
    input: {
      models: ['llama2-7b-2048', 'llama2-70b-4096', 'mixtral-8x7b-32768', 'gemma-7b-it'],
      test_model: 'llama2-7b-2048'
    },
    inputLabel: {
      provider_models_list: '從Groq獲取模型列表'
    },
    modelGroup: 'Groq'
  },
  32: {
    input: {
      models: [
        'claude-instant-1.2',
        'claude-2.0',
        'claude-2.1',
        'claude-3-opus-20240229',
        'claude-3-sonnet-20240229',
        'claude-3-haiku-20240307'
      ],
      test_model: 'claude-3-haiku-20240307'
    },
    prompt: {
      key: '老版本Bedrock按照如下格式輸入：Region|AccessKeyID|SecretAccessKey|SessionToken 其中SessionToken可不填空,新版本Bedrock按照如下格式輸入：Region|Token(其中Token不能為空，Token前往新版本Bedrock控制檯創建API密鑰)'
    },
    modelGroup: 'Anthropic'
  },
  33: {
    input: {
      models: ['yi-34b-chat-0205', 'yi-34b-chat-200k', 'yi-vl-plus'],
      test_model: 'yi-34b-chat-0205'
    },
    modelGroup: 'Lingyiwanwu'
  },
  34: {
    input: {
      models: [
        'mj_imagine',
        'mj_variation',
        'mj_reroll',
        'mj_blend',
        'mj_modal',
        'mj_zoom',
        'mj_shorten',
        'mj_high_variation',
        'mj_low_variation',
        'mj_pan',
        'mj_inpaint',
        'mj_custom_zoom',
        'mj_describe',
        'mj_upscale',
        'swap_face',
        'mj_upload'
      ]
    },
    prompt: {
      key: '密鑰填寫midjourney-proxy的密鑰，如果沒有設置密鑰，可以隨便填',
      base_url: '地址填寫midjourney-proxy部署的地址',
      test_model: '',
      model_mapping: ''
    },
    modelGroup: 'Midjourney'
  },
  35: {
    input: {
      models: [
        '@cf/stabilityai/stable-diffusion-xl-base-1.0',
        '@cf/lykon/dreamshaper-8-lcm',
        '@cf/bytedance/stable-diffusion-xl-lightning',
        '@cf/qwen/qwen1.5-7b-chat-awq',
        '@cf/qwen/qwen1.5-14b-chat-awq',
        '@hf/google/gemma-7b-it',
        '@hf/thebloke/deepseek-coder-6.7b-base-awq',
        '@hf/thebloke/llama-2-13b-chat-awq',
        '@cf/openai/whisper'
      ],
      test_model: '@hf/google/gemma-7b-it'
    },
    prompt: {
      key: '按照如下格式輸入：CLOUDFLARE_ACCOUNT_ID|CLOUDFLARE_API_TOKEN',
      base_url: ''
    },
    modelGroup: 'Cloudflare AI'
  },
  36: {
    input: {
      models: ['command-r', 'command-r-plus'],
      test_model: 'command-r'
    },
    inputLabel: {
      provider_models_list: '從Cohere獲取模型列表'
    },
    modelGroup: 'Cohere'
  },
  37: {
    input: {
      models: ['sd3', 'sd3-turbo', 'stable-image-core']
    },
    prompt: {
      test_model: ''
    },
    modelGroup: 'Stability AI'
  },
  38: {
    input: {
      models: ['coze-*']
    },
    prompt: {
      models: '模型名稱為coze-{bot_id}，你也可以直接使用 coze-* 通配符來匹配所有coze開頭的模型',
      model_mapping:
        '模型名稱映射， 你可以取一個容易記憶的名字來代替coze-{bot_id}，例如：{"coze-translate": "coze-xxxxx"},注意：如果使用了模型映射，那麼上面的模型名稱必須使用映射前的名稱，上述例子中，你應該在模型中填入coze-translate(如果已經使用了coze-*，可以忽略)。'
    },
    modelGroup: 'Coze'
  },
  39: {
    inputLabel: {
      provider_models_list: '從Ollama獲取模型'
    },
    input: {
      base_url: 'https://ollama.com',
      models: [
        'glm-4.6',
        'kimi-k2:1t',
        'qwen3-coder:480b',
        'deepseek-v3.1:671b',
        'gpt-oss:120b',
        'gpt-oss:20b',
        'qwen3-vl:235b',
        'minimax-m2'
      ]
    },
    prompt: {
      base_url:
        '請輸入你部署的Ollama地址或者Ollama Cloud地址，例如：http://127.0.0.1:11434或者https://ollama.com，如果你使用了cloudflare Zero Trust，可以在下方Header配置填入授權信息',
      key: '本地部署可以隨便填，Ollama Cloud請填寫API KEY，獲取地址https://ollama.com/settings/keys'
    }
  },
  40: {
    input: {
      models: ['hunyuan-lite', 'hunyuan-pro', 'hunyuan-standard-256K', 'hunyuan-standard'],
      test_model: 'hunyuan-lite'
    },
    prompt: {
      key: '按照如下格式輸入：SecretId|SecretKey'
    },
    modelGroup: 'Hunyuan'
  },
  41: {
    input: {
      models: ['suno_lyrics', 'chirp-v3-0', 'chirp-v3-5']
    },
    prompt: {
      key: '密鑰填寫Suno-API的密鑰，如果沒有設置密鑰，可以隨便填',
      base_url: '地址填寫Suno-API部署的地址',
      test_model: '',
      model_mapping: ''
    },
    modelGroup: 'Suno'
  },
  42: {
    input: {
      models: ['claude-3-opus-20240229', 'claude-3-sonnet-20240229', 'claude-3-haiku-20240307']
    },
    prompt: {
      key: '請參考wiki中的文檔獲取key. https://github.com/bentwnghk/one-api/wiki/VertexAI',
      other: 'Region|ProjectID',
      base_url: ''
    },
    modelGroup: 'VertexAI'
  },
  45: {
    input: {
      base_url: '',
      models: ['black-forest-labs/FLUX.1-dev', 'black-forest-labs/FLUX.1-schnell']
    },
    inputLabel: {
      base_url: '渠道API地址',
      provider_models_list: '從Siliconflow獲取模型列表'
    },
    prompt: {
      base_url: '官方api地址https://api.siliconflow.com即將停用，請使用https://api.siliconflow.cn'
    },
    modelGroup: 'Siliconflow'
  },
  47: {
    input: {
      models: ['jina-reranker-v2-base-multilingual']
    },
    prompt: {
      test_model: ''
    },
    modelGroup: 'Jina'
  },
  49: {
    input: {
      models: ['gpt-4o', 'gpt-4o-mini', 'text-embedding-3-large', 'text-embedding-3-small', 'Cohere-command-r-plus', 'Cohere-command-r'],
      test_model: 'gpt-4o-mini'
    },
    inputLabel: {
      provider_models_list: '從Github獲取模型列表'
    },
    prompt: {
      key: '密鑰信息請參考https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens',
      base_url: 'https://models.inference.ai.azure.com'
    },
    modelGroup: 'Github'
  },
  51: {
    input: {
      models: [
        'recraftv3',
        'recraft20b',
        'recraft_vectorize',
        'recraft_removeBackground',
        'recraft_clarityUpscale',
        'recraft_generativeUpscale',
        'recraft_styles'
      ]
    }
  },
  53: {
    input: {
      models: [
        'kling-video_kling-v1_std_5',
        'kling-video_kling-v1_std_10',
        'kling-video_kling-v1_pro_5',
        'kling-video_kling-v1_pro_10',

        'kling-video_kling-v1-5_std_5',
        'kling-video_kling-v1-5_std_10',
        'kling-video_kling-v1-5_pro_5',
        'kling-video_kling-v1-5_pro_10',

        'kling-video_kling-v1-10_std_5',
        'kling-video_kling-v1-10_std_10',
        'kling-video_kling-v1-10_pro_5',
        'kling-video_kling-v1-10_pro_10'
      ]
    },
    prompt: {
      key: '官方密鑰格式： accessKey|secretKey'
    },
    modelGroup: 'Kling'
  },
  54: {
    inputLabel: {
      base_url: 'Azure Databricks Endpoint',
      key: 'DATABRICKS_TOKEN'
    },
    prompt: {
      base_url: '請填寫Azure Databricks Endpoint',
      key: '請輸入DATABRICKS_TOKEN'
    }
  },
  20: {
    inputLabel: {
      provider_models_list: '從OR獲取模型列表'
    }
  }
};

export { defaultConfig, typeConfig };

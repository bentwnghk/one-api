const PaymentType = {
  epay: '易支付',
  alipay: '支付宝',
  wxpay: '微信支付',
  stripe: 'Stripe',
};

const CurrencyType = {
  CNY: '人民幣',
  USD: '美元'
};

const PaymentConfig = {
  epay: {
    pay_domain: {
      name: '支付域名',
      description: '支付域名',
      type: 'text',
      value: ''
    },
    partner_id: {
      name: '商戶號',
      description: '商戶號',
      type: 'text',
      value: ''
    },
    key: {
      name: '密鑰',
      description: '密鑰',
      type: 'text',
      value: ''
    },
    pay_type: {
      name: '支付類型',
      description: '支付類型，如果需要跳轉到易支付收銀台，請選擇收銀台',
      type: 'select',
      value: '',
      options: [
        {
          name: '收銀台',
          value: ''
        },
        {
          name: '支付宝',
          value: 'alipay'
        },
        {
          name: '微信',
          value: 'wxpay'
        },
        {
          name: 'QQ',
          value: 'qqpay'
        },
        {
          name: '京東',
          value: 'jdpay'
        },
        {
          name: '銀聯',
          value: 'bank'
        },
        {
          name: 'PayPal',
          value: 'paypal'
        },
        {
          name: 'USDT',
          value: 'usdt'
        }
      ]
    }
  },
  alipay: {
    app_id: {
      name: '應用ID',
      description: '支付宝應用ID',
      type: 'text',
      value: ''
    },
    private_key: {
      name: '應用私鑰',
      description: '應用私鑰，開發者自己生成，詳細參考官方文檔 https://opendocs.alipay.com/common/02kipl?pathHash=84adb0fd',
      type: 'text',
      value: ''
    },
    public_key: {
      name: '支付宝公鑰',
      description: '支付宝公鑰，詳細參考官方文檔 https://opendocs.alipay.com/common/02kdnc?pathHash=fb0c752a',
      type: 'text',
      value: ''
    },
    pay_type: {
      name: '支付類型',
      description: '支付類型，需要您再支付宝開發者中心開通相關權限才可以使用對應類型支付方式',
      type: 'select',
      value: '',
      options: [
        {
          name: '當面付',
          value: 'facepay'
        },
        {
          name: '電腦網站支付',
          value: 'pagepay'
        },
        {
          name: '手機網站支付',
          value: 'wappay'
        }
      ]
    }
  },
  wxpay: {
    app_id: {
      name: 'AppID',
      description: '應用ID 詳見 https://pay.weixin.qq.com/wiki/doc/apiv3/open/pay/chapter2_7_1.shtml',
      type: 'text',
      value: ''
    },
    mch_id: {
      name: '商戶號',
      description: '微信商戶號 詳見 https://pay.weixin.qq.com/wiki/doc/apiv3/open/pay/chapter2_7_1.shtml',
      type: 'text',
      value: ''
    },
    mch_certificate_serial_number: {
      name: '商戶証書序列號',
      description: '商戶証書序列號 詳見 https://pay.weixin.qq.com/wiki/doc/apiv3/open/pay/chapter2_7_1.shtml',
      type: 'text',
      value: ''
    },
    mch_apiv3_key: {
      name: '商戶APIv3密鑰',
      description: '商戶APIv3密鑰 詳見 https://pay.weixin.qq.com/wiki/doc/apiv3/open/pay/chapter2_7_1.shtml',
      type: 'text',
      value: ''
    },
    mch_private_key: {
      name: '商戶私鑰',
      description: '商戶私鑰 詳見 https://pay.weixin.qq.com/wiki/doc/apiv3/open/pay/chapter2_7_1.shtml',
      type: 'text',
      value: ''
    },
    pay_type: {
      name: '支付類型',
      description: '支付類型',
      type: 'select',
      value: '',
      options: [
        {
          name: 'Native 支付',
          value: 'Native'
        }
      ]
    },
  },
  stripe: {
    secret_key: {
      name: 'SecretKey',
      description: 'API 私钥',
      type: 'text',
      value: ''
    },
    webhook_secret: {
      name: 'WebHookSecret',
      description: '回调验证密钥，不用填写，创建网关后会自动在stripe后台创建webhook并获取webhook密钥',
      type: 'text',
      value: ''
    },
  }
};

export { PaymentConfig, PaymentType, CurrencyType };

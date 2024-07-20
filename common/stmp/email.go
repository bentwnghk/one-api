package stmp

import (
	"fmt"
	"one-api/common"
	"one-api/common/config"
	"one-api/common/utils"
	"strings"

	"github.com/wneessen/go-mail"
)

type StmpConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewStmp(host string, port int, username string, password string, from string) *StmpConfig {
	if from == "" {
		from = username
	}

	return &StmpConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

func (s *StmpConfig) Send(to, subject, body string) error {
	message := mail.NewMsg()
	message.From(s.From)
	message.To(to)
	message.Subject(subject)
	message.SetGenHeader("References", s.getReferences())
	message.SetBodyString(mail.TypeTextHTML, body)
	message.SetUserAgent(fmt.Sprintf("Mr.🆖 AI Hub %s // https://api.mr5ai.com", config.Version))

	client, err := mail.NewClient(
		s.Host,
		mail.WithPort(s.Port),
		mail.WithUsername(s.Username),
		mail.WithPassword(s.Password),
		mail.WithSMTPAuth(mail.SMTPAuthPlain))

	if err != nil {
		return err
	}

	switch s.Port {
	case 465:
		client.SetSSL(true)
	case 587:
		client.SetTLSPolicy(mail.TLSMandatory)
		client.SetSMTPAuth(mail.SMTPAuthLogin)
	}

	if err := client.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

func (s *StmpConfig) getReferences() string {
	froms := strings.Split(s.From, "@")
	return fmt.Sprintf("<%s.%s@%s>", froms[0], utils.GetUUID(), froms[1])
}

func (s *StmpConfig) Render(to, subject, content string) error {
	body := getDefaultTemplate(content)

	return s.Send(to, subject, body)
}

func GetSystemStmp() (*StmpConfig, error) {
	if config.SMTPServer == "" || config.SMTPPort == 0 || config.SMTPAccount == "" || config.SMTPToken == "" {
		return nil, fmt.Errorf("SMTP 信息未配置")
	}

	return NewStmp(config.SMTPServer, config.SMTPPort, config.SMTPAccount, config.SMTPToken, config.SMTPFrom), nil
}

func SendPasswordResetEmail(userName, email, link string) error {
	stmp, err := GetSystemStmp()

	if err != nil {
		return err
	}

	contentTemp := `<p style="font-size: 30px">Hi <strong>%s,</strong></p>
	<p>
		您正在進行密碼重置。點擊下方按鈕以重置密碼。
	</p>
	
	<p style="text-align: center; font-size: 13px;">
		<a target="__blank" href="%s" class="button" style="color: #ffffff;">重置密碼</a>
	</p>
	
	<p style="color: #858585; padding-top: 15px;">
		如果鏈接無法點擊，請嘗試點擊下面的鏈接或將其複製到瀏覽器中打開<br> %s
	</p>
	<p style="color: #858585;">重置鏈接 %d 分鐘內有效，如果不是本人操作，請忽略。</p>`

	subject := fmt.Sprintf("%s密碼重置", config.SystemName)
	content := fmt.Sprintf(contentTemp, userName, link, link, common.VerificationValidMinutes)

	return stmp.Render(email, subject, content)
}

func SendVerificationCodeEmail(email, code string) error {
	stmp, err := GetSystemStmp()

	if err != nil {
		return err
	}

	contentTemp := `
	<p>
		您正在進行電郵地址驗証。您的驗証碼為: 
	</p>
	
	<p style="text-align: center; font-size: 30px; color: #58a6ff;">
		<strong>%s</strong>
	</p>
	
	<p style="color: #858585; padding-top: 15px;">
		驗証碼 %d 分鐘內有效，如果不是本人操作，請忽略。
	</p>`

	subject := fmt.Sprintf("%s電郵地址驗証", config.SystemName)
	content := fmt.Sprintf(contentTemp, code, common.VerificationValidMinutes)

	return stmp.Render(email, subject, content)
}

func SendQuotaWarningCodeEmail(userName, email string, quota int, noMoreQuota bool) error {
	stmp, err := GetSystemStmp()

	if err != nil {
		return err
	}

	contentTemp := `<p style="font-size: 30px">Hi <strong>%s,</strong></p>
		<p>
			%s，當前剩餘額度為 %d，為了不影響您的使用，請及時充值。
		</p>
		
		<p style="text-align: center; font-size: 13px;">
			<a target="__blank" href="%s" class="button" style="color: #ffffff;">點擊充值</a>
		</p>
		
		<p style="color: #858585; padding-top: 15px;">
			如果鏈接無法點擊，請嘗試點擊下面的鏈接或將其複製到瀏覽器中打開<br> %s
		</p>`

	subject := "您的額度即將用盡"
	if noMoreQuota {
		subject = "您的額度已用盡"
	}
	topUpLink := fmt.Sprintf("%s/topup", config.ServerAddress)

	content := fmt.Sprintf(contentTemp, userName, subject, quota, topUpLink, topUpLink)

	return stmp.Render(email, subject, content)
}

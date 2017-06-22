package main

import (
	"fmt"
	"github.com/cloudfoundry-community/gautocloud"
	_ "github.com/cloudfoundry-community/gautocloud/connectors/databases"
	"github.com/cloudfoundry-community/gautocloud/connectors/databases/dbtype"
	_ "github.com/cloudfoundry-community/gautocloud/connectors/smtp"
	"github.com/cloudfoundry-community/gautocloud/connectors/smtp/client"
	"github.com/cloudfoundry-community/gautocloud/connectors/smtp/smtptype"
	"github.com/cloudfoundry-community/gautocloud/loader"
	"github.com/cloudfoundry-incubator/notifications/application"
	"log"
	"os"
)

func main() {
	var dbConfig dbtype.MysqlDatabase
	err := gautocloud.Inject(&dbConfig)
	exitOnError(err)
	os.Setenv("DATABASE_URL", fmt.Sprintf("tcp://%s:%s@%s:%d/%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database))

	var smtpConf smtptype.Smtp
	err = gautocloud.Inject(&smtpConf)
	if _, ok := err.(*loader.ErrGiveService); !ok {
		exitOnError(err)
	}
	loadSmtp(smtpConf)
	if os.Getenv("SMTP_AUTH_MECHANISM") == "" {
		os.Setenv("SMTP_AUTH_MECHANISM", "plain")
	}
	os.Setenv("ROOT_PATH", os.Getenv("HOME"))
	os.Setenv("DEFAULT_UAA_SCOPES", "cloud_controller.read,cloud_controller.write,openid,approvals.me,cloud_controller_service_permissions.read,scim.me,uaa.user,password.write,scim.userids,oauth.approvals")
	env, err := application.NewEnvironment()
	exitOnError(err)
	mother := application.NewMother(env)
	app := application.NewApplication(env, mother)
	defer app.Crash()

	app.Boot()
}
func exitOnError(err error) {
	if err == nil {
		return
	}
	log.Fatal(err.Error())
}
func loadSmtp(smtpConf smtptype.Smtp) error {
	if smtpConf.Host == "" {
		return nil
	}
	os.Setenv("SMTP_HOST", smtpConf.Host)
	os.Setenv("SMTP_PORT", fmt.Sprintf("%d", smtpConf.Port))
	os.Setenv("SMTP_USER", smtpConf.User)
	os.Setenv("SMTP_PASS", smtpConf.Password)
	errorMessage := ""
	_, err := client.SmtpConnector{}.GetSmtp(smtpConf, true, true, true)
	os.Setenv("SMTP_TLS", "true")
	if err != nil {
		errorMessage += "\t- tls with startls: " + err.Error() + "\n"
		_, err = client.SmtpConnector{}.GetSmtp(smtpConf, true, true, false)
	}
	if err != nil {
		errorMessage += "\t- tls: " + err.Error() + "\n"
		os.Setenv("SMTP_TLS", "false")
		_, err = client.SmtpConnector{}.GetSmtp(smtpConf, true, false, false)
	}
	if err != nil {
		errorMessage += "\t- plain auth: " + err.Error() + "\n"
		_, err = client.SmtpConnector{}.GetSmtp(smtpConf, false, false, false)
	}
	if err != nil {
		errorMessage += "\t- no auth: " + err.Error() + "\n"
		return fmt.Errorf("No smtp are reachable (trying: tls with starttls, tls, plain auth and no auth):\n" + errorMessage)
	}
	return nil

}

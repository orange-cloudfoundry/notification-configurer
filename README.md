# Notification-configurer

Wrap [notifications](https://github.com/cloudfoundry-incubator/notifications) app to make it run simply on your cloud foundry.

The bosh-release associated to this app also deploys the app as a cloud foundry app but does not autoconnect services whereas this configurer will do.

## Deploy

1. Clone this repo
2. Create a mysql service on your cloud foundry (e.g.: `cf cs p-mysql 100mb mysql-notification`)
3. Create a smtp service as a sendgrid instance (e.g.: `cf cs sendgrid free smtp-notification`) or by using an user provided service, 
to do this modify the `smtp.json` file and run `cf cups smtp-notification -p smtp.json`
4. Modify the `manifest.yml` file with your own configuration
5. Run `cf push`

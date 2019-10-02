# Create Nice Cards From Github Issues

Create an access token at:  https://github.com/settings/tokens

```shell script
export ACCESS_TOKEN=db015666.
go run ./gen --owner argoproj --repo argo-cd --exclude-label 'wontfix' --exclude-label 'workaround' --exclude-label 'help wanted' > enhancements.html
go run ./gen --owner argoproj --repo argo-cd --label 'bug' --exclude-label 'wontfix' --exclude-label 'workaround' --exclude-label 'help wanted' > bugs.html
```

![cards](docs/images/cards.png)
# Create Nice Cards From Github Issues

Create an access token at:  https://github.com/settings/tokens

```shell script
export ACCESS_TOKEN=db015666.
go run ./gen --owner argoproj --repo argo-cd --label enhancement --label helm > cards.html
```

![cards](docs/images/cards.png)
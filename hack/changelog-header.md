### Linux

```shell
curl -L https://github.com/jenkins-x-quickstarts-labs/step-goreleaser/releases/download/v{{.Version}}/step-goreleaser-linux-amd64.tar.gz | tar xzv
sudo mv step-goreleaser /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x-quickstarts-labs/step-goreleaser/releases/download/v{{.Version}}/step-goreleaser-darwin-amd64.tar.gz | tar xzv
sudo mv step-goreleaser /usr/local/bin
```


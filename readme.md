## Build Image

```bash
tag=$(git rev-parse --short HEAD)
make docker-build docker-push IMG=ishenle/mmchatgpt:$tag
```

## Requirements

- OpenAI Key
- Mattermost instance && mattermost token

## Deploy

```bash
cp config.yaml.example config.yaml
dockr run --name mmgpt -d --rm \
  -v $(pwd)/config.yaml:/config.yaml \
  ishenle/mmchatgpt
```

## References

- https://docs.mattermost.com/developer/personal-access-tokens.html
- https://platform.openai.com/account/api-keys



```
helm install stable/nginx-ingress \
  --set controller.defaultBackendService=true \
  --set defaultBackend.image.repository="fujexo/potz-holzoepfel-und-zipfelchape" \
  --set defaultBackend.port=80 \
  --set controller.service.type=NodePort \
  --set defaultBackend.image.tag=latest
```

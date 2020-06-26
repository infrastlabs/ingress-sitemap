# ingress-sitemap

K8S集群Ingress导航

## Feats

- 可视导航页、hosts自定义域名快速添加
- 自动由svc生成ingress 少一层ingress维护

## QuickStart

- install

```bash
kc apply -k https://gitee.com/infrastlabs/ingress-sitemap//deploy
```

- service
  - ingressed label: `auto-ingress/enabled: enabled`
  - un-ingressed label: `auto-ingress/enabled: disabled`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
  labels:
    app: myapp
    auto-ingress/enabled: 'enabled'
spec:
  type: LoadBalancer
  ports:
  - name: http
    port: 80
    targetPort: http
  selector:
    app: myapp
```

## Dev

```bash
sh go-build.sh

export AUTO_INGRESS_SERVER_NAME=demo1.cn
export GW_HTTPS_PORT=31714
go run *.go
./ingsitemap
```

## TODO

- ~~罗列ingress, 使用goTemplates生成index.html模板 (ref: aigb-swagger)~~ Done.
- ~~自动由svc生成ing: 识别label, +env自定义域名;~~  (1.更好定制domain后缀；2.少一层ingress维护)
- 更改域名清理已有ingress(loopJudge)
- ingress命名分组展示

## ref

- https://gitee.com/infrastlabs/k8s-jumpserver
- https://gitee.com/infrastlabs/dh-pages
- http://git.ali.devcn.fun:81/g-dev2/fk-aigb-swagger
- 
- kubernetes-auto-ingress https://github.com/hxquangnhat/kubernetes-auto-ingress
- kubetop https://github.com/siadat/kubetop

# Yeet

Yeet source code to production.
Just for personal testing, don't use it.

```yaml
# app.yaml
push:
  build: true # tag image at hash
  stages:
    - name: dev1
      groups:
        - us-west1
        - us-east2
tag:
  build: false # re-tag the image version at the hash as semver
  stages:
    - name: stag
      groups:
        us: {}
      wait: 2d
    - name: prod-eu
      groups:
        uk: {}
        nl1: {}
        nl2: {}
        eu: {}
      wait: 1d
    - name: prod-us
      groups:
        us-west1:
          manual: true
        us-west2-clustera: {}
        us-west2-clusterb: {}
        us-north1: {}
        us-north2: {}
```

![Logo](assets/lego_headphone.png)

# What is this
A simple tool. It will check your pulbic ip, based on the ip it will ask or register your workday in a .txt file located in for example ~/.workreg/2024/January.txt:
```text
5 home
8 home
10 work
11 work
15 work
17 home
18 work
22 home
24 work
25 work
```
It will show a dialog or a notification:

![Dialog](img/dialog.png)
![Dialog](img/notification.png)

# Workreg

Config your config.yaml file after the first run.
Set a config like this:
```yaml
home:
  ip: 1.1.1.1
  ask: True
work:
  ip: 8.8.8.8
  ask: False
```

# V1
The first version for Linux is here: https://github.com/zoutepopcorn/workreg/releases/tag/v1.0

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
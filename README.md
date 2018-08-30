# Monitoring

## Install

`go get github.com/bborbe/monitoring/cmd/monitoring-check`

`go get github.com/bborbe/monitoring/cmd/monitoring-cron`

## Check

```
monitoring-check \ 
-logtostderr \
-v=2 \
-config sample_config.xml
```

## Cron

```
monitoring-cron \ 
-logtostderr \
-v=2 \
-config sample_config.xml
```

## Config

```
<?xml version="1.0"?>
<nodes>
  <node check="tcp" host="www.benjamin-borbe.de" port="80">
    <node check="http" url="http://www.benjamin-borbe.de" expectstatuscode="200" expecttitle="Benjamin Borbe"/>
  </node>
</nodes>
```
## Available checks

Silent check

`silent="true"`

Disable check

`disabled="true"`

### TCP

```
<node check="tcp" host="www.benjamin-borbe.de" port="80">
```

### HTTP

```
<node check="http" url="http://www.benjamin-borbe.de"></node>
<node check="http" url="http://www.benjamin-borbe.de" expectstatuscode="200"></node>
<node check="http" url="http://www.benjamin-borbe.de" expectbody="ks.cfg" ></node>
<node check="http" url="http://www.benjamin-borbe.de" expectcontent="ks.cfg" ></node>
<node check="http" url="http://www.benjamin-borbe.de" expecttitle="Benjamin Borbe"></node>
<node check="http" url="http://aptly.benjamin-borbe.de/api/version" username="api" passwordfile="/etc/aptly_api_password"/>
<node check="http" url="http://aptly.benjamin-borbe.de/api/version" username="api" password="secret"/>
```

### Webdriver

```
<node check="webdriver" url="http://sonar.benjamin-borbe.de/">
  <action type="click" strategy="xpath"  query="//a[@id='login-link']"/>
  <action type="fill" strategy="xpath"  query="//input[@id='login']" value="bborbe"/>
  <action type="fill" strategy="xpath"  query="//input[@id='password']" value="test123"/>
  <action type="submit" strategy="xpath"  query="//input[@type='submit']"/>
  <action type="expecttitle" value="SonarQube"/>
  <action type="notexists" strategy="xpath"  query="//a[@class='login-link']"/>
</node>
```

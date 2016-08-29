# Monitoring

## Install

`go get github.com/bborbe/monitoring/bin/monitoring_check`

`go get github.com/bborbe/monitoring/bin/monitoring_cron`

## Check

```
monitoring_check \ 
-logtostderr \
-v=2 \
-config sample_config.xml
```

## Cron

```
monitoring_cron \ 
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

## Continuous integration

[Jenkins](https://www.benjamin-borbe.de/jenkins/job/Go-Monitoring/)

## Copyright and license

    Copyright (c) 2016, Benjamin Borbe <bborbe@rocketnews.de>
    All rights reserved.
    
    Redistribution and use in source and binary forms, with or without
    modification, are permitted provided that the following conditions are
    met:
    
       * Redistributions of source code must retain the above copyright
         notice, this list of conditions and the following disclaimer.
       * Redistributions in binary form must reproduce the above
         copyright notice, this list of conditions and the following
         disclaimer in the documentation and/or other materials provided
         with the distribution.

    THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
    "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
    LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
    A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
    OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
    SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
    LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
    DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
    THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
    (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
    OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

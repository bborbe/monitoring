#!/bin/sh


export DEBFULLNAME="Benjamin Borbe"
export EMAIL=bborbe@rocketnews.de

export DEB_SERVER=misc.rn.benjamin-borbe.de
export TARGET_DIR=opt/monitoring/bin

export NAME=monitoring
export BINS="monitoring_server"
export INSTALLS="github.com/bborbe/monitoring/bin/monitoring_server"
export SOURCEDIRECTORY="github.com/bborbe/monitoring"

export MAJOR=0
export MINOR=1
export BUGFIX=0

# exec
sh src/github.com/bborbe/jenkins/jenkins.sh
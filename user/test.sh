#!/bin/bash
PROJECTPATH=event/booking
go test -coverpkg=$PROJECTPATH/aggregates/...,$PROJECTPATH/commands/...,$PROJECTPATH/event/...,$PROJECTPATH/gokit/...,$PROJECTPATH/repositories/...,$PROJECTPATH/service/... $PROJECTPATH/testing/tests/... -coverprofile=coverage.out

#!/bin/bash
git tag -a v$1 -m "tag $1" && git push origin --tags
#!/usr/bin/env bash
set -e

standard-version --dry-run --skip
read -p "Continue (y/n)? " -n 1 -r
echo ;
if [[ $REPLY =~ ^[Yy]$ ]]; then
	standard-version -s ;
fi

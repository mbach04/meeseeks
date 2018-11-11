#!/bin/bash
#Installs dependencies required by this repository
declare -a deps=("github.com/dgrijalva/jwt-go"
                 "github.com/labstack/echo"
                 "github.com/labstack/echo/middleware"
                 "github.com/spf13/viper"
                 "github.com/kabukky/httpscerts"
		)


for dep in "${deps[@]}"
do
	go get "$dep"
done

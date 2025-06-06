cmake_minimum_required(VERSION 3.10)

project(myproject)

add_executable(myapp
    ./main.cpp
    ./native-host.cpp
    ./wa-app.c
)

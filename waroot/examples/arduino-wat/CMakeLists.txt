cmake_minimum_required(VERSION 4.0)

project(WA_APP)

include_directories(
  ${CMAKE_SOURCE_DIR}
)

add_executable(waapp
  ./native/main.cpp
	./native/arduino-host-debug.cpp
  ./wa-app.c
)

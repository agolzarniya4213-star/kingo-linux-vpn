cmake_minimum_required(VERSION 3.10)
project(kingo-tray)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_AUTOMOC ON)

find_package(Qt5 REQUIRED COMPONENTS Widgets Svg Network)

add_executable(kingo-tray
    main.cpp
    trayicon.cpp
    trayicon.h
    ipcclient.cpp
    ipcclient.h
)

target_link_libraries(kingo-tray
    Qt5::Widgets
    Qt5::Svg
    Qt5::Network
)

install(TARGETS kingo-tray DESTINATION bin)

cmake_minimum_required(VERSION 3.14)
project(kingo-tray VERSION 1.0.0 LANGUAGES CXX)

# Industrial standard compiler flags (NASA/Google level strictness)
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)
add_compile_options(-Wall -Wextra -Wpedantic -Werror)

set(CMAKE_AUTOMOC ON)

# Required components: Widgets (UI), Svg (Icons), Network (IPC Sockets)
find_package(Qt5 REQUIRED COMPONENTS Widgets Svg Network)

# IPC Client Layer
add_library(ipcclient STATIC
    ipcclient.cpp
    ipcclient.h
)
target_link_libraries(ipcclient PUBLIC Qt5::Network)

# Main Tray Application
add_executable(kingo-tray
    main.cpp
    trayicon.cpp
    trayicon.h
)
target_link_libraries(kingo-tray 
    PRIVATE 
    Qt5::Widgets 
    Qt5::Svg 
    ipcclient
)

# Installation rules
install(TARGETS kingo-tray RUNTIME DESTINATION bin)

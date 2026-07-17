cmake_minimum_required(VERSION 3.14)
project(kingo-tray VERSION 1.0.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)
add_compile_options(-Wall -Wextra -Wpedantic -Werror)

set(CMAKE_AUTOMOC ON)

find_package(Qt5 REQUIRED COMPONENTS Widgets Svg Network)

add_library(ipcclient STATIC
    ipcclient.cpp
    ipcclient.h
)
target_link_libraries(ipcclient PUBLIC Qt5::Network)

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

install(TARGETS kingo-tray RUNTIME DESTINATION bin)

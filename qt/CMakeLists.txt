cmake_minimum_required(VERSION 3.14)
project(kingo-tray VERSION 2.0.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)
add_compile_options(-Wall -Wextra -Wpedantic -Werror)

set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)

# Added QuickControls2 to the component list
find_package(Qt5 REQUIRED COMPONENTS Quick Gui Qml Network QuickControls2)

add_library(ipcclient STATIC
    ipcclient.cpp
    ipcclient.h
)
target_link_libraries(ipcclient PUBLIC Qt5::Network)

add_executable(kingo-tray
    main.cpp
    appcore.cpp
    appcore.h
    resources.qrc
)

# Linked Qt5::QuickControls2 to provide QML plugin paths at runtime
target_link_libraries(kingo-tray 
    PRIVATE 
    Qt5::Quick
    Qt5::Gui
    Qt5::Qml
    Qt5::QuickControls2
    ipcclient
)

install(TARGETS kingo-tray RUNTIME DESTINATION bin)

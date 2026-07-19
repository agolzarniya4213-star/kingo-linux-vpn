cmake_minimum_required(VERSION 3.16)
project(kingo-linux-vpn-ui LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_AUTOMOC ON)

find_package(Qt6 REQUIRED COMPONENTS Quick QuickControls2 Gui Network Widgets)

# Process resources.qrc
qt_add_resources(RES_SOURCES "resources.qrc")

qt_add_executable(kingo-linux-vpn-ui
    main.cpp
    ipcclient.h ipcclient.cpp
    vpncontroller.h vpncontroller.cpp
    trayicon.h trayicon.cpp
    ${RES_SOURCES} # FIX: This was missing! Injects QML into binary.
)

target_link_libraries(kingo-linux-vpn-ui PRIVATE Qt6::Quick Qt6::QuickControls2 Qt6::Gui Qt6::Network Qt6::Widgets)

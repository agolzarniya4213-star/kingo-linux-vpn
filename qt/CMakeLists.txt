cmake_minimum_required(VERSION 3.16)
project(kingo-linux-vpn-ui LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_AUTOMOC ON)

find_package(Qt6 COMPONENTS Quick QuickControls2 Gui REQUIRED)

qt_add_executable(kingo-linux-vpn-ui main.cpp)

qt_add_qml_module(kingo-linux-vpn-ui
    URI KingoVPN
    VERSION 1.0
    QML_FILES main.qml
)

target_link_libraries(kingo-linux-vpn-ui PRIVATE Qt6::Quick Qt6::QuickControls2 Qt6::Gui)

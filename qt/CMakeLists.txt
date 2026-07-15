cmake_minimum_required(VERSION 3.21)
project(KingoLinuxVpnUI VERSION 0.1.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(Qt6 REQUIRED COMPONENTS Quick QuickControls2 Network)

qt_standard_project_setup()

qt_add_executable(kingo-linux-vpn-ui
    src/main.cpp
    src/appcontroller.cpp
    src/appcontroller.h
    src/models/serverlistmodel.cpp
    src/models/serverlistmodel.h
)

qt_add_qml_module(kingo-linux-vpn-ui
    URI KingoVpn
    VERSION 1.0
    QML_FILES
        qml/Main.qml
        qml/components/ActionButton.qml
        qml/components/InfoCard.qml
        qml/components/ServerRow.qml
    RESOURCES
        assets/logo.svg
)

target_link_libraries(kingo-linux-vpn-ui
    PRIVATE Qt6::Quick Qt6::QuickControls2 Qt6::Network
)

cmake_minimum_required(VERSION 3.10)
project(kingo-tray)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_AUTOMOC ON)

# پیدا کردن کتابخانه‌های Qt
find_package(Qt5 REQUIRED COMPONENTS Widgets Svg)

# تعریف فایل‌های اجرایی
add_executable(kingo-tray
    main.cpp
    trayicon.cpp
    trayicon.h
)

# لینک کردن کتابخانه‌ها
target_link_libraries(kingo-tray
    Qt5::Widgets
    Qt5::Svg
)

# دستورات نصب (اختیاری برای آینده)
install(TARGETS kingo-tray DESTINATION bin)

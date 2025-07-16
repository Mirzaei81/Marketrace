cmake_minimum_required(VERSION 3.14)
project(MyProject)

# Set C++ standard
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(openssl CONFIG REQUIRED)
# Add your source files
set(SOURCES
    dllmain.cpp
)

set(HEADERS
    pch.h
)

# Create the executable
add_executable(${PROJECT_NAME} ${SOURCES} ${HEADERS})

# Link libraries
target_link_libraries(${PROJECT_NAME}
    winhttp
    libcrypto 
    libssl
)


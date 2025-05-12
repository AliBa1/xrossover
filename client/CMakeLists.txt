cmake_minimum_required(VERSION 3.20)
project(xrossover LANGUAGES CXX)

set(EXPORT_COMPILE_COMMANDS TRUE)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(SDL3 REQUIRED CONFIG REQUIRED COMPONENTS SDL3-shared)

# Add include directory
include_directories(include)

add_executable(xrossover src/main.cpp src/Game.cpp)

# target_include_directories(xrossover PUBLIC ${SDL3_INCLUDE_DIRS})
# target_compile_definitions(xrossover PRIVATE SDL_MAIN_USE_CALLBACKS)

# Link to the actual SDL3 library.
target_link_libraries(xrossover PRIVATE SDL3::SDL3)

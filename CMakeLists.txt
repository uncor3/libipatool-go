cmake_minimum_required(VERSION 3.16)
project(ipatool-go)

# Find Go compiler
find_program(GO_EXECUTABLE go REQUIRED)

# Set Go build flags for dynamic library
set(GO_BUILD_FLAGS 
    -buildmode=c-shared
)

# Define the Go source file
set(GO_SOURCE ${CMAKE_CURRENT_SOURCE_DIR}/main.go)

# Define output files
set(GO_LIB_OUTPUT ${CMAKE_CURRENT_BINARY_DIR}/libipatool-go.so)
set(GO_HEADER_OUTPUT ${CMAKE_CURRENT_BINARY_DIR}/libipatool-go.h)

# Custom command to build Go shared library
add_custom_command(
    OUTPUT ${GO_LIB_OUTPUT} ${GO_HEADER_OUTPUT}
    COMMAND ${GO_EXECUTABLE} build ${GO_BUILD_FLAGS} -o ${GO_LIB_OUTPUT} ${GO_SOURCE}
    DEPENDS ${GO_SOURCE}
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    COMMENT "Building Go shared library"
)

# Create a custom target for the Go library
add_custom_target(go_build ALL
    DEPENDS ${GO_LIB_OUTPUT} ${GO_HEADER_OUTPUT}
)

# Create an imported library target
add_library(ipatool-go SHARED IMPORTED GLOBAL)
set_target_properties(ipatool-go PROPERTIES
    IMPORTED_LOCATION ${GO_LIB_OUTPUT}
    INTERFACE_INCLUDE_DIRECTORIES ${CMAKE_CURRENT_BINARY_DIR}
)

# Make sure the Go library is built before the imported target is used
add_dependencies(ipatool-go go_build)

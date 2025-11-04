cmake_minimum_required(VERSION 3.16)
project(ipatool-go)

# Find Go compiler
find_program(GO_EXECUTABLE go)

# If not found in PATH, check common Windows location as a fallback
if(NOT GO_EXECUTABLE AND WIN32)
    find_program(GO_EXECUTABLE go PATHS "C:/msys64/mingw64/bin" NO_DEFAULT_PATH)
endif()

# If still not found, fail the configuration
if(NOT GO_EXECUTABLE)
    message(FATAL_ERROR "Go compiler (go.exe) not found.")
else()
    message(STATUS "Found Go compiler: ${GO_EXECUTABLE}")
endif()

# Set Go build flags for static library
set(GO_BUILD_FLAGS 
    -buildmode=c-archive
    -ldflags "-s -w"
)

# Define the Go source file
set(GO_SOURCE ${CMAKE_CURRENT_SOURCE_DIR}/main.go)

# Define output files
set(GO_LIB_OUTPUT ${CMAKE_CURRENT_BINARY_DIR}/libipatool-go.a)
set(GO_HEADER_OUTPUT ${CMAKE_CURRENT_BINARY_DIR}/libipatool-go.h)

# Custom command to build Go static library
add_custom_command(
    OUTPUT ${GO_LIB_OUTPUT} ${GO_HEADER_OUTPUT}
    COMMAND ${GO_EXECUTABLE} build ${GO_BUILD_FLAGS} -o ${GO_LIB_OUTPUT} ${GO_SOURCE}
    DEPENDS ${GO_SOURCE}
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    COMMENT "Building Go static library"
)

# Create a custom target for the Go library
add_custom_target(go_build ALL
    DEPENDS ${GO_LIB_OUTPUT} ${GO_HEADER_OUTPUT}
)

# Create an imported library target
add_library(ipatool-go STATIC IMPORTED GLOBAL)
set_target_properties(ipatool-go PROPERTIES
    IMPORTED_LOCATION ${GO_LIB_OUTPUT}
    INTERFACE_INCLUDE_DIRECTORIES ${CMAKE_CURRENT_BINARY_DIR}
)
if(APPLE)
    target_link_libraries(ipatool-go INTERFACE
        "-framework Security"
    )
endif()

# Make sure the Go library is built before the imported target is used
add_dependencies(ipatool-go go_build)
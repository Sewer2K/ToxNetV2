# Simplified Dependencies.cmake for cross-compilation

# Threads
find_package(Threads REQUIRED)

# libsodium - we provide the path
if(NOT SODIUM_LIBRARY OR NOT SODIUM_INCLUDE_DIR)
    find_library(SODIUM_LIBRARY NAMES sodium)
    find_path(SODIUM_INCLUDE_DIR NAMES sodium.h)
endif()

# Set variables
set(CMAKE_THREAD_LIBS_INIT ${CMAKE_THREAD_LIBS_INIT})
set(CORE_PKGS "")
set(AV_PKGS "")

# No pkgconfig checks
set(PKG_CONFIG_FOUND FALSE)
set(PKG_CONFIG_EXECUTABLE "")

@REM set ANDROID_HOME=C:/Sdk/android-sdk
set ANDROID_HOME=D:/sdk
@REM set ANDROID_NDK_HOME=C:/Sdk/android-sdk/ndk/25.1.8937393
set ANDROID_NDK_HOME=D:/sdk/ndk/ndk-bundle
gomobile bind -ldflags "-s -w" -v -androidapi 21 "gostc-sub/lib"
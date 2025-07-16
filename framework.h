#pragma once

#define WIN32_LEAN_AND_MEAN             // Exclude rarely-used stuff from Windows headers
// Windows Header Files
#include <windows.h>

typedef void (__stdcall *PFN_CALLBACK1)(const LPCWSTR str);
extern "C" __declspec(dllexport)   wchar_t InputToken[512];
extern "C" __declspec(dllexport) LRESULT DisplayMyMessage(HWND hwndOwner);
extern "C" __declspec(dllexport) LRESULT OnSuccessCallBack(PFN_CALLBACK1 callback1);


// dllmain.cpp : Defines the entry point for the DLL application.
#include "pch.h"
#include <openssl/ssl.h>
#include <winhttp.h>
#include <sstream>

#include <strsafe.h>
#include <vector>

// I made this struct to more conveniently store the
// positions / size of each window in the dialog
typedef struct SizeAndPos_s
{
    int x, y, width, height;
} SizeAndPos_t;


void DisplayErrorBox(const wchar_t* lpszFunction);
HINSTANCE  hInstance = 0;
// Typically these would be #defines, but there
// is no reason to not make them constants
const WORD ID_labelInfo = 1;
const WORD ID_btnQUIT = 2;
const WORD ID_CheckBox = 3;
const WORD ID_txtEdit = 4;
const WORD ID_btnShow = 5;

const SizeAndPos_t mainWindow   =   { 150,    150,    310,    510 };

const SizeAndPos_t labelInfo      =   { 200,    10,    280,    30 };
const SizeAndPos_t txtEdit      =   { 10,     50,    280,    375 };
const SizeAndPos_t btnQuit      =   { 10,    425,     80,     35 };

HWND txtEditHandle = NULL;
PFN_CALLBACK1 callback = NULL;
EVP_PKEY* pKey = NULL;
// hwnd:    All window processes are passed the handle of the window
//         that they belong to in hwnd.
// msg:     Current message (e.g., WM_*) from the OS.
// wParam:  First message parameter, note that these are more or less
//          integers, but they are really just "data chunks" that
//          you are expected to memcpy as raw data to float, etc.
// lParam:  Second message parameter, same deal as above.
LRESULT CALLBACK WndProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam)
{

	switch (msg)
	{

	case WM_CREATE:
		// Note that the "parent window" is the dialog itself. Since we are
		// in the dialog's WndProc, the dialog's handle is passed into hwnd.
		//
		//CreateWindow( lpClassName,        lpWindowName,       dwStyle,                                x,          y,          nWidth,         nHeight,            hWndParent,     hMenu,              hInstance,      lpParam
		//CreateWindow( windowClassName,    initial text,       style (flags),                          xPos,       yPos,       width,          height,             parentHandle,   menuHandle,         instanceHandle, param);
		CreateWindow(TEXT("Static"), TEXT("توکن احراز هویت : "), WS_VISIBLE | WS_CHILD | WS_TABSTOP, labelInfo.x, labelInfo.y, labelInfo.width, labelInfo.height, hwnd, (HMENU)ID_labelInfo, NULL, NULL);
		txtEditHandle = CreateWindow(TEXT("Edit"), TEXT(""), WS_CHILD | WS_VISIBLE | WS_BORDER | ES_MULTILINE | WS_VSCROLL, txtEdit.x, txtEdit.y, txtEdit.width, txtEdit.height, hwnd, (HMENU)ID_txtEdit, NULL, NULL);
		CreateWindow(TEXT("Button"), TEXT("کنسل"), WS_VISIBLE | WS_CHILD, btnQuit.x, btnQuit.y, btnQuit.width, btnQuit.height, hwnd, (HMENU)ID_btnQUIT, NULL, NULL);
		CreateWindow(TEXT("Button"), TEXT("ثبت"), WS_VISIBLE | WS_CHILD, 200, btnQuit.y, btnQuit.width, btnQuit.height, hwnd, (HMENU)ID_btnShow, NULL, NULL);

		break;

		// For more information about WM_COMMAND, see
		// https://msdn.microsoft.com/en-us/library/windows/desktop/ms647591(v=vs.85).aspx
	case WM_COMMAND:

		// The LOWORD of wParam identifies which control sent
		// the WM_COMMAND message. The WM_COMMAND message is
		// sent when the button has been clicked.
		if (LOWORD(wParam) == ID_btnQUIT)
		{
			PostQuitMessage(0);
		}
		else if (LOWORD(wParam) == ID_btnShow)
		{
			size_t sPlainSize = 0;
			int textLength_WithNUL = GetWindowTextLength(txtEditHandle) + 1;
			// WARNING: If you are compiling this for C, please remember to remove the (TCHAR*) cast.
			LPSTR textBoxText = (LPSTR)LocalAlloc(0, textLength_WithNUL);


			if (!GetWindowTextA(txtEditHandle, textBoxText, textLength_WithNUL)) {
				DisplayErrorBox(L"خطا هنگام پرداز متن داخل باکس");
				EVP_PKEY_free(pKey);
				return FALSE;
			}
			EVP_PKEY_CTX* ctx = EVP_PKEY_CTX_new(pKey, nullptr);
			if (!ctx) {
				DisplayErrorBox(L"خطا هنگام پرداز متن داخل باکس context");
				EVP_PKEY_free(pKey);
				return FALSE;
			}
			if (EVP_PKEY_decrypt_init(ctx) <= 0) {
				DisplayErrorBox(L"خطا هنگام پرداز متن داخل باکس init");
				EVP_PKEY_free(pKey);
				EVP_PKEY_CTX_free(ctx);
				return FALSE;
			}
			if (EVP_PKEY_decrypt(ctx, nullptr, &sPlainSize, reinterpret_cast<PUCHAR>(textBoxText), textLength_WithNUL) <= 0) {
				DisplayErrorBox(L"خطا هنگام پرداز متن داخل باکس decrpyt");
				EVP_PKEY_free(pKey);
				EVP_PKEY_CTX_free(ctx);
				return FALSE;
			}

			if (sPlainSize <= 0) {
				DisplayErrorBox(L"خطا هنگام پردازش private_key!");
				return FALSE;
			}
			PUCHAR plainText = (PUCHAR)LocalAlloc(0, sPlainSize);
			if (EVP_PKEY_decrypt(ctx, plainText, &sPlainSize, reinterpret_cast<PUCHAR>(textBoxText), textLength_WithNUL) <= 0) {
				DisplayErrorBox(L"خطا هنگام پرداز متن داخل باکس decrpyt");
				EVP_PKEY_free(pKey);
				EVP_PKEY_CTX_free(ctx);
				return FALSE;
			}
			HINTERNET hIntOpen = WinHttpOpen(L"WinHTTP Example/1.0", WINHTTP_ACCESS_TYPE_DEFAULT_PROXY, WINHTTP_NO_PROXY_NAME, WINHTTP_NO_PROXY_BYPASS, 0);
			if (!hIntOpen) {
				DisplayErrorBox(L"خطا هنگام  اتصال به اینترنت.");
				return FALSE;
			}

			// Connect to server
			HINTERNET hConnect = WinHttpConnect(hIntOpen, L"https://seller.digikala.com", INTERNET_DEFAULT_HTTPS_PORT, 0);
			if (!hConnect) {
				DisplayErrorBox(L"خطا هنگام  اتصال به اینترنت.");
				WinHttpCloseHandle(hIntOpen);
				EVP_PKEY_free(pKey);
				EVP_PKEY_CTX_free(ctx);
				return FALSE;
			}
			    // Create request handle
			HINTERNET hRequest = WinHttpOpenRequest(hConnect, L"POST", L"/open-api/v1/auth/token",
				NULL, WINHTTP_NO_REFERER,
				WINHTTP_DEFAULT_ACCEPT_TYPES,
				WINHTTP_FLAG_SECURE);
			if (!hRequest) {

				DisplayErrorBox(L"خطا هنگام  اتصال به اینترنت.");
				WinHttpCloseHandle(hIntOpen);
				WinHttpCloseHandle(hConnect);
				EVP_PKEY_free(pKey);
				EVP_PKEY_CTX_free(ctx);
				return FALSE;
			}
			PCSTR rgpszAcceptTypes[] = { "application/json", NULL };
			LPCWSTR headers = L"Content-Type: application/json\r\n";
			if (!WinHttpAddRequestHeaders(hRequest, headers, -1l, WINHTTP_ADDREQ_FLAG_ADD)) {
				DisplayErrorBox(L"خطا هنگام  اتصال به اینترنت headers.");
			}

			PCSTR cstrHeader = "Content-Type: application/json";
			rsize_t dataSize = sizeof("{authorization_code: '' }" + sPlainSize);
			PCHAR rawData = (PCHAR)LocalAlloc(0, dataSize);
			int iBuffSize = sprintf_s(rawData, dataSize, "{authorization_code: '%s' }", reinterpret_cast<PCSTR>(plainText));
			if (!WinHttpSendRequest(hRequest,
				WINHTTP_NO_ADDITIONAL_HEADERS, 0,
				rawData, dataSize, dataSize, 0)) {
				DisplayErrorBox(L"خطا هنگام  اتصال به اینترنت SendRequest.");
				WinHttpCloseHandle(hIntOpen);
				WinHttpCloseHandle(hConnect);
				EVP_PKEY_free(pKey);
				EVP_PKEY_CTX_free(ctx);
				return FALSE;
			}
			if (!WinHttpReceiveResponse(hRequest, NULL)) {
				DisplayErrorBox(L"خطا هنگام  اتصال به اینترنت ReciveResponse.");
				WinHttpCloseHandle(hIntOpen);
				WinHttpCloseHandle(hConnect);
				EVP_PKEY_free(pKey);
				EVP_PKEY_CTX_free(ctx);
				return FALSE;
			}
			DWORD bytesAvailable = 0;
			std::wstringstream  wss; 
			do {
				WinHttpQueryDataAvailable(hRequest, &bytesAvailable);
				if (bytesAvailable > 0) {
					std::vector<PWCHAR> buffer(bytesAvailable + 1, 0);
					DWORD bytesRead = 0;
					WinHttpReadData(hRequest, buffer.data(), bytesAvailable, &bytesRead);
					wss<<buffer.data();
				}
			} while (bytesAvailable > 0);
			if (callback) {
				callback(wss.str().c_str());
			}
			MessageBoxW(NULL, L"موفقیت!", L"توکن با موفقیت ثبت شد!", 0);
			LocalFree(textBoxText);
			LocalFree(plainText);
			EVP_cleanup();
			PostQuitMessage(0);
			break;
		}

	case WM_DESTROY:
		PostQuitMessage(0);
		break;
		}

	return DefWindowProc(hwnd, msg, wParam, lParam);
}


// hInstance: This handle refers to the running executable
// hPrevInstance: Not used. See https://blogs.msdn.microsoft.com/oldnewthing/20040615-00/?p=38873
// lpCmdLine: Command line arguments.
// nCmdShow: a flag that says whether the main application window
//           will be minimized, maximized, or shown normally.
//
// Note that it's necessary to use _tWinMain to make it
// so that command line arguments will work, both
// with and without UNICODE / _UNICODE defined.

BOOL WINAPI DllMain(
    HINSTANCE hinstDLL,  // handle to DLL module
    DWORD fdwReason,     // reason for calling function
    LPVOID lpvReserved   // reserved
){
	hInstance = hinstDLL;
	FILE* pFile = NULL;
    errno_t err =  fopen_s(&pFile,"private_key.pem","r");
    OpenSSL_add_all_algorithms();
	OPENSSL_init_crypto(OPENSSL_INIT_LOAD_CONFIG, nullptr);

    if (pFile == NULL) {
        DisplayErrorBox(L"فایل private_key.pem موجود نیست!");
        return FALSE;
    }
	if (PEM_read_PrivateKey == NULL) {
        DisplayErrorBox(L"خطا nv ssl!");
        return FALSE;
	}
	pKey = (*PEM_read_PrivateKey)(pFile, &pKey, NULL, NULL);
    if (!pKey) {
        DisplayErrorBox(L"خطا هنگام پردازش private_key!");
        return FALSE;
    }
    fclose(pFile);

    return TRUE;
}
extern "C" __declspec(dllexport) LRESULT OnSuccessCallBack(PFN_CALLBACK1 callback1) {
    callback = callback1;
	return 0;
}
extern "C" __declspec(dllexport) LRESULT DisplayMyMessage(HWND hwndOwner) {
    WNDCLASS mainWindowClass = { 0 };

    // You can set the main window name to anything, but
    // typically you should prefix custom window classes
    // with something that makes it unique.
    mainWindowClass.lpszClassName = TEXT("JRH.MainWindow");

    mainWindowClass.hInstance = hInstance;
    mainWindowClass.hbrBackground = GetSysColorBrush(COLOR_3DFACE);
    mainWindowClass.lpfnWndProc = WndProc;
    mainWindowClass.hCursor = LoadCursor(0, IDC_ARROW);

    RegisterClass(&mainWindowClass);
    //CreateWindow( windowClassName,                windowName,             style,                            xPos,         yPos,       width,              height,            parentHandle,   menuHandle,  instanceHandle, param);
    CreateWindow(   mainWindowClass.lpszClassName,  TEXT("Main Window"),    WS_OVERLAPPEDWINDOW | WS_VISIBLE, mainWindow.x, mainWindow.y, mainWindow.width, mainWindow.height, NULL,           0,           hInstance
        ,      NULL);
    MSG  msg;
    while (GetMessage(&msg, NULL, 0, 0))
    {
        TranslateMessage(&msg);
        DispatchMessage(&msg);
    }

    return (int)msg.wParam;
}

void DisplayErrorBox(LPCTSTR lpszFunction) 
{ 
    // Retrieve the system error message for the last-error code

    LPVOID lpMsgBuf;
    LPVOID lpDisplayBuf;
    DWORD dw = GetLastError(); 

    FormatMessage(
        FORMAT_MESSAGE_ALLOCATE_BUFFER | 
        FORMAT_MESSAGE_FROM_SYSTEM |
        FORMAT_MESSAGE_IGNORE_INSERTS,
        NULL,
        dw,
        MAKELANGID(LANG_NEUTRAL, SUBLANG_DEFAULT),
        (LPTSTR) &lpMsgBuf,
        0, NULL );

    // Display the error message and clean up

    lpDisplayBuf = (LPVOID)LocalAlloc(LMEM_ZEROINIT, 
        (lstrlen((LPCTSTR)lpMsgBuf)+lstrlen((LPCTSTR)lpszFunction)+40)*sizeof(TCHAR)); 
    StringCchPrintf((LPTSTR)lpDisplayBuf, 
        LocalSize(lpDisplayBuf) / sizeof(TCHAR),
        TEXT("%s failed with error %d: %s"), 
        lpszFunction, dw, lpMsgBuf); 
    MessageBox(NULL, (LPCTSTR)lpDisplayBuf, TEXT("Error"), MB_OK); 

    LocalFree(lpMsgBuf);
    LocalFree(lpDisplayBuf);
}

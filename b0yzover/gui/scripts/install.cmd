@echo off
chcp 65001 >nul 2>&1
title b0yzover - Kurulum Sihirbazi

rem ============================================================
rem  b0yzover - Modern Kurulum Sihirbazi (Windows)
rem  Bagimliliklari otomatik algilar ve hizlica kurar.
rem ============================================================

rem --- ANSI renk kurulumu (Windows 10+) - setlocal'den ONCE tanimlanir,
rem --- boylece alt rutinlerdeki setlocal/endlocal renkleri silmez.
for /f "delims=" %%A in ('echo prompt $E^| cmd') do set "ESC=%%A"
set "R=%ESC%[91m"
set "G=%ESC%[92m"
set "Y=%ESC%[93m"
set "C=%ESC%[96m"
set "W=%ESC%[97m"
set "GR=%ESC%[90m"
set "BD=%ESC%[1m"
set "N=%ESC%[0m"

set "APP=b0yzover"
set "FAILED=0"
set /a STEPS_TOTAL=4
set /a STEP_DONE=0
set "ST1=wait"
set "ST2=wait"
set "ST3=wait"
set "ST4=wait"
set "PKG=algilaniyor..."
set "GOVER=kontrol ediliyor..."
set "PLAT=Windows"

pushd "%~dp0.." >nul 2>&1
set "PROJECT_DIR=%CD%"
popd >nul 2>&1

setlocal EnableExtensions EnableDelayedExpansion

call :nowcs
set "TSTART=%WCS%"

cls

rem ============================ ADIM 1: ortam ============================
set "ST1=run"
call :repaint
call :nowcs
set "TS=%WCS%"
call :detect_env
call :fin 1 0 !TS!

rem ============================ ADIM 2: Go ===============================
set "ST2=run"
call :repaint
call :nowcs
set "TS=%WCS%"
where go >nul 2>&1
if errorlevel 1 (
  call :install_go
  if errorlevel 1 (
    call :fin 2 1 !TS!
    goto summary
  )
)
for /f "delims=" %%V in ('go version 2^>nul') do set "GOVER=%%V"
call :fin 2 0 !TS!

rem ==================== ADIM 3: bagimliliklar ============================
set "ST3=run"
call :repaint
call :nowcs
set "TS=%WCS%"
pushd "%PROJECT_DIR%"
set "RUN_CMD=go mod download"
call :runspin "go mod download"
set "RC=%SPINRC%"
popd
call :fin 3 !RC! !TS!
if not "!RC!"=="0" (
  call :showlog
  goto summary
)

rem ======================= ADIM 4: derleme ===============================
set "ST4=run"
call :repaint
call :nowcs
set "TS=%WCS%"
pushd "%PROJECT_DIR%"
set "RUN_CMD=go build -o b0yzover.exe ."
call :runspin "go build"
set "RC=%SPINRC%"
popd
call :fin 4 !RC! !TS!
if not "!RC!"=="0" (
  call :showlog
  goto summary
)

:summary
call :nowcs
set /a TOT=WCS-TSTART
set /a SE=TOT/10, DE=TOT%%10
set "TOTALT=!SE!.!DE!s"
call :repaint
call :summarybox
echo.
pause
endlocal & exit /b %FAILED%

rem ------------------------------------------------------------ yardimci

:nowcs
set "T=%TIME%"
if "%T:~0,1%"==" " set "T=0%T:~1%"
set /a HH=1%T:~0,2%-100, MM=1%T:~3,2%-100, SS=1%T:~6,2%-100, CC=1%T:~9,2%-100
set /a WCS=((HH*60+MM)*60+SS)*100+CC
exit /b 0

:fin
rem %1=adim no  %2=rc  %3=baslangic cs
call :nowcs
set /a EL=%WCS%-%3
set /a SE=EL/10, DE=EL%%10
set "TT=%SE%.%DE%s"
if "%~2"=="0" (
  set "ST%~1=ok"
) else (
  set "ST%~1=fail"
  set "FAILED=1"
)
set "T%~1=%TT%"
set /a STEP_DONE+=1
exit /b 0

:detect_env
set "HAVE_WINGET="
set "HAVE_CHOCO="
set "HAVE_SCOOP="
where winget >nul 2>&1 && set "HAVE_WINGET=1"
where choco  >nul 2>&1 && set "HAVE_CHOCO=1"
where scoop  >nul 2>&1 && set "HAVE_SCOOP=1"
set "PKGLIST="
if defined HAVE_WINGET set "PKGLIST=%PKGLIST% winget"
if defined HAVE_CHOCO  set "PKGLIST=%PKGLIST% choco"
if defined HAVE_SCOOP  set "PKGLIST=%PKGLIST% scoop"
if defined PKGLIST (set "PKG=%PKGLIST:~1%") else set "PKG=kurucu bulunamadi"
where go >nul 2>&1 || (set "GOVER=kurulu degil" & exit /b 0)
set "GOVER="
for /f "delims=" %%V in ('go version 2^>nul') do set "GOVER=%%V"
if not defined GOVER set "GOVER=kurulu degil"
exit /b 0

:refreshpath
if exist "%USERPROFILE%\go\bin" set "PATH=%USERPROFILE%\go\bin;%PATH%"
if exist "C:\Program Files\Go\bin" set "PATH=C:\Program Files\Go\bin;%PATH%"
if exist "%LOCALAPPDATA%\Programs\Go\bin" set "PATH=%LOCALAPPDATA%\Programs\Go\bin;%PATH%"
if exist "%USERPROFILE%\scoop\apps\go\current\bin" set "PATH=%USERPROFILE%\scoop\apps\go\current\bin;%PATH%"
exit /b 0

:install_go
call :try_winget && exit /b 0
call :try_choco  && exit /b 0
call :try_scoop  && exit /b 0
set "GOVER=otomatik kurulamadi"
exit /b 1

:try_winget
if not "%HAVE_WINGET%"=="1" exit /b 1
set "RUN_CMD=winget install -e --id GoLang.Go --accept-source-agreements --accept-package-agreements"
call :runspin "Go kuruluyor  ·  winget"
if not "%SPINRC%"=="0" exit /b 1
call :refreshpath
where go >nul 2>&1 || exit /b 1
exit /b 0

:try_choco
if not "%HAVE_CHOCO%"=="1" exit /b 1
set "RUN_CMD=choco install golang -y"
call :runspin "Go kuruluyor  ·  chocolatey"
if not "%SPINRC%"=="0" exit /b 1
call :refreshpath
where go >nul 2>&1 || exit /b 1
exit /b 0

:try_scoop
if not "%HAVE_SCOOP%"=="1" exit /b 1
set "RUN_CMD=scoop install go"
call :runspin "Go kuruluyor  ·  scoop"
if not "%SPINRC%"=="0" exit /b 1
call :refreshpath
where go >nul 2>&1 || exit /b 1
exit /b 0

:runspin
rem %1=spinner basligi ; komut global RUN_CMD icinde
set "RUNTITLE=%~1"
set "SPINFLAG=%TEMP%\%APP%_%RANDOM%%RANDOM%.flag"
set "SPINLOG=%TEMP%\%APP%_last.log"
set "SPINBAT=%TEMP%\%APP%_task.cmd"
> "%SPINBAT%" (
  echo @echo off
  echo %RUN_CMD% ^> "%SPINLOG%" 2^>^&1
  echo (echo %%errorlevel%%) ^> "%SPINFLAG%"
)
if exist "%SPINFLAG%" del "%SPINFLAG%" >nul 2>&1
start "" /b "%SPINBAT%"
set /a SP_I=0
:spinloop
if exist "%SPINFLAG%" goto spindone
set /a SP_I+=1
set /a SP_M=SP_I %% 10
if %SP_M%==0 set "SPC=⠋"
if %SP_M%==1 set "SPC=⠙"
if %SP_M%==2 set "SPC=⠹"
if %SP_M%==3 set "SPC=⠸"
if %SP_M%==4 set "SPC=⠼"
if %SP_M%==5 set "SPC=⠴"
if %SP_M%==6 set "SPC=⠦"
if %SP_M%==7 set "SPC=⠧"
if %SP_M%==8 set "SPC=⠇"
if %SP_M%==9 set "SPC=⠏"
<nul set /p "=  %C%%SPC%%N%  %BD%%RUNTITLE%%N%  %GR%·  %SP_I%s%N%%ESC%[K"
ping -n 2 127.0.0.1 >nul
goto spinloop
:spindone
set /p SPINRC=<"%SPINFLAG%"
<nul set /p "=%ESC%[2K%ESC%[1G"
rem echo komutunun ekledigi sondaki boslugu temizle
set "SPINRC=%SPINRC: =%"
if not defined SPINRC set "SPINRC=1"
del "%SPINFLAG%" >nul 2>&1
exit /b 0

:repaint
cls
echo.
echo  %BD%%W%◆ %APP% — KURULUM%N%  %GR%·  Alt Alan Adi Ele Gecirme Tarayici · %PLAT%%N%
echo.
echo  %GR%┌─ ORTAM ─────────────────────────────────────────────────────%N%
echo  %GR%│%N%  %C%Dizin%N%   %W%%PROJECT_DIR%%N%
echo  %GR%│%N%  %C%Paket%N%   %W%%PKG%%N%
echo  %GR%│%N%  %C%Go%N%     %W%%GOVER%%N%
echo  %GR%└─────────────────────────────────────────────────────────────%N%
echo.
call :progressbar
echo.
call :step 1 "Ortam kontrolu"
call :step 2 "Go arac zinciri"
call :step 3 "Bagimliliklar  ·  go mod download"
call :step 4 "Derleme  ·  go build"
echo.
exit /b 0

:progressbar
set /a FIL=STEP_DONE*30/STEPS_TOTAL
set /a EMP=30-FIL
set /a PCT=STEP_DONE*100/STEPS_TOTAL
set "BAR="
for /l %%I in (1,1,%FIL%) do call set "BAR=%%BAR%%█"
for /l %%I in (1,1,%EMP%) do call set "BAR=%%BAR%%░"
echo  %G%%BAR%%N% %W%%PCT%%%N%
exit /b 0

:step
rem %1=adim no  %2=etiket
set "IDX=%~1"
set "LBL=%~2"
call set "ST=%%ST%IDX%%%"
call set "TT=%%T%IDX%%%"
set "ICON=%GR%○%N%"
set "TAIL="
if "%ST%"=="run"  (set "ICON=%C%●%N%" & set "TAIL=%C%  calisiyor...%N%")
if "%ST%"=="ok"   (set "ICON=%G%✔%N%" & set "TAIL=  %G%%TT%%N%")
if "%ST%"=="fail" (set "ICON=%R%✖%N%" & set "TAIL=  %R%%TT%%N%")
echo    %ICON%  %BD%%LBL%%N%%TAIL%
exit /b 0

:showlog
echo  %R%┌─ HATA GUNLUGU ───────────────────────────────────────────────%N%
echo  %R%│%N%
type "%TEMP%\%APP%_last.log" 2>nul
echo  %R%│%N%
echo  %R%└──────────────────────────────────────────────────────────────%N%
echo.
exit /b 0

:summarybox
if "%FAILED%"=="0" goto summ_ok
echo  %R%┌─ SONUC ──────────────────────────────────────────────────────%N%
echo  %R%│%N%  %BD%%R%✖ KURULUM BASARISIZ%N%  %GR%(%TOTALT%)%N%
echo  %R%│%N%  %Y%Manuel kurulum: https://go.dev/dl/%N%
echo  %R%│%N%  %Y%Go'yu kurduktan sonra bu scripti tekrar calistirin.%N%
echo  %R%└──────────────────────────────────────────────────────────────%N%
exit /b 0
:summ_ok
echo  %G%┌─ SONUC ──────────────────────────────────────────────────────%N%
echo  %G%│%N%  %BD%%G%✔ KURULUM TAMAMLANDI%N%  %GR%(%TOTALT%)%N%
echo  %G%│%N%  %C%▸%N% %W%%PROJECT_DIR%\b0yzover.exe%N%
echo  %G%│%N%  %C%▸%N% b0yzover.exe run --target hedef.com
echo  %G%└──────────────────────────────────────────────────────────────%N%
exit /b 0

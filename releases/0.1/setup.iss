; setup.iss — Inno Setup script for Jit

[Setup]
AppName=Jit
AppVersion=0.1.0
AppPublisher=Plycedes
AppPublisherURL=https://github.com/plycedes/jit
DefaultDirName={pf}\Jit
DefaultGroupName=Jit
UninstallDisplayIcon={app}\jit.exe
OutputBaseFilename=JitInstaller
Compression=lzma
SolidCompression=yes
WizardStyle=modern
PrivilegesRequired=admin

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
Source: "jit.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "README.txt"; DestDir: "{app}"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "LICENSE.txt"; DestDir: "{app}"; Flags: ignoreversion recursesubdirs createallsubdirs

[Icons]
Name: "{group}\Jit CLI"; Filename: "{app}\jit.exe"
Name: "{commondesktop}\Jit CLI"; Filename: "{app}\jit.exe"; Tasks: desktopicon

[Tasks]
Name: "desktopicon"; Description: "Create a desktop icon"; GroupDescription: "Additional icons:"; Flags: unchecked

[Registry]
; Add the Jit install path to the system PATH variable
Root: HKLM; Subkey: "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"; \
ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}"; Flags: preservestringtype

[UninstallDelete]
; Remove directory on uninstall
Type: filesandordirs; Name: "{app}"

[Run]
Filename: "{app}\jit.exe"; Description: "Run Jit now"; Flags: nowait postinstall skipifsilent

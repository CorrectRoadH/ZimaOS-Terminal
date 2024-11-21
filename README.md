# What is this?
This is a ZimaOS Module Application.

# How to use?
Run the following command in ZimaOS(>1.2.5)
```
zpkg install zimaos-terminal
```

# How to Build a ZimaOS Module Application?
ref: https://github.com/CorrectRoadH/ZimaOS-Terminal
## Module Structure Components
### Module Description Files
`raw/usr/lib/extension-release.d/extension-release.zimaos_terminal`
The `zimaos_terminal` in the filename is the final module name, which should match the packaged `zimaos_terminal.raw`
Its content should always be:
```
ID=_any
```
This is used by systemd-sysext to identify the module

`raw/usr/share/casaos/modules/zimaos_terminal.json`
```
{
    "name": "zimaos_terminal", // module name
    "ui": {
        "name": "zimaos_terminal", // should match the frontend folder name under raw/usr/share/casaos/www/modules
        "title": {
          "en_us": "Terminal" // module title
        },
        "prefetch": true,
        "show": true,
        "entry": "/modules/zimaos_terminal/index.html", // frontend entry
        "icon": "/modules/zimaos_terminal/appicon.ico", // frontend icon
        "description": "", 
        "formality": {
          "type": "newtab",
          "props": {
            "width": "100vh",
            "height": "100vh",
            "hasModalCard": true,
            "animation": "zoom-in"
          }
        }
    },
    "services": [
        {
            "name": "zimaos-terminal" // backend systemd service name
        }
    ]     
}
```

## Backend Components
1. Binary Files

Place executable binary files in the `raw/usr/bin` directory

2. Systemd Service File (Optional)

If your application needs to start on boot, you need to create a systemd service file and place it in the `raw/etc/systemd/system` directory. See `raw/usr/lib/systemd/system/zimaos-terminal.service` for reference

## Frontend
Frontend files should be automatically compiled to the `raw/usr/share/casaos/www/modules/` directory

## Packaging
After all the above work is ready, execute `mksquashfs raw/ zimaos_terminal.raw` to generate the final module package

## Best Practices
### Automatic Release
Configure your GitHub Actions workflow to publish your module, see https://github.com/CorrectRoadH/ZimaOS-Terminal/blob/main/.github/workflows/release-raw.yml

You need configure the Action secrets to include your GitHub token, which can be obtained from your GitHub account settings.

Your module package should be published to GitHub releases so that it can be installed via zpkg

### Publishing to App Store
Add your module application to `https://github.com/IceWhaleTech/Mod-Store/blob/main/mod-v2.json` by submitting a PR.

Format:
```
{
    "name": "zimaos_terminal",  // should match the final raw name
    "title": "Zimaos Terminal", // title
    "repo": "CorrectRoadH/ZimaOS-Terminal" // repository address
}
```

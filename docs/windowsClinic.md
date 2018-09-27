# Running full clinic deployment on Windows

### Requirements

Setup was tested on Windows 10 Home Edition.

## 1. Install Docker Toolbox for Windows.

*   Unless you have them already installed choose to install Virtual Box and Git for Windows.
*   Run Docker Toolbox QuickStart Shell to create and setup dockem machine.

## 2. Checkout IRYO WWM repo.

Before checkout, while configuring git you should disable auto conversion to CRLF line endings.

```
git config --global core.autocrlf false
```

## 3. Add IRYO WWM dir to shared folders for docker machine VM.

Run in powershell (as admin) while being in IRYO WWM dir:

```
docker-machine stop
vboxmanage sharedfolder add default --name "iryo" --hostpath "${PWD}" --automount
docker-machine start
```

## 4. Generate certificates and import root certificate to Windows root truststore.

*   Enter `docs/windowsClinic` directory in the admin-mode powershell.
*   Run `generateAndImportCerts.ps1` script:

```
powershell -ExecutionPolicy ByPass -File .\generateAndImportCerts.ps1
```

## 5. Setup location, clinic on cloud and import certificates.

1.  Setup location and clinic on cloud deployment that you intend to connect clinic to. Write down location ID and clinic ID.
2.  Configure certificates for authSync, storageSync and batchStorageSync generated in previous step so they will be accepted as valid by `cloudAuth`.

## 6. Set configuration values.

1.  Edit `frontendConfig.json` to include correct `clinicId` and `locationId`.
2.  Set environment variables. You can do it by running following commands in powershell (replace placeholders with correct values!).

```
$Env:IRYO_TAG = "v0.4.2"
$Env:CLINIC_ID = "<CLINIC_ID>"
$Env:LOCATION_ID = "<LOCATION_ID>"
$Env:CLOUD_AUTH_HOST = "<CLOUD_AUTH_HOST>"
$Env:CLOUD_STORAGE_HOST = "<CLOUD_STORAGE_HOST>"
$Env:AUTH_STORAGE_ENCRYPTION_KEY = "<AUTH_STORAGE_ENCRYPTION_KEY>"
$Env:SYMMETRIC_REGISTRATION_URL = "<CLOUD_SYMMETRIC_REGISTRATION_URL>"
```

`AUTH_STORAGE_ENCRYPTION_KEY` needs to be the same as at your chosen cloud deployment's `cloudAuth`. Otherwise `locatAuth` won't be able to decrypt received auth DB file.

## 7. Start clinic

While being in `docs/windowsClinic` run in powershell:

```
docker-compose up -d
```

Now you should be able to access clinic web interface at `https://iryo.local`.

## Known issues

1.  `locaNats` often does not start on the first try. When `docker-compose up` is called once again it finally works. The issue needs to be investigated.
2.  `localDiscovery` often does not start on the first try due to `postgres` not being online yet. Waiting script needs to be implemented.
3.  `localPrometheus` container is forced to run as root as on default it's run as user `nobody` and cannot access data volume.
4.  `localPrometheus` expression browser is not easily accessible at the moment.

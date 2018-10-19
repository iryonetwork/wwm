# Running full clinic deployment on Windows

### Requirements

Setup was tested on Windows 10 Home & Pro Editions.

## 1. Initial setup.

### Windows 10 Pro

*   Install Git for Windows to be able to checkout the project.
*   Install and setup Docker for Windows. Choose default options to use Linux Containers on Windows.
*   Install Bonjour for Windows to get support for mDNS.
    *   The easiest way to get latest Bonjour for Windows is to install iTunes. It is possible to extract Bonjour-only _\*.msi_ installer by unzipping iTunes installer file.
*   Checkout IRYO WWM repo:
    Before checkout, while configuring git you should disable auto conversion to CRLF line endings.

        ```
        git config --global core.autocrlf false
        ```

    You should checkout IRYO WWM to `C:\iryo\wwm` not to have to change default `IRYO_WWM_DIR` docker-compose file environment variable

*   Go to Docker for Windows settings. In tab `Shared Drives` add the drive on which you checked out IRYO WWM repo to list of drives that can be available to Docker containers.

### Windows 10 Home

*   Download Docker Toolbox for Windows.
*   Unless you have them already installed choose to install Virtual Box and Git for Windows.
*   Run Docker Toolbox QuickStart Shell to create and setup docker machine.
*   Install Bonjour for Windows to get support for mDNS.
    *   The easiest way to get latest Bonjour for Windows is to install iTunes. It is possible to extract Bonjour-only _\*.msi_ installer by unzipping iTunes installer file.
*   Checkout IRYO WWM repo:
    Before checkout, while configuring git you should disable auto conversion to CRLF line endings.

    ```
    git config --global core.autocrlf false
    ```

*   Add IRYO WWM dir to shared folders for docker machine VM.
    The easiest way to do it currently is to open Virtual Box GUI and add path to WWM dir to visible there docker machine VM. You should mount WWM dir under `/iryo` in the VM not to have to change default `IRYO_WWM_DIR` docker-compose file environment variable.

## 2. Generate certificates and import root certificate to Windows root truststore.

*   Enter `docs/windowsClinic` directory in the admin-mode powershell.
*   Run `generateAndImportCerts.ps1` script:

```
powershell -ExecutionPolicy ByPass -File .\generateAndImportCerts.ps1
```

## 3. Setup location, clinic on cloud and import certificates.

1.  Setup location and clinic on cloud deployment that you intend to connect clinic to. Write down location ID and clinic ID.
2.  Configure certificates for authSync, storageSync and batchStorageSync generated in previous step so they will be accepted as valid by `cloudAuth`.

## 4. Set configuration values.

1.  Edit `frontendConfig.json` to include correct `clinicId` and `locationId`.
2.  Edit `.env` environment variables file for `docker-compose` it's included together with `docker-compose` in folders speciifc for Windows edition.
    The .env files contain default values for test Windows clinic deployment that is connecting to stagingcloud deployment.
    The values that are not filled in and has to be added before running the clinic are:
    *   `CLOUDSYMMETRIC_BASIC_AUTH_USERNAME` and `CLOUDSYMMETRIC_BASIC_AUTH_PASSWORD`.
        It needs to be set to correct username and password setup for `cloudSymmetric` server endpoints at the chosen cloud deployment.
    *   `AUTH_STORAGE_ENCRYPTION_KEY`
        It needs to be the same as at the chosen cloud deployment's `cloudAuth`. Otherwise `locatAuth` won't be able to decrypt received auth DB file.

## 6. Start clinic

While being in `docs/windowsClinic/home` (for Windows 10 Home) or `docs/windowsClinic/home` (for Windows 10 Pro) run in powershell:

```
docker-compose up -d
```

Now you should be able to access clinic web interface at `https://iryo.local`.

## Known issues

### Windows 10 Pro

1.  Clinic site is not accessible from Microsoft Edge browser due to Microsoft Edge not being able to correctly resolve local domains.

### Windows 10 Home

1.  Clinic site is not accessible from Microsoft Edge browser.
2.  `locaNats` often does not start on the first try. When `docker-compose up` is called once again it finally works. The issue needs to be investigated.
3.  `localDiscovery` often does not start on the first try due to `postgres` not being online yet. Waiting script needs to be implemented.
4.  `localPrometheus` container is forced to run as root as on default it's run as user `nobody` and cannot access data volume.
5.  `localPrometheus` expression browser is not easily accessible at the moment.
